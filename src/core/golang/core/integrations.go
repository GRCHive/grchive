package core

type IntegrationType int32

const (
	ITSapErp IntegrationType = 1
)

type GenericIntegration struct {
	Id          int64           `db:"id"`
	Type        IntegrationType `db:"type"`
	OrgId       int32           `db:"org_id"`
	Name        string          `db:"name"`
	Description string          `db:"description"`
}
