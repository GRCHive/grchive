package core

import "strings"

// Takes in an input address in the form of USER@DOMAIN
// and returns (USER, DOMAIN).
func ParseEmailAddress(email string) (string, string) {
	split := strings.Split(email, "@")
	if len(split) != 2 {
		return "", ""
	}
	return split[0], split[1]
}
