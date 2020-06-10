package core

import (
	"encoding/json"
	"time"
)

type SapErpRfc struct {
	Id            int64  `db:"id"`
	IntegrationId int64  `db:"integration_id"`
	Function      string `db:"function_name"`
}

type SapErpRfcVersion struct {
	Id           int64            `db:"id"`
	RfcId        int64            `db:"rfc_id"`
	CreatedTime  time.Time        `db:"created_time"`
	FinishedTime NullTime         `db:"finished_time"`
	RawData      NullString       `db:"data" json:"-"`
	Data         *json.RawMessage `db:"-"`
	Success      bool             `db:"success"`
	Logs         NullString       `db:"-"`
}
