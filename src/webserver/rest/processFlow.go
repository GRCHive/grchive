package rest

import (
	"encoding/json"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
	"net/http"
	"strconv"
	"time"
)

func getAllProcessFlows(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	queryVals := r.URL.Query()

	organizationName, ok := queryVals["organization"]
	if !ok || len(organizationName) == 0 {
		core.Warning("Failed to get process flows (no organization)")
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	organization, err := database.FindOrganizationFromGroupName(organizationName[0])
	if err != nil {
		core.Warning("No organization: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	// Ensure that the user has access to this organization.
	userParsedData, err := webcore.FindSessionParsedDataInContext(r.Context())
	if err != nil {
		core.Warning("No user session data: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	if userParsedData.Org.OktaGroupId != organization.OktaGroupId {
		core.Warning("Unauthorized access")
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	requestedId, ok := queryVals["requested"]

	var flows []*core.ProcessFlow
	var index uint32 = 0

	if !ok || len(requestedId) == 0 {
		flows, err = database.FindOrganizationProcessFlows(organization)
		if err != nil {
			core.Warning("Database error [0]: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			jsonWriter.Encode(struct{}{})
			return
		}
	} else {
		intRequestedId, err := strconv.ParseUint(requestedId[0], 10, 32)
		if err != nil {
			core.Warning("Invalid requested id: " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			jsonWriter.Encode(struct{}{})
			return
		}

		flows, index, err = database.FindOrganizationProcessFlowsWithIndex(organization, uint32(intRequestedId))
		if err != nil {
			core.Warning("Database error [1]: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			jsonWriter.Encode(struct{}{})
			return
		}
	}

	jsonWriter.Encode(struct {
		Flows          []*core.ProcessFlow
		RequestedIndex uint32
	}{
		Flows:          flows,
		RequestedIndex: index,
	})
}

func newProcessFlow(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	// Retrieve name, description, and organization ID from the post data.
	if err := r.ParseForm(); err != nil || len(r.PostForm) == 0 {
		core.Warning("Failed to parse form data: " + core.ErrorString(err))
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	nameData := r.PostForm["name"]
	descriptionData := r.PostForm["description"]
	orgIdData := r.PostForm["organization"]

	if len(nameData) == 0 || len(descriptionData) == 0 || len(orgIdData) == 0 {
		core.Warning("Empty name or description or organization.")
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	org, err := database.FindOrganizationFromGroupName(orgIdData[0])
	if err != nil {
		core.Warning("No organization: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	// Ensure that the user has access to this organization.
	userParsedData, err := webcore.FindSessionParsedDataInContext(r.Context())
	if err != nil {
		core.Warning("No user session data: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	if userParsedData.Org.OktaGroupId != org.OktaGroupId {
		core.Warning("Unauthorized access")
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	newFlow := core.ProcessFlow{
		Name:            nameData[0],
		Org:             org,
		Description:     descriptionData[0],
		CreationTime:    time.Now(),
		LastUpdatedTime: time.Now(),
	}

	err = database.InsertNewProcessFlow(&newFlow)
	if err != nil {
		if database.IsDuplicateDBEntry(err) {
			w.WriteHeader(http.StatusBadRequest)
			jsonWriter.Encode(struct {
				IsDuplicate bool
			}{
				true,
			})
		} else {
			core.Warning("Database error: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			jsonWriter.Encode(struct{}{})
		}
		return
	}

	// Echo the ID and Name back to the requester.
	w.WriteHeader(http.StatusOK)
	jsonWriter.Encode(struct {
		Name string
		Id   uint32
	}{
		newFlow.Name,
		newFlow.Id,
	})
}
