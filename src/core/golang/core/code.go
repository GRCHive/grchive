package core

import (
	"time"
)

type ManagedCode struct {
	Id           int64     `db:"id"`
	OrgId        int32     `db:"org_id"`
	GitHash      string    `db:"git_hash"`
	ActionTime   time.Time `db:"action_time"`
	GitPath      string    `db:"git_path"`
	GiteaFileSha string    `db:"gitea_file_sha"`
}
