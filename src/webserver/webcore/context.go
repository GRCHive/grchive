package webcore

import (
	"context"
	"gitlab.com/b3h47pte/audit-stuff/core"
)

type ContextKey string

const (
	UserSessionContextKey ContextKey = "CONTEXTKEY"
)

func AddSessionToContext(session *core.UserSession, ctx context.Context) context.Context {
	return context.WithValue(ctx, UserSessionContextKey, session)
}
