package core

import (
	"time"
)

type ProcessFlow struct {
	Id              uint32        `db:"id" json:"-"`
	Name            string        `db:"name"`
	Org             *Organization `db:"org" json:"-"`
	Description     string        `db:"description"`
	CreationTime    time.Time     `db:"created_time"`
	LastUpdatedTime time.Time     `db:"last_updated_time"`
}

type ProcessFlowNode struct {
}

type ProcessFlowEdge struct {
}
