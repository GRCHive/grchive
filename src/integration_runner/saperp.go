package main

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/sap_api"
	"gitlab.com/grchive/grchive/vault_api"
	"gitlab.com/grchive/grchive/webcore"
)

func handleSapErpVersionWrapper(versionId int64, rfcId int64) error {
	data, err := handleSapErpVersion(versionId, rfcId)

	var success bool
	var storedData core.NullString

	if err != nil {
		success = false
		encryptedError, rerr := vault.TransitEncrypt(
			webcore.SapErpEncryptionPath,
			[]byte(err.Error()),
		)
		if rerr == nil {
			storedData = core.CreateNullString(string(encryptedError))
		}
	} else {
		success = true
		encryptedData, rerr := vault.TransitEncrypt(
			webcore.SapErpEncryptionPath,
			data,
		)
		if rerr == nil {
			storedData = core.CreateNullString(string(encryptedData))
		}
	}

	tx := database.CreateTx()
	return database.WrapTx(tx, func() error {
		return database.CompleteSapErpRfcVersionWithTx(tx, versionId, success, storedData)
	})
}

func handleSapErpVersion(versionId int64, rfcId int64) ([]byte, error) {
	rfc, err := database.GetSapErpRfc(rfcId)
	if err != nil {
		return nil, err
	}

	sapConnection, err := database.GetSapErpIntegration(rfc.IntegrationId)
	if err != nil {
		return nil, err
	}

	sapConnection.Password, err = webcore.DecryptEncryptedPassword(sapConnection.Password)
	if err != nil {
		return nil, err
	}

	if sapConnection.RealHostname.NullString.Valid {
		GlobalHostManager.AddOverride(sapConnection.Host, sapConnection.RealHostname.NullString.String)
	}
	defer func() {
		if sapConnection.RealHostname.NullString.Valid {
			GlobalHostManager.RemoveOverride(sapConnection.Host)
		}
	}()

	client, err := sap_api.CreateSapClient(*sapConnection)
	if err != nil {
		return nil, err
	}

	results, err := client.RunRfc(rfc.Function)
	if err != nil {
		return nil, err
	}

	rawJsonData, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}

	return rawJsonData, nil
}
