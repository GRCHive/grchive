package core_test

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/grchive/grchive/core"
	"testing"
)

func TestParseEmail(t *testing.T) {
	for _, test := range []struct {
		email  string
		user   string
		domain string
	}{
		{"bob@hello.com", "bob", "hello.com"},
		{"@hello.com", "", "hello.com"},
		{"bob@", "bob", ""},
		{"@", "", ""},
		{"cheese@test@.com", "", ""},
	} {
		parsedUser, parsedDomain := core.ParseEmailAddress(test.email)
		assert.Equal(t, test.user, parsedUser)
		assert.Equal(t, test.domain, parsedDomain)
	}

}
