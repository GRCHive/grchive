package rest

import (
	"encoding/json"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
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
	OrgId          int32          `webcore:"orgId"`
	DeploymentType core.NullInt32 `webcore:"deploymentType,optional"`
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
	DbId       int64  `json:"dbId"`
	OrgId      int32  `json:"orgId"`
	ConnString string `json:"connectionString"`
	Username   string `json:"username"`
	Password   string `json:"password"`
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
		dbs, err = database.GetAllDatabasesForOrg(org.Id, role)
	} else {
		dbs, err = database.GetAllDatabasesForOrgWithDeployment(org.Id, inputs.DeploymentType.NullInt32.Int32, role)
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

	apiKey, err := webcore.GetAPIKeyFromRequest(r)
	if apiKey == nil || err != nil {
		core.Warning("No API Key: " + core.ErrorString(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

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

	encPassword, salt, err := webcore.CreateSaltedEncryptedPassword(inputs.Password)
	if err != nil {
		core.Warning("Failed to encrypt password: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	conn := core.DatabaseConnection{
		DbId:       inputs.DbId,
		OrgId:      org.Id,
		ConnString: inputs.ConnString,
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
