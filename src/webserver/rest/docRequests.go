package rest

import (
	"encoding/json"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
	"net/http"
	"time"
)

type NewDocumentRequestInputs struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	CatId           int64  `json:"catId"`
	OrgId           int32  `json:"orgId"`
	RequestedUserId int64  `json:"requestedUserId"`
	VendorProductId int64  `json:"vendorProductId"`
}

type UpdateDocumentRequestInputs struct {
	RequestId       int64  `json:"requestId"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	CatId           int64  `json:"catId"`
	OrgId           int32  `json:"orgId"`
	RequestedUserId int64  `json:"requestedUserId"`
}

type GetDocumentRequestInputs struct {
	RequestId int64 `webcore:"requestId"`
	OrgId     int32 `webcore:"orgId"`
}

type DeleteDocumentRequestInputs struct {
	RequestId int64 `json:"requestId"`
	OrgId     int32 `json:"orgId"`
}

type CompleteDocumentRequestInputs struct {
	RequestId int64 `json:"requestId"`
	OrgId     int32 `json:"orgId"`
	Complete  bool  `json:"complete"`
}

type AllDocumentRequestsInputs struct {
	OrgId           int32          `webcore:"orgId"`
	CatId           core.NullInt64 `webcore:"catId,optional"`
	VendorProductId core.NullInt64 `webcore:"vendorProductId,optional"`
}

func newDocumentRequest(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := NewDocumentRequestInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, inputs.OrgId)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	request := core.DocumentRequest{
		Name:            inputs.Name,
		Description:     inputs.Description,
		CatId:           inputs.CatId,
		OrgId:           inputs.OrgId,
		RequestedUserId: inputs.RequestedUserId,
		RequestTime:     time.Now().UTC(),
	}

	tx := database.CreateTx()

	err = database.CreateNewDocumentRequestWithTx(&request, role, tx)
	if err != nil {
		tx.Rollback()
		core.Warning("Failed to create new doc request: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if inputs.VendorProductId != -1 {
		err = database.LinkRequestToVendorProductWithTx(inputs.VendorProductId, request.Id, request.CatId, request.OrgId, role, tx)
		if err != nil {
			tx.Rollback()
			core.Warning("Failed to link request to vendor product: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		core.Warning("Failed to commit new doc request: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(request)
}

func updateDocumentRequest(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := UpdateDocumentRequestInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, inputs.OrgId)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	request := core.DocumentRequest{
		Id:              inputs.RequestId,
		Name:            inputs.Name,
		Description:     inputs.Description,
		CatId:           inputs.CatId,
		OrgId:           inputs.OrgId,
		RequestedUserId: inputs.RequestedUserId,
		RequestTime:     time.Now().UTC(),
	}

	err = database.UpdateDocumentRequest(&request, role)
	if err != nil {
		core.Warning("Failed to update doc request: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(request)
}

func getDocumentRequest(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetDocumentRequestInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, inputs.OrgId)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	req, err := database.GetDocumentRequest(inputs.RequestId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get single doc request: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cat, err := database.GetDocumentationCategory(req.CatId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get doc category: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fileIds, err := database.GetFulfilledFileIdsForDocRequest(inputs.RequestId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get relevant file IDs for request: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	files := make([]*core.ControlDocumentationFile, len(fileIds))
	for i, id := range fileIds {
		files[i], err = database.GetControlDocumentation(id, inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to get file metadata: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	jsonWriter.Encode(struct {
		Request  *core.DocumentRequest
		Files    []*core.ControlDocumentationFile
		Category *core.ControlDocumentationCategory
	}{
		Request:  req,
		Files:    files,
		Category: cat,
	})
}

func allDocumentRequests(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := AllDocumentRequestsInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, inputs.OrgId)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var reqs []*core.DocumentRequest
	if inputs.VendorProductId.NullInt64.Valid {
		reqs, err = database.GetAllDocumentRequestsForVendorProduct(inputs.VendorProductId.NullInt64.Int64, inputs.OrgId, role)
	} else if inputs.CatId.NullInt64.Valid {
		reqs, err = database.GetAllDocumentRequestsForDocCat(inputs.CatId.NullInt64.Int64, inputs.OrgId, role)
	} else {
		reqs, err = database.GetAllDocumentRequestsForOrganization(inputs.OrgId, role)
	}

	if err != nil {
		core.Warning("Failed to get all doc requests: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(reqs)
}

func deleteDocumentRequest(w http.ResponseWriter, r *http.Request) {
	inputs := DeleteDocumentRequestInputs{}
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

	err = database.DeleteDocumentRequest(inputs.RequestId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to delete doc request: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func completeDocumentRequest(w http.ResponseWriter, r *http.Request) {
	inputs := CompleteDocumentRequestInputs{}
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

	err = database.CompleteDocumentRequest(inputs.RequestId, inputs.OrgId, inputs.Complete, role)
	if err != nil {
		core.Warning("Failed to complete/reopen doc request: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
