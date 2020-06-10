package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/sap_api"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
)

type NewSapErpInput struct {
	Integration core.GenericIntegration         `json:"integration"`
	Setup       sap_api.SapRfcConnectionOptions `json:"setup"`
}

func newSapErpIntegration(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := NewSapErpInput{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	org, err := webcore.FindOrganizationInContext(r.Context())
	if err != nil {
		core.Warning("Can't find org in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	system, err := webcore.FindSystemInContext(r.Context())
	if err != nil {
		core.Warning("Can't find system in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inputs.Integration.OrgId = org.Id
	inputs.Integration.Type = core.ITSapErp

	inputs.Setup.Password, err = webcore.CreateEncryptedPassword(inputs.Setup.Password)
	if err != nil {
		core.Warning("Failed to encrypted password: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tx := database.CreateTx()
	err = database.WrapTx(tx, func() error {
		return database.NewIntegrationWithTx(tx, &inputs.Integration)
	}, func() error {
		return database.NewSapErpIntegrationWithTx(tx, inputs.Integration.Id, inputs.Setup)
	}, func() error {
		return database.LinkIntegrationWithSystemWithTx(tx, inputs.Integration.Id, system.Id, system.OrgId)
	})

	if err != nil {
		core.Warning("Failed to create integration: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(inputs.Integration)
}

type EditSapErpInput struct {
	Setup sap_api.SapRfcConnectionOptions `json:"setup"`
}

func editSapErpIntegration(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	integration, err := webcore.FindGenericIntegrationInContext(r.Context())
	if err != nil {
		core.Warning("Can't find generic integration in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inputs := EditSapErpInput{}
	err = webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	retPassword := inputs.Setup.Password
	inputs.Setup.Password, err = webcore.CreateEncryptedPassword(inputs.Setup.Password)
	if err != nil {
		core.Warning("Failed to encrypt password: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tx := database.CreateTx()
	err = database.WrapTx(tx, func() error {
		return database.EditSapErpIntegrationWithTx(tx, integration.Id, inputs.Setup)
	})
	if err != nil {
		core.Warning("Failed to edit SAP ERP integration: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	inputs.Setup.Password = retPassword
	jsonWriter.Encode(inputs.Setup)
}

func getSapErpIntegration(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	integration, err := webcore.FindGenericIntegrationInContext(r.Context())
	if err != nil {
		core.Warning("Can't find generic integration in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sap, err := database.GetSapErpIntegration(integration.Id)
	if err != nil {
		core.Warning("Can't get SAP ERP integration: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sap.Password, err = webcore.DecryptEncryptedPassword(sap.Password)
	if err != nil {
		core.Warning("Failed to decrypt password: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(sap)
}
