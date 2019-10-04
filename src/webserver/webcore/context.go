package webcore

import (
	"context"
	"errors"
	"gitlab.com/b3h47pte/audit-stuff/core"
)

type ContextKey string

const (
	// Session (cookie)
	UserSessionContextKey           ContextKey = "SESSION"
	UserSessionParsedDataContextKey            = "PARSEDDATA"
	// From Request
	OrganizationContextKey = "ORGANIZATION"
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

func AddSessionParsedDataToContext(data *core.UserSessionParsedData, ctx context.Context) context.Context {
	return context.WithValue(ctx, UserSessionParsedDataContextKey, data)
}

func FindSessionParsedDataInContext(ctx context.Context) (*core.UserSessionParsedData, error) {
	data, ok := ctx.Value(UserSessionParsedDataContextKey).(*core.UserSessionParsedData)
	if !ok || data == nil {
		return nil, errors.New("Failed to find session parsed data in context.")
	}
	return data, nil
}

func AddOrganizationInfoToContext(org *core.Organization, ctx context.Context) context.Context {
	return context.WithValue(ctx, OrganizationContextKey, org)
}

func FindOrganizationInContext(ctx context.Context) (*core.Organization, error) {
	org, ok := ctx.Value(OrganizationContextKey).(*core.Organization)
	if !ok || org == nil {
		return nil, errors.New("Failed to find organization in context.")
	}
	return org, nil
}
