package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
)

type NewServerInputs struct {
	OrgId           int32  `json:"orgId"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	IpAddress       string `json:"ip"`
	OperatingSystem string `json:"os"`
	Location        string `json:"location"`
}

type UpdateServerInputs struct {
	ServerId        int64  `webcore:"serverId"`
	OrgId           int32  `json:"orgId"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	IpAddress       string `json:"ip"`
	OperatingSystem string `json:"os"`
	Location        string `json:"location"`
}

type DeleteServerInputs struct {
	ServerId int64 `webcore:"serverId"`
	OrgId    int32 `webcore:"orgId"`
}

type AllServersInputs struct {
	OrgId int32 `webcore:"orgId"`
}

type GetServerInputs struct {
	ServerId int64 `webcore:"serverId"`
	OrgId    int32 `webcore:"orgId"`
}

func newServer(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := NewServerInputs{}
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

	s := core.Server{
		OrgId:           inputs.OrgId,
		Name:            inputs.Name,
		Description:     inputs.Description,
		IpAddress:       inputs.IpAddress,
		OperatingSystem: inputs.OperatingSystem,
		Location:        inputs.Location,
	}

	err = database.NewServer(&s, role)
	if err != nil {
		core.Warning("Failed to create new server: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(s)
}

func updateServer(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := UpdateServerInputs{}
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

	s := core.Server{
		Id:              inputs.ServerId,
		OrgId:           inputs.OrgId,
		Name:            inputs.Name,
		Description:     inputs.Description,
		IpAddress:       inputs.IpAddress,
		OperatingSystem: inputs.OperatingSystem,
		Location:        inputs.Location,
	}

	err = database.UpdateServer(&s, role)
	if err != nil {
		core.Warning("Failed to update server: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(s)
}

func allServers(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := AllServersInputs{}
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

	servers, err := database.AllServersForOrganization(inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to retrieve servers: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(servers)
}

func deleteServer(w http.ResponseWriter, r *http.Request) {
	inputs := DeleteServerInputs{}
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

	err = database.DeleteServer(inputs.ServerId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to delete server: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getServer(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetServerInputs{}
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

	server, err := database.GetServer(inputs.ServerId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to retrieve server: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	relevantSystems, err := database.GetSystemsLinkedToServer(server.Id, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get systems linked to server: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	relevantDbs, err := database.GetDatabasesLinkedToServer(server.Id, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get dbs linked to server: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(struct {
		Server          *core.Server
		RelevantSystems []*core.System
		RelevantDbs     []*core.Database
	}{
		Server:          server,
		RelevantSystems: relevantSystems,
		RelevantDbs:     relevantDbs,
	})
}
