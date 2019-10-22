package core

type Risk struct {
	Id              int64  `db:"id"`
	Name            string `db:"name"`
	Description     string `db:"description"`
	RelevantNodeIds []int64
}
