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
