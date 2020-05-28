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
) (string, *core.FileStorageData, error) {
	var err error

	core.Debug("\tCreate Metadata")
	if !useExistingMetadata {
		err = database.CreateControlDocumentationFileWithTx(file, tx, role)
		if err != nil {
			return "", nil, err
		}
	}

	// This little sequence is kind shitty since
	// we want to put in the version number into the filename
	// but we aren't able to create a new file version until
	// the storage object is created as well. So we need to
	// store the FileStorageData in two steps. The first step
	// we put in a temporary StorageId which will get updated
	// after creating the FileVersion object.
	storage := core.FileStorageData{
		MetadataId:   file.Id,
		StorageName:  fileName,
		OrgId:        file.OrgId,
		BucketId:     bucket,
		StorageId:    "",
		UploadTime:   time.Now().UTC(),
		UploadUserId: uploadUserId,
	}

	core.Debug("\tCreate Storage")
	err = database.CreateFileStorageWithTx(&storage, tx, role)
	if err != nil {
		return "", nil, err
	}

	var latestVersion *core.FileVersion
	if addToFileVersion {
		latestVersion, err = database.AddFileVersionWithTx(file, &storage, tx, role)
	} else {
		latestVersion, err = database.GetLatestNonPreviewFileVersion(file.Id, file.OrgId, role)
	}
	if err != nil {
		return "", nil, err
	}

	if storageName == nil {
		tmp := file.StorageFilename(org, latestVersion.VersionNumber)
		storageName = &tmp
	}

	core.Debug("\tUpdate Database")
	storage.StorageId = *storageName
	err = database.UpdateFileStorageStorageIdWithTx(storage.Id, file.OrgId, storage.StorageId, tx, role)
	if err != nil {
		return "", nil, err
	}

	core.Debug("\tCreate Engine Key")
	transitKey := file.UniqueKey()
	err = vault.TransitCreateNewEngineKey(transitKey)
	if err != nil {
		return "", nil, err
	}

	core.Debug("\tEncrypt")
	encryptedFile, err := vault.TransitEncrypt(transitKey, buffer)
	if err != nil {
		return "", nil, err
	}

	core.Debug("\tUpload")
	_, err = gcloud.Upload(bucket, *storageName, encryptedFile, core.EnvConfig.HmacKey)
	if err != nil {
		return "", nil, err
	}

	return *storageName, &storage, nil
}
