package rest

import (
	"encoding/json"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
	"net/http"
)

type NewVendorProductInputs struct {
	OrgId       int32  `json:"orgId"`
	VendorId    int64  `json:"vendorId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Url         string `json:"url"`
}

func newVendorProduct(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := NewVendorProductInputs{}
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

	product := core.VendorProduct{
		VendorId:    inputs.VendorId,
		OrgId:       inputs.OrgId,
		Name:        inputs.Name,
		Description: inputs.Description,
		Url:         inputs.Url,
	}

	err = database.NewVendorProduct(&product, role)
	if err != nil {
		core.Warning("Failed to create new vendor product: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(product)
}

type AllVendorProductInputs struct {
	OrgId    int32 `webcore:"orgId"`
	VendorId int64 `webcore:"vendorId"`
}

func allVendorProducts(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := AllVendorProductInputs{}
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

	products, err := database.AllVendorProductsForVendor(inputs.VendorId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to create get vendor products: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(products)
}

type GetVendorProductInputs struct {
	ProductId int64 `webcore:"productId"`
	VendorId  int64 `webcore:"vendorId"`
	OrgId     int32 `webcore:"orgId"`
}

func getVendorProduct(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetVendorProductInputs{}
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

	product, err := database.GetVendorProduct(inputs.ProductId, inputs.VendorId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get vendor product: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	socFiles, err := database.GetSocDocumentationForVendorProduct(inputs.ProductId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get vendor product SOC files: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(struct {
		Product  *core.VendorProduct
		SocFiles []*core.ControlDocumentationFile
	}{
		Product:  product,
		SocFiles: socFiles,
	})
}

type UpdateVendorProductInputs struct {
	ProductId   int64  `json:"productId"`
	OrgId       int32  `json:"orgId"`
	VendorId    int64  `json:"vendorId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Url         string `json:"url"`
}

func updateVendorProduct(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := UpdateVendorProductInputs{}
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

	product := core.VendorProduct{
		Id:          inputs.ProductId,
		VendorId:    inputs.VendorId,
		OrgId:       inputs.OrgId,
		Name:        inputs.Name,
		Description: inputs.Description,
		Url:         inputs.Url,
	}

	err = database.UpdateVendorProduct(&product, role)
	if err != nil {
		core.Warning("Failed to update vendor product: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(product)
}

type DeleteVendorProductInputs struct {
	ProductId int64 `json:"productId"`
	VendorId  int64 `json:"vendorId"`
	OrgId     int32 `json:"orgId"`
}

func deleteVendorProduct(w http.ResponseWriter, r *http.Request) {
}

type LinkVendorProductSocInputs struct {
	ProductId int64                                  `json:"productId"`
	OrgId     int32                                  `json:"orgId"`
	SocFiles  []*core.ControlDocumentationFileHandle `json:"socFiles"`
}

func linkVendorProductSoc(w http.ResponseWriter, r *http.Request) {
	inputs := LinkVendorProductSocInputs{}
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

	err = database.LinkVendorProductToSocFiles(inputs.ProductId, inputs.OrgId, inputs.SocFiles, role)
	if err != nil {
		core.Warning("Failed to link vendor product to soc files: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func unlinkVendorProductSoc(w http.ResponseWriter, r *http.Request) {
	inputs := LinkVendorProductSocInputs{}
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

	err = database.UnlinkVendorProductFromSocFiles(inputs.ProductId, inputs.OrgId, inputs.SocFiles, role)
	if err != nil {
		core.Warning("Failed to unlink vendor product from soc files: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
