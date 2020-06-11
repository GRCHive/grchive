package utility

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type SqlQueryResult struct {
	Columns []string
	CsvText string
}

func CreateSqlQueryResultFromRows(rows *sql.Rows) (*SqlQueryResult, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	writer := csv.NewWriter(buf)
	record := make([]string, len(columns))

	dest := map[string]interface{}{}
	for rows.Next() {
		err = sqlx.MapScan(rows, dest)
		if err != nil {
			return nil, err
		}

		for i, col := range columns {
			val := dest[col]
			switch val.(type) {
			case []uint8:
				record[i] = string(val.([]uint8))
			default:
				record[i] = fmt.Sprintf("%v", dest[col])
			}
		}

		// record represents a single row so we need to write
		// for each row.
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
