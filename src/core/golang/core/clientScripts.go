package core

import (
	"fmt"
	"github.com/gosimple/slug"
)

type ClientScript struct {
	Id          int64  `db:"id"`
	OrgId       int32  `db:"org_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

func (s ClientScript) Filename(ext string) string {
	return fmt.Sprintf("%s-%d.%s", slug.Make(s.Name), s.Id, ext)
}
