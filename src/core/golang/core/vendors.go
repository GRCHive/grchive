package core

import (
	"fmt"
)

type Vendor struct {
	Id          int64  `db:"id"`
	OrgId       int32  `db:"org_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Url         string `db:"url"`
	DocCatId    int64  `db:"doc_cat_id"`
}

type VendorProduct struct {
	Id          int64  `db:"id"`
	VendorId    int64  `db:"vendor_id"`
	OrgId       int32  `db:"org_id"`
	Name        string `db:"product_name"`
	Description string `db:"description"`
	Url         string `db:"url"`
}

func (v Vendor) CreateDocumentationCategory() ControlDocumentationCategory {
	return ControlDocumentationCategory{
		Name:        fmt.Sprintf("%s Files", v.Name),
		Description: "",
		OrgId:       v.OrgId,
	}
}
