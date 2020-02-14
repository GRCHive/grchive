package utility

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type SqlQueryResult struct {
	Columns []string
	CsvText string
}

func CreateSqlQueryResultFromRows(rows *sqlx.Rows) (*SqlQueryResult, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	writer := csv.NewWriter(buf)
	record := make([]string, len(columns))

	dest := map[string]interface{}{}
	for rows.Next() {
		err = rows.MapScan(dest)
		if err != nil {
			return nil, err
		}

		for i, col := range columns {
			record[i] = fmt.Sprintf("%v", dest[col])
		}

		err = writer.Write(record)
		if err != nil {
			return nil, err
		}
	}

	writer.Flush()

	return &SqlQueryResult{
		Columns: columns,
		CsvText: buf.String(),
	}, nil
}
