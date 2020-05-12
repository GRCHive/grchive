package core

import (
	"fmt"
	"github.com/gosimple/slug"
	"github.com/iancoleman/strcase"
)

type ClientScript struct {
	Id          int64  `db:"id"`
	OrgId       int32  `db:"org_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

func (s ClientScript) Filename(ext string) string {
	return fmt.Sprintf("src/main/kotlin/scripts/%s-%d.%s", strcase.ToCamel(s.Name), s.Id, ext)
}

func (s ClientScript) MetadataFilename() string {
	return fmt.Sprintf("src/main/resources/scripts/metadata-%s-%d.yaml", slug.Make(s.Name), s.Id)
}
