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
	case core.ResourceUser:
		return fmt.Sprintf("%s %s (%s)",
			data["first_name"].(string),
			data["last_name"].(string),
			data["email"].(string),
		), nil
	case core.ResourceControl:
		return fmt.Sprintf("%s (%s)",
			data["name"].(string),
			data["identifier"].(string),
		), nil
	case core.ResourceDocRequest:
		return data["name"].(string), nil
	case core.ResourceSqlQueryRequest:
		return data["name"].(string), nil
	}

	return "", errors.New("Unsupported resource type: " + typ)
}
