package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/sap_api"
	"gitlab.com/grchive/grchive/vault_api"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
)

type NewSapErpInput struct {
	Integration core.GenericIntegration         `json:"integration"`
	Setup       sap_api.SapRfcConnectionOptions `json:"setup"`
}

func newSapErpIntegration(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := NewSapErpInput{}
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

	system, err := webcore.FindSystemInContext(r.Context())
	if err != nil {
		core.Warning("Can't find system in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inputs.Integration.OrgId = org.Id
	inputs.Integration.Type = core.ITSapErp

	inputs.Setup.Password, err = webcore.CreateEncryptedPassword(inputs.Setup.Password)
	if err != nil {
		core.Warning("Failed to encrypted password: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tx := database.CreateTx()
	err = database.WrapTx(tx, func() error {
		return database.NewIntegrationWithTx(tx, &inputs.Integration)
	}, func() error {
		return database.NewSapErpIntegrationWithTx(tx, inputs.Integration.Id, inputs.Setup)
	}, func() error {
		return database.LinkIntegrationWithSystemWithTx(tx, inputs.Integration.Id, system.Id, system.OrgId)
	})

	if err != nil {
		core.Warning("Failed to create integration: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(inputs.Integration)
}

type EditSapErpInput struct {
	Setup sap_api.SapRfcConnectionOptions `json:"setup"`
}

func editSapErpIntegration(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	integration, err := webcore.FindGenericIntegrationInContext(r.Context())
	if err != nil {
		core.Warning("Can't find generic integration in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inputs := EditSapErpInput{}
	err = webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	retPassword := inputs.Setup.Password
	inputs.Setup.Password, err = webcore.CreateEncryptedPassword(inputs.Setup.Password)
	if err != nil {
		core.Warning("Failed to encrypt password: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tx := database.CreateTx()
	err = database.WrapTx(tx, func() error {
		return database.EditSapErpIntegrationWithTx(tx, integration.Id, inputs.Setup)
	})
	if err != nil {
		core.Warning("Failed to edit SAP ERP integration: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	inputs.Setup.Password = retPassword
	jsonWriter.Encode(inputs.Setup)
}

func getSapErpIntegration(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	integration, err := webcore.FindGenericIntegrationInContext(r.Context())
	if err != nil {
		core.Warning("Can't find generic integration in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sap, err := database.GetSapErpIntegration(integration.Id)
	if err != nil {
		core.Warning("Can't get SAP ERP integration: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sap.Password, err = webcore.DecryptEncryptedPassword(sap.Password)
	if err != nil {
		core.Warning("Failed to decrypt password: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(sap)
}

func allSapErpRfc(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	integration, err := webcore.FindGenericIntegrationInContext(r.Context())
	if err != nil {
		core.Warning("Can't find generic integration in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rfc, err := database.AllSapErpRfc(integration.Id)
	if err != nil {
		core.Warning("Can't get all SAP ERP RFC: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(rfc)
}

type NewSapErpRfcInput struct {
	Function string `json:"function"`
}

func newSapErpRfc(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := NewSapErpRfcInput{}
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

	rfc := core.SapErpRfc{
		IntegrationId: integration.Id,
		Function:      inputs.Function,
	}

	var version *core.SapErpRfcVersion
	tx := database.CreateTx()
	err = database.WrapTx(tx, func() error {
		return database.NewSapErpRfcWithTx(tx, &rfc)
	}, func() error {
		var ierr error
		version, ierr = database.NewSapErpRfcVersionWithTx(tx, rfc.Id)
		return ierr
	})

	if err != nil {
		core.Warning("Failed to create SAP ERP RFC: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	webcore.DefaultRabbitMQ.SendMessage(webcore.PublishMessage{
		Exchange: webcore.DEFAULT_EXCHANGE,
		Queue:    webcore.SAP_ERP_RFC_QUEUE,
		Body: webcore.SapErpRfcMessage{
			RfcId:     rfc.Id,
			VersionId: version.Id,
		},
	})

	jsonWriter.Encode(rfc)
}

func deleteSapErpRfc(w http.ResponseWriter, r *http.Request) {
	rfc, err := webcore.FindSapErpRfcInContext(r.Context())
	if err != nil {
		core.Warning("Can't find SAP ERP RFC in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tx := database.CreateTx()
	err = database.WrapTx(tx, func() error {
		return database.DeleteSapErpRfcWithTx(tx, rfc.Id)
	})

	if err != nil {
		core.Warning("Failed to delete SAP ERP RFC: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func allSapErpRfcVersions(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	rfc, err := webcore.FindSapErpRfcInContext(r.Context())
	if err != nil {
		core.Warning("Can't find SAP ERP RFC in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	versions, err := database.AllSapErpRfcVersions(rfc.Id)
	if err != nil {
		core.Warning("Failed to get SAP ERP RFC versions: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(versions)
}

func newSapErpRfcVersion(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	rfc, err := webcore.FindSapErpRfcInContext(r.Context())
	if err != nil {
		core.Warning("Can't find SAP ERP RFC in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var version *core.SapErpRfcVersion
	tx := database.CreateTx()
	err = database.WrapTx(tx, func() error {
		var ierr error
		version, ierr = database.NewSapErpRfcVersionWithTx(tx, rfc.Id)
		return ierr
	})

	if err != nil {
		core.Warning("Failed to create new SAP ERP RFC version: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	webcore.DefaultRabbitMQ.SendMessage(webcore.PublishMessage{
		Exchange: webcore.DEFAULT_EXCHANGE,
		Queue:    webcore.SAP_ERP_RFC_QUEUE,
		Body: webcore.SapErpRfcMessage{
			RfcId:     rfc.Id,
			VersionId: version.Id,
		},
	})

	jsonWriter.Encode(version)
}

func getSapErpRfcVersion(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	version, err := webcore.FindSapErpRfcVersionInContext(r.Context())
	if err != nil {
		core.Warning("Can't find SAP ERP RFC in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if version.RawData.NullString.Valid {
		decryptedData, err := vault.TransitDecrypt(
			webcore.SapErpEncryptionPath,
			[]byte(version.RawData.NullString.String),
		)
		if err != nil {
			core.Warning("Failed to decrypt RFC data: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if version.Success {
			version.Data = &json.RawMessage{}
			err = json.Unmarshal(decryptedData, &version.Data)
			if err != nil {
				core.Warning("Failed to parse JSON: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			version.Logs = core.CreateNullString(string(decryptedData))
		}
	}

	jsonWriter.Encode(version)
}
