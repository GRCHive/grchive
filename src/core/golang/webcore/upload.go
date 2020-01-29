package webcore

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/backblaze_api"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/vault_api"
	"time"
)

func UploadNewFileWithTx(
	file *core.ControlDocumentationFile,
	fileName string,
	buffer []byte,
	role *core.Role,
	org *core.Organization,
	uploadUserId int64,
	b2Auth *backblaze.B2AuthToken,
	tx *sqlx.Tx,
) (*backblaze.B2File, error) {
	err := database.CreateControlDocumentationFileWithTx(file, tx, role)
	if err != nil {
		return nil, err
	}

	transitKey := file.UniqueKey()
	err = vault.TransitCreateNewEngineKey(transitKey)
	if err != nil {
		return nil, err
	}

	encryptedFile, err := vault.TransitEncrypt(transitKey, buffer)
	if err != nil {
		return nil, err
	}

	b2Filename := file.StorageFilename(org)

	b2File, err := backblaze.UploadFile(b2Auth,
		core.EnvConfig.Backblaze.ControlDocBucketId,
		b2Filename,
		encryptedFile)
	if err != nil {
		return nil, err
	}

	storage := core.FileStorageData{
		MetadataId:   file.Id,
		StorageName:  fileName,
		OrgId:        file.OrgId,
		BucketId:     b2File.BucketId,
		StorageId:    b2File.FileId,
		UploadTime:   time.Now().UTC(),
		UploadUserId: uploadUserId,
	}

	err = database.CreateFileStorageWithTx(&storage, tx, role)
	if err != nil {
		backblaze.DeleteFile(b2Auth, file.StorageFilename(org), b2File)
		return nil, err
	}

	err = database.AddFileVersionWithTx(file, &storage, tx, role)
	if err != nil {
		backblaze.DeleteFile(b2Auth, file.StorageFilename(org), b2File)
		return nil, err
	}

	return &b2File, nil
}
