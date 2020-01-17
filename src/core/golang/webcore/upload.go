package webcore

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/b3h47pte/audit-stuff/backblaze_api"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"gitlab.com/b3h47pte/audit-stuff/vault_api"
)

func UploadNewFileWithTx(file *core.ControlDocumentationFile, buffer []byte, role *core.Role, org *core.Organization, b2Auth *backblaze.B2AuthToken, tx *sqlx.Tx) error {
	err := database.CreateControlDocumentationFileWithTx(file, tx, role)
	if err != nil {
		return err
	}

	transitKey := file.UniqueKey()
	err = vault.TransitCreateNewEngineKey(transitKey)
	if err != nil {
		return err
	}

	encryptedFile, err := vault.TransitEncrypt(transitKey, buffer)
	if err != nil {
		return err
	}

	b2Filename := file.StorageFilename(org)

	b2File, err := backblaze.UploadFile(b2Auth,
		core.EnvConfig.Backblaze.ControlDocBucketId,
		b2Filename,
		encryptedFile)
	if err != nil {
		return err
	}

	file.BucketId = b2File.BucketId
	file.StorageId = b2File.FileId
	err = database.UpdateControlDocumentation(file, tx, role)
	if err != nil {
		backblaze.DeleteFile(b2Auth, file.StorageFilename(org), b2File)
		return err
	}

	return nil
}
