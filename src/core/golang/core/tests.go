package core

type CodeRunTestSummary struct {
	SuccessfulTests int32
	TotalTests      int32
}

type TrackedSource struct {
	Id     int64  `db:"id"`
	RunId  int64  `db:"run_id"`
	OrgId  int32  `db:"org_id"`
	DataId int64  `db:"data_id"`
	Src    string `db:"src"`
}

type TrackedData struct {
	Id       int64  `db:"id"`
	SourceId int64  `db:"source_id"`
	Data     string `db:"data"`
}

type TrackedTest struct {
	Id      int64     `db:"id"`
	DataAId NullInt64 `db:"data_a_id"`
	DataBId NullInt64 `db:"data_b_id"`
	Ok      bool      `db:"ok"`
	Action  string    `db:"action"`
	Field   string    `db:"field"`
}
