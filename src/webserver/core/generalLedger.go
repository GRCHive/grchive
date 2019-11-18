package core

type GeneralLedgerAccount struct {
	Id                  int64  `db:"id"`
	OrgId               int32  `db:"org_id"`
	ParentCategoryId    int64  `db:"parent_category_id"`
	AccountId           string `db:"account_identifier"`
	AccountName         string `db:"account_name"`
	AccountDescription  string `db:"account_description"`
	FinanciallyRelevant bool   `db:"financially_relevant"`
}

type GeneralLedgerCategory struct {
	Id               int64     `db:"id"`
	OrgId            int32     `db:"org_id"`
	ParentCategoryId NullInt64 `db:"parent_category_id"`
	Name             string    `db:"name"`
	Description      string    `db:"description"`
}
