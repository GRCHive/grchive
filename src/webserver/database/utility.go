package database

import (
	"github.com/lib/pq"
)

// Checks whether the error indicates a duplicate entry on INSERT.
func IsDuplicateDBEntry(err error) bool {
	if err == nil {
		return false
	}

	switch err.(type) {
	case *pq.Error:
		return err.(*pq.Error).Code == "23505"
	}
	return false
}
