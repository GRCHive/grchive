package webcore

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/gcloud_api/storage"
	"gitlab.com/grchive/grchive/vault_api"
	"time"
)

func UploadNewFileWithTx(
	gcloud storage.GCloudStorageApi,
	file *core.ControlDocumentationFile,
	bucket string,
	storageName *string,
	fileName string,
	buffer []byte,
	role *core.Role,
	uploadUserId int64,
	tx *sqlx.Tx,
	org *core.Organization,
	useExistingMetadata bool,
	addToFileVersion bool,
) (*core.FileStorageData, error) {
	var err error

	if !useExistingMetadata {
		err = database.CreateControlDocumentationFileWithTx(file, tx, role)
		if err != nil {
			return nil, err
		}
	}

	if storageName == nil {
		tmp := file.StorageFilename(org)
		storageName = &tmp
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

	err = gcloud.Upload(bucket, *storageName, encryptedFile)
	if err != nil {
		return nil, err
	}

	storage := core.FileStorageData{
		MetadataId:   file.Id,
		StorageName:  fileName,
		OrgId:        file.OrgId,
		BucketId:     bucket,
		StorageId:    *storageName,
		UploadTime:   time.Now().UTC(),
		UploadUserId: uploadUserId,
	}

	err = database.CreateFileStorageWithTx(&storage, tx, role)
	if err != nil {
		gcloud.Delete(bucket, *storageName)
		return nil, err
	}

	if addToFileVersion {
		err = database.AddFileVersionWithTx(file, &storage, tx, role)
		if err != nil {
			gcloud.Delete(bucket, *storageName)
			return nil, err
		}
	}

	return &storage, nil
}
