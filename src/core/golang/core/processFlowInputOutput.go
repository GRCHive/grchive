package core

type ProcessFlowIOType struct {
	Id   int64  `db:"id"`
	Name string `db:"name"`
}

type ProcessFlowInputOutput struct {
	Id           int64  `db:"id"`
	Name         string `db:"name"`
	ParentNodeId int64  `db:"parent_node_id"`
	TypeId       int32  `db:"io_type_id"`
}
