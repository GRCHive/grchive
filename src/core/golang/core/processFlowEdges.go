package core

type ProcessFlowEdge struct {
	Id         int64 `db:"id"`
	InputIoId  int64 `db:"input_id"`
	OutputIoId int64 `db:"output_id"`
}
