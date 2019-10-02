package webcore

import (
	"context"
	"errors"
	"gitlab.com/b3h47pte/audit-stuff/core"
)

type ContextKey string

const (
	UserSessionContextKey ContextKey = "CONTEXTKEY"
)

func AddSessionToContext(session *core.UserSession, ctx context.Context) context.Context {
	return context.WithValue(ctx, UserSessionContextKey, session)
}

func FindSessionInContext(ctx context.Context) (*core.UserSession, error) {
	session, ok := ctx.Value(UserSessionContextKey).(*core.UserSession)
	if !ok || session == nil {
		return nil, errors.New("Failed to find session in context.")
	}
	return session, nil
}
