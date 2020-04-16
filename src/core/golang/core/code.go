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
	Name    string `db:"name" yaml:"name"`
	ParamId int32  `db:"param_type" yaml:"paramId"`
}

type ScriptRun struct {
	Id              int64      `db:"id"`
	LinkId          int64      `db:"link_id"`
	StartTime       time.Time  `db:"start_time"`
	BuildStartTime  NullTime   `db:"build_start_time"`
	BuildFinishTime NullTime   `db:"build_finish_time"`
	BuildSuccess    bool       `db:"build_success"`
	RunStartTime    NullTime   `db:"run_start_time"`
	RunFinishTime   NullTime   `db:"run_finish_time"`
	RunSuccess      bool       `db:"run_success"`
	BuildLog        NullString `db:"build_log"`
	RunLog          NullString `db:"run_log"`
	UserId          int64      `db:"user_id"`
}
