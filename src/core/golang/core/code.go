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

type CodeBuildStatus struct {
	Pending bool
	Success bool
}

type SupportedCodeParameterType struct {
	Id         int32  `db:"id"`
	Name       string `db:"name"`
	GolangType string `db:"golang_type" json:"-"`
	KotlinType string `db:"kotlin_type" json:"-"`
}

type CodeParameter struct {
	LinkId  int64  `db:"link_id" yaml:"-" json:"-"`
	Name    string `db:"name"`
	ParamId int32  `db:"param_type"`
}
