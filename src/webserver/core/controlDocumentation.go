package core

import (
	"fmt"
	"time"
)

type ControlDocumentationCategory struct {
	Id          int64  `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	ControlId   int64  `db:"control_id"`
}

type ControlDocumentationFile struct {
	Id           int64     `db:"id"`
	BucketId     string    `db:"bucket_id" json:"-"`
	StorageId    string    `db:"storage_id" json:"-"`
	StorageName  string    `db:"storage_name"`
	RelevantTime time.Time `db:"relevant_time"`
	UploadTime   time.Time `db:"upload_time"`
	CategoryId   int64     `db:"category_id"`
}

func (file ControlDocumentationFile) UniqueKey() string {
	// It's tempting to also use the bucket id and storage id here but
	// the file ID is the only thing we control so it's probably safe only to use that.
	return fmt.Sprintf("controlDocFile-%d", file.Id)
}
