package core

import "time"

type EventVerb string

const (
	VerbGettingStarted EventVerb = "requested to get started"
	VerbAssign                   = "assigned"
	VerbComplete                 = "completed"
	VerbReopen                   = "reopened"
	VerbComment                  = "commented on"
	VerbApprove                  = "approved"
	VerbReject                   = "rejected"
)

type Event struct {
	Subject        User
	Verb           EventVerb
	Object         interface{}
	IndirectObject interface{}
	Timestamp      time.Time
}
