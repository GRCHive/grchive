package core

import (
	"fmt"
	"github.com/gosimple/slug"
)

type ClientData struct {
	Id          int64  `db:"id"`
	OrgId       int32  `db:"org_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

func (c ClientData) Filename(ext string) string {
	return fmt.Sprintf("src/main/kotlin/data/%s-%d.%s", slug.Make(c.Name), c.Id, ext)
}

type ClientDataVersion struct {
	Id      int64  `db:"id"`
	OrgId   int32  `db:"org_id"`
	DataId  int64  `db:"data_id"`
	Version int32  `db:"version"`
	Kotlin  string `db:"kotlin"`
}

type DataSourceOption struct {
	Id          int64  `db:"id"`
	Name        string `db:"name"`
	KotlinClass string `db:"kotlin_class"`
}

type DataSourceLink struct {
	OrgId        int32    `db:"org_id"`
	DataId       int64    `db:"data_id"`
	SourceId     SourceId `db:"source_id"`
	SourceTarget map[string]interface{}
}

type SourceId int64

const (
	SourceGrchive    SourceId = 1
	SourceDbPostgres          = 2
)

type FullClientDataWithLink struct {
	Data ClientData     `db:"data"`
	Link DataSourceLink `db:"link"`
}
