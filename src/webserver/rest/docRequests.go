package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
	"time"
)

type NewDocumentRequestInputs struct {
	Name            string         `json:"name"`
	Description     string         `json:"description"`
	CatId           int64          `json:"catId"`
	ControlId       core.NullInt64 `json:"controlId"`
	FolderId        core.NullInt64 `json:"folderId"`
	OrgId           int32          `json:"orgId"`
	RequestedUserId int64          `json:"requestedUserId"`
	VendorProductId int64          `json:"vendorProductId"`
	AssigneeUserId  core.NullInt64 `json:"assigneeUserId"`
	DueDate         core.NullTime  `json:"dueDate"`
}

type UpdateDocumentRequestInputs struct {
	RequestId       int64          `json:"requestId"`
	Name            string         `json:"name"`
	Description     string         `json:"description"`
	CatId           int64          `json:"catId"`
	OrgId           int32          `json:"orgId"`
	RequestedUserId int64          `json:"requestedUserId"`
	AssigneeUserId  core.NullInt64 `json:"assigneeUserId"`
	DueDate         core.NullTime  `json:"dueDate"`
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
		OrgId:           inputs.OrgId,
		RequestedUserId: inputs.RequestedUserId,
		RequestTime:     time.Now().UTC(),
		AssigneeUserId:  inputs.AssigneeUserId,
		DueDate:         inputs.DueDate,
	}

	tx, err := database.CreateAuditTrailTx(role)
	if err != nil {
		core.Warning("Failed to create audit trail TX: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = database.WrapTx(tx, func() error {
		return database.CreateNewDocumentRequestWithTx(&request, role, tx)
	}, func() error {
		// We either need to be linking to a vendor or linking to a control/folder.
		if inputs.VendorProductId != -1 {
			return database.LinkRequestToVendorProductWithTx(inputs.VendorProductId, request.Id, request.OrgId, role, tx)
		} else if inputs.ControlId.NullInt64.Valid && inputs.FolderId.NullInt64.Valid {
			err := database.AddDocRequestControlLinkWithTx(request.Id, inputs.ControlId.NullInt64.Int64, inputs.OrgId, role, tx)
			if err != nil {
				return err
			}
			return database.AddDocRequestControlFolderLinkWithTx(
				request.Id,
				inputs.ControlId.NullInt64.Int64,
				inputs.FolderId.NullInt64.Int64,
				inputs.OrgId,
				role,
				tx,
			)
		}
		return nil
	}, func() error {
		return database.AddDocRequestDocCatLinkWithTx(request.Id, inputs.CatId, inputs.OrgId, role, tx)
	})

	if err != nil {
		core.Warning("Failed to create new doc request: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(struct {
		Request *core.DocumentRequest
	}{
		Request: &request,
	})
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
		OrgId:           inputs.OrgId,
		RequestedUserId: inputs.RequestedUserId,
		RequestTime:     time.Now().UTC(),
		AssigneeUserId:  inputs.AssigneeUserId,
		DueDate:         inputs.DueDate,
	}

	err = database.UpdateDocumentRequest(&request, role)
	if err != nil {
		core.Warning("Failed to update doc request: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(struct {
		Request *core.DocumentRequest
	}{
		Request: &request,
	})
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

	vendorId, vendorProductId, err := database.GetVendorProductIdForDocRequest(inputs.RequestId, inputs.OrgId)
	if err != nil {
		core.Warning("Failed to get vendor product id: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(struct {
		Request         *core.DocumentRequest
		Files           []*core.ControlDocumentationFile
		VendorProductId int64
		VendorId        int64
	}{
		Request:         req,
		Files:           files,
		VendorProductId: vendorProductId,
		VendorId:        vendorId,
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
		reqs, err = database.FindDocRequestsLinkedToDocCat(inputs.CatId.NullInt64.Int64, inputs.OrgId, role)
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

type NewDocumentRequestFileLinksInputs struct {
	Files []int64 `json:"files"`
}

func newDocRequestFileLinks(w http.ResponseWriter, r *http.Request) {
	inputs := NewDocumentRequestFileLinksInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	request, err := webcore.FindDocumentRequestInContext(r.Context())
	if err != nil {
		core.Warning("Can't find doc request: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tx := database.CreateTx()
	err = database.WrapTx(tx, func() error {
		for _, fid := range inputs.Files {
			err := database.FulfillDocumentRequestWithTx(request.Id, fid, request.OrgId, core.ServerRole, tx)
			if err != nil {
				return err
			}
		}
		return nil
	}, func() error {
		return database.MarkDocumentRequestProgressWithTx(request.Id, request.OrgId, tx)
	})

	if err != nil {
		core.Warning("Failed to link files to request: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
