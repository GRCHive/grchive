package core

import (
	"time"
)

type ShellTypeId int32

const (
	BashShellId  ShellTypeId = 1
	PowerShellId             = 2
)

type ShellScript struct {
	Id          int64       `db:"id"`
	OrgId       int32       `db:"org_id"`
	TypeId      ShellTypeId `db:"type_id"`
	Name        string      `db:"name"`
	Description string      `db:"description"`
	BucketId    string      `db:"bucket_id"`
	StorageId   string      `db:"storage_id"`
}

type ShellScriptVersion struct {
	Id            int64     `db:"id"`
	ShellId       int64     `db:"shell_id"`
	OrgId         int32     `db:"org_id"`
	UploadTime    time.Time `db:"upload_time"`
	UploadUserId  int64     `db:"upload_user_id"`
	GcsGeneration int64     `db:"gcs_generation"`
}

type ShellScriptRun struct {
	Id              int64     `db:"id"`
	ScriptVersionId int64     `db:"script_version_id"`
	RunUserId       int64     `db:"run_user_id"`
	CreateTime      time.Time `db:"create_time"`
	RunTime         NullTime  `db:"run_time"`
	EndTime         NullTime  `db:"end_time"`
}

type ShellScriptRunPerServer struct {
	RunId        int64      `db:"run_id"`
	OrgId        int32      `db:"org_id"`
	ServerId     int64      `db:"server_id"`
	EncryptedLog NullString `db:"encrypted_log"`
	RunTime      NullTime   `db:"run_time"`
	EndTime      NullTime   `db:"end_time"`
	Success      bool       `db:"success"`
}
