package rest

import (
	"encoding/json"
	"fmt"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
)

type NewDatabaseInputs struct {
	Name      string `json:"name"`
	OrgId     int32  `json:"orgId"`
	TypeId    int32  `json:"typeId"`
	OtherType string `json:"otherType"`
	Version   string `json:"version"`
}

type GetAllDatabaseInputs struct {
	OrgId          int32                   `webcore:"orgId"`
	DeploymentType core.NullInt32          `webcore:"deploymentType,optional"`
	Filter         core.DatabaseFilterData `webcore:"filter"`
}

type GetDatabaseInputs struct {
	DbId  int64 `webcore:"dbId"`
	OrgId int32 `webcore:"orgId"`
}

type EditDatabaseInputs struct {
	DbId      int64  `json:"dbId"`
	Name      string `json:"name"`
	OrgId     int32  `json:"orgId"`
	TypeId    int32  `json:"typeId"`
	OtherType string `json:"otherType"`
	Version   string `json:"version"`
}

type DeleteDatabaseInputs struct {
	DbId  int64 `json:"dbId"`
	OrgId int32 `json:"orgId"`
}

type NewDbConnectionInputs struct {
	DbId       int64             `json:"dbId"`
	OrgId      int32             `json:"orgId"`
	Host       string            `json:"host"`
	Port       int32             `json:"port"`
	DbName     string            `json:"dbName"`
	Parameters map[string]string `json:"parameters"`
	Username   string            `json:"username"`
	Password   string            `json:"password"`
}

type DeleteDatabaseConnectionInputs struct {
	ConnId int64 `json:"connId"`
	DbId   int64 `json:"dbId"`
	OrgId  int32 `json:"orgId"`
}

type LinkSystemsInputs struct {
	DbId   int64   `json:"dbId"`
	OrgId  int32   `json:"orgId"`
	SysIds []int64 `json:"sysIds"`
}

func newDb(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := NewDatabaseInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	org, err := database.FindOrganizationFromId(inputs.OrgId)
	if err != nil {
		core.Warning("No organization: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, org.Id)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	db := core.Database{
		Name:      inputs.Name,
		OrgId:     inputs.OrgId,
		TypeId:    inputs.TypeId,
		OtherType: inputs.OtherType,
		Version:   inputs.Version,
	}

	err = database.InsertNewDatabase(&db, role)
	if err != nil {
		core.Warning("Can't insert new database: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(db)
}

func getAllDb(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetAllDatabaseInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	org, err := database.FindOrganizationFromId(inputs.OrgId)
	if err != nil {
		core.Warning("No organization: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, org.Id)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var dbs []*core.Database

	if !inputs.DeploymentType.NullInt32.Valid {
		dbs, err = database.GetAllDatabasesForOrg(org.Id, inputs.Filter, role)
	} else {
		dbs, err = database.GetAllDatabasesForOrgWithDeployment(org.Id, inputs.DeploymentType.NullInt32.Int32, inputs.Filter, role)
	}
	if err != nil {
		core.Warning("Can't get all databases: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonWriter.Encode(dbs)
}

func getDbTypes(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	types, err := database.GetAllSupportedDatabaseTypes(core.ServerRole)
	if err != nil {
		core.Warning("Can't get types: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonWriter.Encode(types)
}

func getDb(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetDatabaseInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	org, err := database.FindOrganizationFromId(inputs.OrgId)
	if err != nil {
		core.Warning("No organization: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, org.Id)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	db, err := database.GetDb(inputs.DbId, org.Id, role)
	if err != nil {
		core.Warning("Can't get database: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	conn, err := database.FindDatabaseConnectionForDatabase(inputs.DbId, org.Id, role)
	if err != nil {
		core.Warning("Can't get database connection: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	allSystems, err := database.GetAllSystemsForOrg(org.Id, role)
	if err != nil {
		core.Warning("Failed to obtain systems: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sysIds, err := database.FindSystemIdsForDatabase(db.Id, org.Id, role)
	if err != nil {
		core.Warning("Failed to find relevant systems: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	deployment, err := database.GetDatabaseDeployment(db.Id, org.Id, role)
	if err != nil {
		core.Warning("Failed to find relevant deployment: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(struct {
		Database          *core.Database
		Connection        *core.DatabaseConnection
		RelevantSystemIds []int64
		AllSystems        []*core.System
		Deployment        *core.FullDeployment
	}{
		Database:          db,
		Connection:        conn,
		RelevantSystemIds: sysIds,
		AllSystems:        allSystems,
		Deployment:        deployment,
	})
}

func editDb(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := EditDatabaseInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	org, err := database.FindOrganizationFromId(inputs.OrgId)
	if err != nil {
		core.Warning("No organization: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, org.Id)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	db := core.Database{
		Id:        inputs.DbId,
		Name:      inputs.Name,
		OrgId:     inputs.OrgId,
		TypeId:    inputs.TypeId,
		OtherType: inputs.OtherType,
		Version:   inputs.Version,
	}

	err = database.EditDb(&db, role)
	if err != nil {
		core.Warning("Can't edit database: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(db)
}

func deleteDb(w http.ResponseWriter, r *http.Request) {
	inputs := DeleteDatabaseInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	org, err := database.FindOrganizationFromId(inputs.OrgId)
	if err != nil {
		core.Warning("No organization: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, org.Id)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = database.DeleteDb(inputs.DbId, org.Id, role)
	if err != nil {
		core.Warning("Can't delete database: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func newDbConnection(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := NewDbConnectionInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, inputs.OrgId)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	encPassword, salt, err := webcore.CreateSaltedEncryptedPassword(inputs.Password)
	if err != nil {
		core.Warning("Failed to encrypt password: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	conn := core.DatabaseConnection{
		DbId:       inputs.DbId,
		OrgId:      inputs.OrgId,
		Host:       inputs.Host,
		Port:       inputs.Port,
		DbName:     inputs.DbName,
		Parameters: inputs.Parameters,
		Username:   inputs.Username,
		Password:   encPassword,
		Salt:       salt,
	}

	err = database.InsertNewDatabaseConnection(&conn, role)
	if err != nil {
		core.Warning("Failed to insert db connection: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Decrypt for sending back to client.
	conn.Password, err = webcore.DecryptSaltedEncryptedPassword(conn.Password, conn.Salt)
	if err != nil {
		core.Warning("Failed to decrypt password: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Doesn't particularly matter if creating the database refresh here fails
	// since the user can just request a refresh later if necessary.
	webcore.CreateNewDatabaseRefresh(inputs.DbId, inputs.OrgId, role)

	jsonWriter.Encode(conn)
}

func deleteDbConnection(w http.ResponseWriter, r *http.Request) {
	inputs := DeleteDatabaseConnectionInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	org, err := database.FindOrganizationFromId(inputs.OrgId)
	if err != nil {
		core.Warning("No organization: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, org.Id)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = database.DeleteDatabaseConnection(inputs.ConnId, inputs.DbId, org.Id, role)
	if err != nil {
		core.Warning("Failed to delete db connection: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func linkSystemsToDatabase(w http.ResponseWriter, r *http.Request) {
	inputs := LinkSystemsInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	org, err := database.FindOrganizationFromId(inputs.OrgId)
	if err != nil {
		core.Warning("No organization: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, org.Id)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = database.LinkSystemsToDatabase(inputs.DbId, org.Id, inputs.SysIds, role)
	if err != nil {
		core.Warning("Failed to link systems to database: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getDatabaseSettings(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	db, err := webcore.FindDatabaseInContext(r.Context())
	if err != nil {
		core.Warning("Failed find database in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	settings, err := database.GetDatabaseSettings(db.Id)
	if err != nil {
		core.Warning("Failed find get database settings: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(settings)
}

type EditDatabaseSettingsInputs struct {
	AutoRefreshEnabled  bool                        `json:"autoRefreshEnabled"`
	AutoRefreshSchedule *core.ScheduledTaskRawInput `json:"autoRefreshSchedule"`
}

func editDatabaseSettings(w http.ResponseWriter, r *http.Request) {
	db, err := webcore.FindDatabaseInContext(r.Context())
	if err != nil {
		core.Warning("Failed find database in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.FindRoleInContext(r.Context())
	if err != nil {
		core.Warning("Failed find get role in context: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	settings, err := database.GetDatabaseSettings(db.Id)
	if err != nil {
		core.Warning("Failed find get database settings: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	inputs := EditDatabaseSettingsInputs{}
	err = webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if inputs.AutoRefreshSchedule != nil {
		inputs.AutoRefreshSchedule.Name = fmt.Sprintf("Refresh DB: %s", db.Name)
	}

	tx := database.CreateTx()

	isEditingAutoRefresh := settings.AutoRefreshEnabled && inputs.AutoRefreshEnabled
	isChangeAutoRefresh := settings.AutoRefreshEnabled != inputs.AutoRefreshEnabled

	// This function needs to handle three cases:
	// 	1) Creating a new task where no task existed before
	//  2) Removing a task when there was a task already
	//  3) Editing a task when there was a task already
	// There doesn't exist functionality to edit a task, so #3 needs to be
	// handles by doing a delete followed creationg of a new task. Therefore
	// we can handle case #3 by performing #2 followed by #1.
	err = database.WrapTx(tx, func() error {
		// Removing a task happens when the user disables auto-refresh
		// and the current refresh setting is enabled or user is editing.
		if !isEditingAutoRefresh && !(isChangeAutoRefresh && !inputs.AutoRefreshEnabled) {
			return nil
		}

		storedId := settings.AutoRefreshTaskId.NullInt64.Int64
		settings.AutoRefreshTaskId = core.NullInt64{}

		webcore.DefaultRabbitMQ.SendMessage(webcore.PublishMessage{
			Exchange: webcore.DEFAULT_EXCHANGE,
			Queue:    webcore.TASK_MANAGER_QUEUE,
			Body: webcore.TaskManagerMessage{
				Action: "Delete",
				TaskId: storedId,
			},
		})

		return database.DeleteScheduledTaskWithTx(tx, storedId)
	}, func() error {
		if !isEditingAutoRefresh && !(isChangeAutoRefresh && inputs.AutoRefreshEnabled) {
			return nil
		}

		metadata, err := webcore.CreateScheduledTaskFromRawInputs(
			tx,
			inputs.AutoRefreshSchedule,
			core.KGrchiveApiTask,
			core.GrchiveApiTaskData{
				Endpoint: webcore.MustGetRouteUrl(webcore.ApiDbRefreshRouteName),
				Method:   "POST",
				Payload: NewRefreshInput{
					DbId:  db.Id,
					OrgId: db.OrgId,
				},
			},
			role.UserId,
			db.OrgId,
			webcore.TaskLinkOptions{
				DbId: core.CreateNullInt64(db.Id),
			},
		)

		if err != nil {
			return err
		}

		settings.AutoRefreshTaskId = core.CreateNullInt64(metadata.Id)
		return nil
	})

	if err != nil {
		core.Warning("Failed to edit database settings: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// This needs to be here just in case the RabbitMQ gets sent and processed
	// before the database transaction is processed.
	if settings.AutoRefreshTaskId.NullInt64.Valid {
		webcore.DefaultRabbitMQ.SendMessage(webcore.PublishMessage{
			Exchange: webcore.DEFAULT_EXCHANGE,
			Queue:    webcore.TASK_MANAGER_QUEUE,
			Body: webcore.TaskManagerMessage{
				Action: "Add",
				TaskId: settings.AutoRefreshTaskId.NullInt64.Int64,
			},
		})
	}
}
