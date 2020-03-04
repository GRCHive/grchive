package database

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
)

func NewVendorWithTx(vendor *core.Vendor, role *core.Role, tx *sqlx.Tx) error {
	if !role.Permissions.HasAccess(core.ResourceVendors, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	err := UpgradeTxToAudit(tx, role)
	if err != nil {
		return err
	}

	rows, err := tx.NamedQuery(`
		INSERT INTO vendors(name, description, url, org_id)
		VALUES (:name, :description, :url, :org_id)
		RETURNING id
	`, vendor)
	if err != nil {
		return err
	}

	rows.Next()
	err = rows.Scan(&vendor.Id)
	if err != nil {
		rows.Close()
		return err
	}
	rows.Close()

	_, err = tx.Exec(`
		INSERT INTO vendor_documentation_category_link (vendor_id, org_id, doc_cat_id)
		VALUES ($1, $2, $3)
	`, vendor.Id, vendor.OrgId, vendor.DocCatId)
	return err
}

func AllVendorsForOrganization(orgId int32, role *core.Role) ([]*core.Vendor, error) {
	if !role.Permissions.HasAccess(core.ResourceVendors, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	vendors := make([]*core.Vendor, 0)
	err := dbConn.Select(&vendors, `
		SELECT vnd.*, link.doc_cat_id
		FROM vendors AS vnd
		INNER JOIN vendor_documentation_category_link AS link
			ON vnd.id = link.vendor_id
		WHERE vnd.org_id = $1
	`, orgId)
	return vendors, err
}

func GetVendorFromId(vendorId int64, orgId int32, role *core.Role) (*core.Vendor, error) {
	if !role.Permissions.HasAccess(core.ResourceVendors, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	vendor := core.Vendor{}
	err := dbConn.Get(&vendor, `
		SELECT vnd.*, link.doc_cat_id
		FROM vendors AS vnd
		INNER JOIN vendor_documentation_category_link AS link
			ON vnd.id = link.vendor_id
		WHERE vnd.id = $1 AND vnd.org_id = $2
	`, vendorId, orgId)
	return &vendor, err
}

func DeleteVendor(vendorId int64, orgId int32, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceVendors, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	// First delete the relevant documentation category and then
	// delete the vendor itself. Everything else should propgate through
	// CASCADE.
	catId, err := GetDocCatIdForVendor(vendorId, orgId, role)

	if err != nil {
		tx.Rollback()
		return err
	}

	err = DeleteControlDocumentationCategoryWithTx(catId, orgId, role, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM vendors
		WHERE id = $1 AND org_id = $2
	`, vendorId, orgId)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func UpdateVendor(vendor *core.Vendor, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceVendors, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(`
		UPDATE vendors
		SET name = :name,
			description = :description,
			url = :url
		WHERE id = :id AND org_id = :org_id
	`, vendor)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func GetDocCatIdForVendor(vendorId int64, orgId int32, role *core.Role) (int64, error) {
	if !role.Permissions.HasAccess(core.ResourceVendors, core.AccessView) {
		return -1, core.ErrorUnauthorized
	}

	catId := int64(0)
	err := dbConn.Get(&catId, `
		SELECT doc_cat_id
		FROM vendor_documentation_category_link
		WHERE vendor_id = $1 AND org_id = $2
	`, vendorId, orgId)
	return catId, err
}
