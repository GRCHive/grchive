package database

import (
	"errors"
	"fmt"
	"gitlab.com/grchive/grchive/core"
)

func GetResourceName(typ string, id int64) (string, error) {
	rows, err := dbConn.Queryx(fmt.Sprintf(`
		SELECT *
		FROM %s
		WHERE id = $1	
	`, typ), id)

	if err != nil {
		return "", err
	}

	defer rows.Close()
	rows.Next()

	data := map[string]interface{}{}
	err = rows.MapScan(data)
	if err != nil {
		return "", err
	}

	switch typ {
	case core.ResourceIdUser:
		return fmt.Sprintf("%s %s (%s)",
			data["first_name"].(string),
			data["last_name"].(string),
			data["email"].(string),
		), nil
	case core.ResourceIdControl:
		return fmt.Sprintf("%s (%s)",
			data["name"].(string),
			data["identifier"].(string),
		), nil
	case core.ResourceIdDocRequest:
		return data["name"].(string), nil
	case core.ResourceIdSqlQueryRequest:
		return data["name"].(string), nil
	case core.ResourceIdDatabase:
		return data["name"].(string), nil
	case core.ResourceIdClientData:
		return data["name"].(string), nil
	case core.ResourceIdGenericRequests:
		return data["name"].(string), nil
	}

	return "", errors.New("Unsupported resource type: " + typ)
}
