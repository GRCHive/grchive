package core

import "time"

type Notification struct {
	Id                 int64     `db:"id"`
	OrgId              int32     `db:"org_id"`
	Time               time.Time `db:"time"`
	SubjectType        string    `db:"subject_type"`
	SubjectId          int64     `db:"subject_id"`
	Verb               string    `db:"verb"`
	ObjectType         string    `db:"object_type"`
	ObjectId           int64     `db:"object_id"`
	IndirectObjectType string    `db:"indirect_object_type"`
	IndirectObjectId   int64     `db:"indirect_object_id"`
}
