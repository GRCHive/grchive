package database

import (
	"gitlab.com/grchive/grchive/core"
)

func GetSupportedCodeParameterTypes() ([]*core.SupportedCodeParameterType, error) {
	typs := make([]*core.SupportedCodeParameterType, 0)
	err := dbConn.Select(&typs, `
		SELECT *
		FROM supported_code_parameter_types
		ORDER BY name ASC
	`)
	return typs, err
}

func GetSingleSupportedCodeParameterTypeFromId(id int32) (*core.SupportedCodeParameterType, error) {
	typ := core.SupportedCodeParameterType{}
	err := dbConn.Get(&typ, `
		SELECT *
		FROM supported_code_parameter_types
		WHERE id = $1
	`, id)
	return &typ, err
}
