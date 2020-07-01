package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
)

func readPbcNotificationCadenceSettingsHelper(constraint string, args ...interface{}) ([]*core.PbcNotificationCadenceSettings, error) {
	rows, err := dbConn.Queryx(fmt.Sprintf(`
		SELECT cs.*, ARRAY_TO_JSON(ARRAY_REMOVE(ARRAY_AGG(cau.user_id), NULL)) AS "additional_users"
		FROM org_pbc_notification_cadence_settings AS cs
		LEFT JOIN org_pbc_notification_cadence_addtl_users AS cau
			ON cau.cadence_id = cs.id
		%s
		GROUP BY 
			cs.id,
			cs.org_id,
			cs.days_before_due,
			cs.send_to_assignee,
			cs.send_to_requester
		ORDER BY cs.days_before_due DESC
	`, constraint), args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ret := make([]*core.PbcNotificationCadenceSettings, 0)
	for rows.Next() {
		data := map[string]interface{}{}
		err := rows.MapScan(data)
		if err != nil {
			return nil, err
		}

		newSetting := core.PbcNotificationCadenceSettings{}
		newSetting.Id = data["id"].(int64)
		newSetting.OrgId = int32(data["org_id"].(int64))
		newSetting.DaysBeforeDue = int32(data["days_before_due"].(int64))
		newSetting.SendToAssignee = data["send_to_assignee"].(bool)
		newSetting.SendToRequester = data["send_to_requester"].(bool)
		newSetting.AdditionalUsers, err = readInt64Array(data["additional_users"].([]uint8))
		if err != nil {
			return nil, err
		}

		ret = append(ret, &newSetting)
	}
	return ret, nil
}

func GetOrgPbcNotificationCadenceSettings(orgId int32) ([]*core.PbcNotificationCadenceSettings, error) {
	return readPbcNotificationCadenceSettingsHelper("WHERE org_id = $1", orgId)
}

func GetPbcNotificationCadenceSettings(id int64) (*core.PbcNotificationCadenceSettings, error) {
	data, err := readPbcNotificationCadenceSettingsHelper("WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return data[0], nil
}

func DeletePbcNotificationCadenceSettingsWithTx(tx *sqlx.Tx, id int64) error {
	_, err := tx.Exec(`
		DELETE FROM org_pbc_notification_cadence_settings
		WHERE id = $1
	`, id)
	return err
}

func CreateNewPbcNotificationCadenceSettingWithTx(tx *sqlx.Tx, orgId int32, daysBefore int32) (*core.PbcNotificationCadenceSettings, error) {
	rows, err := tx.Queryx(`
		INSERT INTO org_pbc_notification_cadence_settings (org_id, days_before_due)
		VALUES ($1, $2)
		RETURNING *
	`, orgId, daysBefore)

	if err != nil {
		return nil, nil
	}

	defer rows.Close()
	rows.Next()

	settings := core.PbcNotificationCadenceSettings{}
	err = rows.StructScan(&settings)
	settings.AdditionalUsers = make([]int64, 0)
	return &settings, err
}

func EditPbcNotificationCadenceSettingWithTx(tx *sqlx.Tx, setting core.PbcNotificationCadenceSettings) error {
	_, err := tx.NamedExec(`
		UPDATE org_pbc_notification_cadence_settings
		SET days_before_due = :days_before_due,
			send_to_assignee = :send_to_assignee,	
			send_to_requester = :send_to_requester
		WHERE org_id = :org_id AND id = :id
	`, setting)

	if err != nil {
		return nil
	}

	// Easiest way to write this - need to sync the values found in multiple rows with an array coming from the user.
	// So just delete -> reinsert.
	_, err = tx.NamedExec(`
		DELETE FROM org_pbc_notification_cadence_addtl_users
		WHERE cadence_id = :id
	`, setting)

	if err != nil {
		return err
	}

	for _, uid := range setting.AdditionalUsers {
		_, err = tx.Exec(`
			INSERT INTO org_pbc_notification_cadence_addtl_users (cadence_id, user_id)
			VALUES ($1, $2)
		`, setting.Id, uid)

		if err != nil {
			return err
		}
	}

	return nil
}

func EditAllPbcNotificationCadenceSettingWithTx(tx *sqlx.Tx, setting core.PbcNotificationCadenceSettings) error {
	_, err := tx.NamedExec(`
		UPDATE org_pbc_notification_cadence_settings
		SET send_to_assignee = :send_to_assignee,	
			send_to_requester = :send_to_requester
		WHERE org_id = :org_id
	`, setting)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM org_pbc_notification_cadence_addtl_users AS cau
		USING org_pbc_notification_cadence_settings AS cs
		WHERE cau.cadence_id = cs.id
			AND cs.org_id = $1
	`, setting.OrgId)

	if err != nil {
		return err
	}

	for _, uid := range setting.AdditionalUsers {
		_, err = tx.Exec(`
			INSERT INTO org_pbc_notification_cadence_addtl_users (cadence_id, user_id)
			SELECT cs.id, $1
			FROM org_pbc_notification_cadence_settings AS cs
			WHERE cs.org_id = $2
		`, uid, setting.OrgId)

		if err != nil {
			return err
		}
	}

	return nil
}

func GetPbcNotificationRecord(orgId int32) (map[core.PbcNotificationRecordKey]bool, error) {
	ret := map[core.PbcNotificationRecordKey]bool{}

	data := make([]core.PbcNotificationRecordKey, 0)
	err := dbConn.Select(&data, `
		SELECT cadence_id, org_id, request_id
		FROM org_pbc_notification_record
		WHERE org_id = $1
	`, orgId)

	if err != nil {
		return nil, err
	}

	for _, d := range data {
		ret[d] = true
	}

	return ret, nil
}

func MarkPbcNotificationRecordWithTx(tx *sqlx.Tx, rec core.PbcNotificationRecordKey) error {
	_, err := tx.NamedExec(`
		INSERT INTO org_pbc_notification_record (cadence_id, org_id, request_id, notif_time)
		VALUES (:cadence_id, :org_id, :request_id, NOW())
	`, rec)
	return err
}
