package core

type ProcessFlowNodeType struct {
	Id          int32  `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

type ProcessFlowNode struct {
	Id            int64  `db:"id"`
	Name          string `db:"name"`
	ProcessFlowId int64  `db:"process_flow_id"`
	Description   string `db:"description"`
	NodeTypeId    int32  `db:"node_type"`
	Inputs        []ProcessFlowInputOutput
	Outputs       []ProcessFlowInputOutput
	RiskIds       []int64
}
