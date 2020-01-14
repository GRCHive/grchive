package core

import (
	"github.com/google/uuid"
)

type UuidGen interface {
	GenStr() string
}

type GoogleUuidGen struct{}

func (g GoogleUuidGen) GenStr() string {
	return uuid.New().String()
}

var DefaultUuidGen = GoogleUuidGen{}
