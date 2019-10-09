package core

type ProcessFlowNodeType struct {
	Id          uint32 `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}
