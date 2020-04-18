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
	UserId       int64     `db:"user_id"`
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
	RequiresBuild   bool       `db:"requires_build"`
	BuildStartTime  NullTime   `db:"build_start_time"`
	BuildFinishTime NullTime   `db:"build_finish_time"`
	BuildSuccess    bool       `db:"build_success"`
	RunStartTime    NullTime   `db:"run_start_time"`
	RunFinishTime   NullTime   `db:"run_finish_time"`
	RunSuccess      bool       `db:"run_success"`
	BuildLog        NullString `db:"build_log" json:"-"`
	RunLog          NullString `db:"run_log" json:"-"`
	UserId          int64      `db:"user_id"`
}

type DroneCiStatus struct {
	CodeId     int64     `db:"code_id"`
	OrgId      int32     `db:"org_id"`
	CommitHash string    `db:"commit_hash"`
	TimeStart  time.Time `db:"time_start"`
	TimeEnd    NullTime  `db:"time_end"`
	Success    bool      `db:"success"`
	Logs       string    `db:"logs" json:"-"`
	Jar        string    `db:"jar" json:"-"`
}
