package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
)

type AllControlGeneralLedgerAccountLinkInputs struct {
	ControlId core.NullInt64 `webcore:"controlId,optional"`
	AccountId core.NullInt64 `webcore:"accountId,optional"`
	OrgId     int32          `webcore:"orgId"`
}

func allControlGeneralLedgerAccountLinks(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := AllControlGeneralLedgerAccountLinkInputs{}
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

	if inputs.ControlId.NullInt64.Valid {
		accounts, err := database.FindGeneralLedgerAccountsLinkedToControl(inputs.ControlId.NullInt64.Int64, inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to get linked accounts: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		categories, err := database.GetOrgGLCategories(inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to get org categories: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsonWriter.Encode(struct {
			Accounts   []*core.GeneralLedgerAccount
			Categories []*core.GeneralLedgerCategory
		}{
			Accounts:   accounts,
			Categories: categories,
		})
	} else if inputs.AccountId.NullInt64.Valid {
		Controls, err := database.FindControlsLinkedToGeneralLedgerAccount(inputs.AccountId.NullInt64.Int64, inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to get linked Controls: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		jsonWriter.Encode(struct {
			Controls []*core.Control
		}{
			Controls: Controls,
		})
	} else {
		core.Warning("Invalid combination of inputs.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
