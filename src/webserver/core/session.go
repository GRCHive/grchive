package core

import (
	"github.com/gorilla/Sessions"
)

var SessionStore = sessions.NewCookieStore(LoadEnvConfig().SessionKeys...)
