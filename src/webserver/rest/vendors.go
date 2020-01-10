package rest

import (
	"encoding/json"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
	"net/http"
)

type NewVendorInputs struct {
	OrgId       int32  `json:"orgId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Url         string `json:"url"`
}

func newVendor(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := NewVendorInputs{}
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

	// Need to create the vendor AND create a document
	// category specific for the vendor.
	vendor := core.Vendor{
		OrgId:       inputs.OrgId,
		Name:        inputs.Name,
		Description: inputs.Description,
		Url:         inputs.Url,
	}

	category := vendor.CreateDocumentationCategory()

	tx := database.CreateTx()

	err = database.NewControlDocumentationCategoryWithTx(&category, role, tx)
	if err != nil {
		tx.Rollback()
		core.Warning("Failed to create new doc category: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	vendor.DocCatId = category.Id

	err = database.NewVendorWithTx(&vendor, role, tx)
	if err != nil {
		tx.Rollback()
		core.Warning("Failed to create new vendor: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if tx.Commit() != nil {
		tx.Rollback()
		core.Warning("Failed to commit changes: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(vendor)
}

type AllVendorInputs struct {
	OrgId int32 `webcore:"orgId"`
}

func allVendors(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := AllVendorInputs{}
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

	vendors, err := database.AllVendorsForOrganization(inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get vendors: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(vendors)
}

type GetVendorInputs struct {
	VendorId int64 `webcore:"vendorId"`
	OrgId    int32 `webcore:"orgId"`
}

func getVendor(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetVendorInputs{}
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

	vendor, err := database.GetVendorFromId(inputs.VendorId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get vendor: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(vendor)
}

type UpdateVendorInputs struct {
	VendorId    int64  `json:"vendorId"`
	OrgId       int32  `json:"orgId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Url         string `json:"url"`
}

func updateVendor(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := UpdateVendorInputs{}
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

	// Need to create the vendor AND create a document
	// category specific for the vendor.
	vendor := core.Vendor{
		Id:          inputs.VendorId,
		OrgId:       inputs.OrgId,
		Name:        inputs.Name,
		Description: inputs.Description,
		Url:         inputs.Url,
	}

	// Need this since we need to return a complete vendor object back.
	vendor.DocCatId, err = database.GetDocCatIdForVendor(inputs.VendorId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get vendor id: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = database.UpdateVendor(&vendor, role)
	if err != nil {
		core.Warning("Failed to update vendor: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(vendor)
}

type DeleteVendorInputs struct {
	VendorId int64 `json:"vendorId"`
	OrgId    int32 `json:"orgId"`
}

func deleteVendor(w http.ResponseWriter, r *http.Request) {
	inputs := DeleteVendorInputs{}
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

	err = database.DeleteVendor(inputs.VendorId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to delete vendor: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func newVendorProduct(w http.ResponseWriter, r *http.Request) {
}

func allVendorProducts(w http.ResponseWriter, r *http.Request) {
}

func getVendorProduct(w http.ResponseWriter, r *http.Request) {
}

func updateVendorProduct(w http.ResponseWriter, r *http.Request) {
}

func deleteVendorProduct(w http.ResponseWriter, r *http.Request) {
}
