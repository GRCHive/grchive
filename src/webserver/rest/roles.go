package rest

import (
	"encoding/json"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
	"net/http"
)

type NewRoleInputs struct {
	OrgId       int32               `json:"orgId"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Permissions core.PermissionsMap `json:"permissions"`
}

type GetOrganizationRolesInputs struct {
	OrgId int32 `webcore:"orgId"`
}

func getAllOrganizationRoles(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetOrganizationRolesInputs{}
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

	availableRoles, err := database.FindRolesForOrg(org.Id, role)
	if err != nil {
		core.Warning("Failed to find roles for org: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(availableRoles)
}

func editRole(w http.ResponseWriter, r *http.Request) {
}

func deleteRole(w http.ResponseWriter, r *http.Request) {
}

func newRole(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := NewRoleInputs{}
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

	newRole := core.Role{
		RoleMetadata: core.RoleMetadata{
			Name:        inputs.Name,
			Description: inputs.Description,
			OrgId:       inputs.OrgId,
			IsDefault:   false,
			IsAdmin:     false,
		},
		Permissions: inputs.Permissions,
	}

	err = database.InsertOrgRole(&newRole.RoleMetadata, &newRole, role)
	if err != nil {
		core.Warning("Failed to insert new role: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(newRole)
}
