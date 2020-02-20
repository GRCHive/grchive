package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
)

type AllRiskGeneralLedgerAccountLinkInputs struct {
	RiskId    core.NullInt64 `webcore:"riskId,optional"`
	AccountId core.NullInt64 `webcore:"accountId,optional"`
	OrgId     int32          `webcore:"orgId"`
}

func allRiskGeneralLedgerAccountLinks(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := AllRiskGeneralLedgerAccountLinkInputs{}
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

	if inputs.RiskId.NullInt64.Valid {
		accounts, err := database.FindGeneralLedgerAccountsLinkedToRisk(inputs.RiskId.NullInt64.Int64, inputs.OrgId, role)
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
		risks, err := database.FindRisksLinkedToGeneralLedgerAccount(inputs.AccountId.NullInt64.Int64, inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to get linked risks: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		jsonWriter.Encode(struct {
			Risks []*core.Risk
		}{
			Risks: risks,
		})
	} else {
		core.Warning("Invalid combination of inputs.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
