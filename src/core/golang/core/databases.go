package core

import (
	"time"
)

type DatabaseType struct {
	Id            int32  `db:"id"`
	Name          string `db:"name"`
	HasSqlSupport bool   `db:"has_sql_support"`
}

type Database struct {
	Id        int64  `db:"id"`
	Name      string `db:"name"`
	OrgId     int32  `db:"org_id"`
	TypeId    int32  `db:"type_id"`
	OtherType string `db:"other_type"`
	Version   string `db:"version"`
}

type DatabaseConnection struct {
	Id         int64             `db:"id"`
	DbId       int64             `db:"db_id"`
	OrgId      int32             `db:"org_id"`
	Host       string            `db:"host"`
	Port       int32             `db:"port"`
	DbName     string            `db:"dbname"`
	Parameters map[string]string `db:"parameters"`
	Username   string            `db:"username"`
	Password   string            `db:"password" json:"-"`
	Salt       string            `db:"salt" json:"-"`
}

type DbRefresh struct {
	Id                int64    `db:"id"`
	DbId              int64    `db:"db_id"`
	OrgId             int32    `db:"org_id"`
	RefreshTime       NullTime `db:"refresh_time"`
	RefreshFinishTime NullTime `db:"refresh_finish_time"`
	RefreshSuccess    bool     `db:"refresh_success"`
	RefreshErrors     string   `db:"refresh_errors"`
}

type DbSchema struct {
	Id         int64  `db:"id"`
	OrgId      int32  `db:"org_id"`
	RefreshId  int64  `db:"refresh_id"`
	SchemaName string `db:"schema_name"`
}

type DbTable struct {
	Id        int64  `db:"id"`
	OrgId     int32  `db:"org_id"`
	SchemaId  int64  `db:"schema_id"`
	TableName string `db:"table_name"`
}

type DbColumn struct {
	Id         int64  `db:"id"`
	OrgId      int32  `db:"org_id"`
	TableId    int64  `db:"table_id"`
	ColumnName string `db:"column_name"`
	ColumnType string `db:"column_type"`
}

type DbFunction struct {
	Id       int64      `db:"id"`
	OrgId    int32      `db:"org_id"`
	SchemaId int64      `db:"schema_id"`
	Name     string     `db:"name"`
	Src      string     `db:"src"`
	RetType  NullString `db:"ret_type"`
}

type DbSqlQueryMetadata struct {
	Id          int64  `db:"id"`
	DbId        int64  `db:"db_id"`
	OrgId       int32  `db:"org_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

type DbSqlQuery struct {
	Id           int64     `db:"id"`
	MetadataId   int64     `db:"metadata_id"`
	Version      int32     `db:"version_number"`
	UploadTime   time.Time `db:"upload_time"`
	UploadUserId int64     `db:"upload_user_id"`
	OrgId        int32     `db:"org_id"`
	Query        string    `db:"query"`
}

type DbSqlQueryRequest struct {
	Id             int64     `db:"id"`
	QueryId        int64     `db:"query_id"`
	UploadTime     time.Time `db:"upload_time"`
	UploadUserId   int64     `db:"upload_user_id"`
	AssigneeUserId NullInt64 `db:"assignee"`
	DueDate        NullTime  `db:"due_date"`
	OrgId          int32     `db:"org_id"`
	Name           string    `db:"name"`
	Description    string    `db:"description"`
}

func (r *DbSqlQueryRequest) UnmarshalJSON(data []byte) error {
	return FlexibleJsonStructUnmarshal(data, r)
}

type DbSqlQueryRequestApproval struct {
	RequestId        int64     `db:"request_id"`
	OrgId            int32     `db:"org_id"`
	ResponseTime     time.Time `db:"response_time"`
	ResponsderUserId int64     `db:"responder_user_id"`
	Response         bool      `db:"response"`
	Reason           string    `db:"reason"`
}

type DbSqlQueryRunCode struct {
	RequestId      int64     `db:"request_id"`
	OrgId          int32     `db:"org_id"`
	ExpirationTime time.Time `db:"expiration_time"`
	UsedTime       NullTime  `db:"used_time"`
	HashedCode     string    `db:"hashed_code"`
	Salt           string    `db:"salt"`
}

type DatabaseFilterData struct {
	Type NumericFilterData
}

var NullDatabaseFilterData DatabaseFilterData = DatabaseFilterData{
	Type: NullNumericFilterData,
}

type DatabaseSettings struct {
	DbId               int64      `db:"db_id"`
	OrgId              int32      `db:"org_id"`
	AutoRefreshTaskId  NullInt64  `db:"auto_refresh_task" json:"-"`
	AutoRefreshEnabled bool       `db:"auto_refresh_enabled"`
	AutoRefreshRRule   NullString `db:"auto_refresh_rrule"`
}
