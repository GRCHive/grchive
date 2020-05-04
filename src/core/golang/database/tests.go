package database

import (
	"gitlab.com/grchive/grchive/core"
)

func GetCodeRunTestSummary(runId int64, orgId int32, role *core.Role) (*core.CodeRunTestSummary, error) {
	rows, err := dbConn.Queryx(`
		WITH tests AS (
			SELECT t.*
			FROM test_tests AS t
			LEFT JOIN test_data AS tda
				ON tda.id = t.data_a_id
			LEFT JOIN test_sources AS tsa
				ON tsa.id = tda.source_id
			LEFT JOIN test_data AS tdb
				ON tdb.id = t.data_b_id
			LEFT JOIN test_sources AS tsb
				ON tsb.id = tdb.source_id
			WHERE 
				tsa.run_id = $1 AND
				tsa.org_id = $2 AND
				tsb.run_id = $1 AND
				tsb.org_id = $2
		)
		SELECT 
			(
				SELECT COUNT(*)
				FROM tests
				WHERE ok = 'true'
			),
			COUNT(tests.*)
		FROM tests
	`, runId, orgId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return &core.CodeRunTestSummary{}, nil
	}

	summary := core.CodeRunTestSummary{}
	err = rows.Scan(&summary.SuccessfulTests, &summary.TotalTests)
	if err != nil {
		return nil, err
	}

	return &summary, nil
}
