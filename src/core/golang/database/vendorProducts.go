package database

import (
	"gitlab.com/grchive/grchive/core"
	"strconv"
)

func NewVendorProduct(product *core.VendorProduct, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceVendors, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	rows, err := tx.NamedQuery(`
		INSERT INTO vendor_products(product_name, description, url, org_id, vendor_id)
		VALUES (:product_name, :description, :url, :org_id, :vendor_id)
		RETURNING id
	`, product)
	if err != nil {
		tx.Rollback()
		return err
	}

	rows.Next()
	err = rows.Scan(&product.Id)
	if err != nil {
		rows.Close()
		tx.Rollback()
		return err
	}
	rows.Close()

	return tx.Commit()
}

func AllVendorProductsForVendor(vendorId int64, orgId int32, role *core.Role) ([]*core.VendorProduct, error) {
	if !role.Permissions.HasAccess(core.ResourceVendors, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	products := make([]*core.VendorProduct, 0)
	err := dbConn.Select(&products, `
		SELECT *
		FROM vendor_products
		WHERE org_id = $1 and vendor_id = $2
	`, orgId, vendorId)

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return nil, err
	}

	for _, p := range products {
		err = LogAuditSelectWithTx(orgId, core.ResourceVendorProduct, strconv.FormatInt(p.Id, 10), role, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	return products, tx.Commit()
}

func GetVendorProduct(productId int64, vendorId int64, orgId int32, role *core.Role) (*core.VendorProduct, error) {
	if !role.Permissions.HasAccess(core.ResourceVendors, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	product := core.VendorProduct{}
	err := dbConn.Get(&product, `
		SELECT *
		FROM vendor_products
		WHERE id = $1 AND vendor_id = $2 AND org_id = $3
	`, productId, vendorId, orgId)
	if err != nil {
		return nil, err
	}
	return &product, LogAuditSelect(orgId, core.ResourceVendorProduct, strconv.FormatInt(product.Id, 10), role)
}

func UpdateVendorProduct(product *core.VendorProduct, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceVendors, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(`
		UPDATE vendor_products
		SET product_name = :product_name,
			description = :description,
			url = :url
		WHERE id = :id
	`, product)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func DeleteVendorProduct(productId int64, vendorId int64, orgId int32, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceVendors, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM vendor_products
		WHERE id = $1 AND vendor_id = $2 AND org_id = $3
	`, productId, vendorId, orgId)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func LinkVendorProductToSocFiles(productId int64, orgId int32, files []*core.ControlDocumentationFileHandle, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceVendors, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()

	for _, file := range files {
		_, err := tx.Exec(`
			INSERT INTO vendor_product_soc_reports (product_id, org_id, file_id)
			VALUES ($1, $2, $3)
		`, productId, orgId, file.Id)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func UnlinkVendorProductFromSocFiles(productId int64, orgId int32, files []*core.ControlDocumentationFileHandle, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceVendors, core.AccessEdit) {
		return core.ErrorUnauthorized
	}
	tx := dbConn.MustBegin()

	for _, file := range files {
		_, err := tx.Exec(`
			DELETE FROM vendor_product_soc_reports
			WHERE product_id = $1
				AND org_id = $2
				AND file_id = $3
				AND cat_id = $4
		`, productId, orgId, file.Id, file.CategoryId)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
