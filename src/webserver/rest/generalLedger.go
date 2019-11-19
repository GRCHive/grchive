package rest

import (
	"encoding/json"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
	"net/http"
)

type NewGLCategoryInputs struct {
	OrgId            int32          `json:"orgId"`
	ParentCategoryId core.NullInt64 `json:"parentCategoryId"`
	Name             string         `json:"name"`
	Description      string         `json:"description"`
}

type EditGLCategoryInputs struct {
	CatId            int64          `json:"catId"`
	OrgId            int32          `json:"orgId"`
	ParentCategoryId core.NullInt64 `json:"parentCategoryId"`
	Name             string         `json:"name"`
	Description      string         `json:"description"`
}

type DeleteGLCategoryInputs struct {
	CatId int64 `json:"catId"`
	OrgId int32 `json:"orgId"`
}

type NewGLAccountInputs struct {
	OrgId               int32  `json:"orgId"`
	ParentCategoryId    int64  `json:"parentCategoryId"`
	AccountName         string `json:"accountName"`
	AccountId           string `json:"accountId"`
	AccountDescription  string `json:"accountDescription"`
	FinanciallyRelevant bool   `json:"financiallyRelevant"`
}

type GetGLInputs struct {
	OrgId int32 `webcore:"orgId"`
}

func deleteGLCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	inputs := DeleteGLCategoryInputs{}
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

	err = database.DeleteGLCategory(inputs.CatId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to delete GL Cat: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func editGLCategory(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := EditGLCategoryInputs{}
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

	cat := core.GeneralLedgerCategory{
		Id:               inputs.CatId,
		OrgId:            inputs.OrgId,
		ParentCategoryId: inputs.ParentCategoryId,
		Name:             inputs.Name,
		Description:      inputs.Description,
	}

	err = database.UpdateGLCategory(&cat, role)
	if err != nil {
		core.Warning("Can't update GL category: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(cat)
}

func createNewGLCategory(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := NewGLCategoryInputs{}
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

	cat := core.GeneralLedgerCategory{
		OrgId:            inputs.OrgId,
		ParentCategoryId: inputs.ParentCategoryId,
		Name:             inputs.Name,
		Description:      inputs.Description,
	}

	err = database.CreateNewGLCategory(&cat, role)
	if err != nil {
		core.Warning("Can't create new GL category: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(cat)
}

func createNewGLAccount(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := NewGLAccountInputs{}
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

	acc := core.GeneralLedgerAccount{
		OrgId:               inputs.OrgId,
		ParentCategoryId:    inputs.ParentCategoryId,
		AccountId:           inputs.AccountId,
		AccountName:         inputs.AccountName,
		AccountDescription:  inputs.AccountDescription,
		FinanciallyRelevant: inputs.FinanciallyRelevant,
	}

	err = database.CreateNewGLAccount(&acc, role)
	if err != nil {
		core.Warning("Can't create new GL account: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(acc)
}

func getGL(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetGLInputs{}
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

	cats, err := database.GetOrgGLCategories(inputs.OrgId, role)
	if err != nil {
		core.Warning("Can't get GL categories: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	accs, err := database.GetOrgGLAccounts(inputs.OrgId, role)
	if err != nil {
		core.Warning("Can't get GL accounts: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	type GetGLOutputs struct {
		Categories []*core.GeneralLedgerCategory
		Accounts   []*core.GeneralLedgerAccount
	}
	outputs := GetGLOutputs{cats, accs}
	jsonWriter.Encode(outputs)
}
