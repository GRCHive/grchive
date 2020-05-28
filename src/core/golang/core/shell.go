package core

import (
	"time"
)

type ShellScript struct {
	Id          int64  `db:"id"`
	OrgId       int32  `db:"org_id"`
	TypeId      int32  `db:"type_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	BucketId    string `db:"bucket_id"`
	StorageId   string `db:"storage_id"`
}

type ShellScriptVersion struct {
	ShellId       int64     `db:"shell_id"`
	OrgId         int32     `db:"org_id"`
	UploadTime    time.Time `db:"upload_time"`
	UploadUserId  int64     `db:"upload_user_id"`
	GcsGeneration int64     `db:"gcs_generation"`
}
