package core

import (
	"fmt"
	"time"
)

type ControlDocumentationCategory struct {
	Id          int64  `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	OrgId       int32  `db:"org_id"`
}

type ControlDocumentationFile struct {
	Id           int64     `db:"id"`
	RelevantTime time.Time `db:"relevant_time"`
	CategoryId   int64     `db:"category_id"`
	OrgId        int32     `db:"org_id"`
	AltName      string    `db:"alt_name"`
	Description  string    `db:"description"`
}

type FileStorageData struct {
	Id           int64     `db:"id"`
	MetadataId   int64     `db:"metadata_id"`
	StorageName  string    `db:"storage_name"`
	OrgId        int32     `db:"org_id"`
	BucketId     string    `db:"bucket_id"`
	StorageId    string    `db:"storage_id"`
	UploadTime   time.Time `db:"upload_time"`
	UploadUserId int64     `db:"upload_user_id"`
}

type FileVersion struct {
	FileId        int64 `db:"file_id"`
	StorageId     int64 `db:"file_storage_id"`
	OrgId         int32 `db:"org_id"`
	VersionNumber int32 `db:"version_number"`
}

type FileStorageAuxData struct {
	IsPreview     bool  `db:"is_preview"`
	VersionNumber int32 `db:"version_number"`
}

type ControlDocumentationFileHandle struct {
	Id         int64
	CategoryId int64
}

func (file ControlDocumentationFile) UniqueKey() string {
	// It's tempting to also use the bucket id and storage id here but
	// the file ID is the only thing we control so it's probably safe only to use that.
	return fmt.Sprintf("controlDocFile-%d", file.Id)
}

func (file ControlDocumentationFile) StorageFilename(org *Organization, version int32) string {
	return fmt.Sprintf("org-%d-%s/%s-v%d", org.Id, org.Name, file.UniqueKey(), version)
}

type FileFolder struct {
	Id    int64  `db:"id"`
	OrgId int32  `db:"org_id"`
	Name  string `db:"name"`
}
