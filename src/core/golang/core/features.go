package core

type FeatureId int64

const (
	AutomationFeature FeatureId = 1
)

type Feature struct {
	Id   FeatureId `db:"id"`
	Name string    `db:"name"`
}
