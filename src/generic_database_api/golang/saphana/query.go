package saphana_api

import (
	"context"
	"gitlab.com/grchive/grchive/db_api/utility"
	"time"
)

func (hdb *SapHanaDriver) RunQuery(query string) (*utility.SqlQueryResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	rows, err := hdb.connection.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return utility.CreateSqlQueryResultFromRows(rows)
}
