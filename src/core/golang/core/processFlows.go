package core

import (
	"time"
)

type ProcessFlow struct {
	Id              int64         `db:"id"`
	Name            string        `db:"name"`
	Org             *Organization `db:"org" json:"-"`
	Description     string        `db:"description"`
	CreationTime    time.Time     `db:"created_time"`
	LastUpdatedTime time.Time     `db:"last_updated_time"`
}

type NodeRiskRelationship struct {
	NodeId int64 `db:"node_id"`
	RiskId int64 `db:"risk_id"`
}

type NodeControlRelationship struct {
	NodeId    int64 `db:"node_id"`
	ControlId int64 `db:"control_id"`
}

type RiskControlRelationship struct {
	RiskId    int64 `db:"risk_id"`
	ControlId int64 `db:"control_id"`
}

type ProcessFlowGraph struct {
	Nodes       []*ProcessFlowNode
	Edges       []*ProcessFlowEdge
	Risks       []*Risk
	Controls    []*Control
	NodeRisk    []*NodeRiskRelationship
	NodeControl []*NodeControlRelationship
	RiskControl []*RiskControlRelationship
}
