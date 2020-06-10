package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
)

func allIntegrations(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	system, err := webcore.FindSystemInContext(r.Context())
	if err != nil {
		core.Warning("Can't find system in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	integrations, err := database.AllGenericIntegrationsForSystem(system.Id)
	if err != nil {
		core.Warning("Can't find integrations for system: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(integrations)
}

type EditGenericIntegrationInput struct {
	Data core.GenericIntegration `json:"data"`
}

func editGenericIntegration(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := EditGenericIntegrationInput{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	integration, err := webcore.FindGenericIntegrationInContext(r.Context())
	if err != nil {
		core.Warning("Can't find generic integration in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	integration.Name = inputs.Data.Name
	integration.Description = inputs.Data.Description

	tx := database.CreateTx()
	err = database.WrapTx(tx, func() error {
		return database.EditGenericIntegrationWithTx(tx, integration)
	})

	if err != nil {
		core.Warning("Failed to edit integration: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(integration)
}

func deleteGenericIntegration(w http.ResponseWriter, r *http.Request) {
	integration, err := webcore.FindGenericIntegrationInContext(r.Context())
	if err != nil {
		core.Warning("Can't find generic integration in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tx := database.CreateTx()
	err = database.WrapTx(tx, func() error {
		return database.DeleteGenericIntegrationWithTx(tx, integration.Id)
	})

	if err != nil {
		core.Warning("Failed to delete generic integration: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
