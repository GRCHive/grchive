package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
)

func getOrgPbcNotificationSettings(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	org, err := webcore.FindOrganizationInContext(r.Context())
	if err != nil {
		core.Warning("No organization: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	settings, err := database.GetOrgPbcNotificationCadenceSettings(org.Id)
	if err != nil {
		core.Warning("Failed to get org pbc notification settings: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(settings)
}

type NewPbcNotificationSettingInputs struct {
	DaysBeforeDue int32 `json:"daysBefore"`
}

func newOrgPbcNotificationSetting(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := NewPbcNotificationSettingInputs{}
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

	var settings *core.PbcNotificationCadenceSettings

	tx := database.CreateTx()
	err = database.WrapTx(tx, func() error {
		var ierr error
		settings, ierr = database.CreateNewPbcNotificationCadenceSettingWithTx(tx, org.Id, inputs.DaysBeforeDue)
		return ierr
	})

	if err != nil {
		core.Warning("Failed to create PBC notification cadence setting: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(settings)
}

func deleteOrgPbcNotificationSetting(w http.ResponseWriter, r *http.Request) {
	setting, err := webcore.FindPbcNotificationCadenceSettingsInContext(r.Context())
	if err != nil {
		core.Warning("Failed to get org pbc notification settings: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tx := database.CreateTx()
	err = database.WrapTx(tx, func() error {
		return database.DeletePbcNotificationCadenceSettingsWithTx(tx, setting.Id)
	})

	if err != nil {
		core.Warning("Failed to delete PBC notification cadence setting: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type EditPbcNotificationSettingInputs struct {
	Setting  core.PbcNotificationCadenceSettings `json:"setting"`
	ApplyAll bool                                `json:"applyAll"`
}

func editOrgPbcNotificationSetting(w http.ResponseWriter, r *http.Request) {
	setting, err := webcore.FindPbcNotificationCadenceSettingsInContext(r.Context())
	if err != nil {
		core.Warning("Failed to get org pbc notification settings: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inputs := EditPbcNotificationSettingInputs{}
	err = webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inputs.Setting.Id = setting.Id
	inputs.Setting.OrgId = setting.OrgId

	tx := database.CreateTx()
	err = database.WrapTx(tx, func() error {
		return database.EditPbcNotificationCadenceSettingWithTx(tx, inputs.Setting)
	}, func() error {
		if !inputs.ApplyAll {
			return nil
		}
		return database.EditAllPbcNotificationCadenceSettingWithTx(tx, inputs.Setting)
	})

	if err != nil {
		core.Warning("Failed to delete PBC notification cadence setting: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
