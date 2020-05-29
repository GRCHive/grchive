package webcore

import (
	"context"
	"errors"
	"gitlab.com/grchive/grchive/core"
)

type ContextKey string

const (
	// Session (cookie)
	UserSessionContextKey           ContextKey = "SESSION"
	UserSessionParsedDataContextKey            = "PARSEDDATA"
	// From Request
	OrganizationContextKey                = "ORGANIZATION"
	UserContextKey                        = "USER"
	ApiKeyContextKey                      = "APIKEY"
	GenericRequestContextKey              = "GENREQ"
	RoleContextKey                        = "ROLE"
	ShellScriptContextKey                 = "SHELLSCRIPT"
	ShellScriptVersionContextKey          = "SHELLSCRIPTVERSION"
	ServerContextKey                      = "SERVER"
	ServerConnectionSshPasswordContextKey = "SERVERCONNECTIONSSHPASSWORD"
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

func AddUserToContext(user *core.User, ctx context.Context) context.Context {
	return context.WithValue(ctx, UserContextKey, user)
}

func FindUserInContext(ctx context.Context) (*core.User, error) {
	user, ok := ctx.Value(UserContextKey).(*core.User)
	if !ok || user == nil {
		return nil, errors.New("Failed to find user in context.")
	}
	return user, nil
}

func AddApiKeyToContext(key *core.ApiKey, ctx context.Context) context.Context {
	return context.WithValue(ctx, ApiKeyContextKey, key)
}

func FindApiKeyInContext(ctx context.Context) (*core.ApiKey, error) {
	key, ok := ctx.Value(ApiKeyContextKey).(*core.ApiKey)
	if !ok || key == nil {
		return nil, errors.New("Failed to find api keyin context.")
	}
	return key, nil
}

func AddRoleToContext(role *core.Role, ctx context.Context) context.Context {
	return context.WithValue(ctx, RoleContextKey, role)
}

func FindRoleInContext(ctx context.Context) (*core.Role, error) {
	role, ok := ctx.Value(RoleContextKey).(*core.Role)
	if !ok || role == nil {
		return nil, errors.New("Failed to find role in context.")
	}
	return role, nil
}

func AddShellScriptToContext(script *core.ShellScript, ctx context.Context) context.Context {
	return context.WithValue(ctx, ShellScriptContextKey, script)
}

func FindShellScriptInContext(ctx context.Context) (*core.ShellScript, error) {
	script, ok := ctx.Value(ShellScriptContextKey).(*core.ShellScript)
	if !ok || script == nil {
		return nil, errors.New("Failed to find shell script in context.")
	}
	return script, nil
}

func FindShellScriptVersionInContext(ctx context.Context) (*core.ShellScriptVersion, error) {
	version, ok := ctx.Value(ShellScriptVersionContextKey).(*core.ShellScriptVersion)
	if !ok || version == nil {
		return nil, errors.New("Failed to find shell script version in context.")
	}
	return version, nil
}

func FindServerInContext(ctx context.Context) (*core.Server, error) {
	server, ok := ctx.Value(ServerContextKey).(*core.Server)
	if !ok || server == nil {
		return nil, errors.New("Failed to find server in context.")
	}
	return server, nil
}

func FindServerSSHPasswordConnectionInContext(ctx context.Context) (*core.ServerSSHPasswordConnection, error) {
	resource, ok := ctx.Value(ServerConnectionSshPasswordContextKey).(*core.ServerSSHPasswordConnection)
	if !ok || resource == nil {
		return nil, errors.New("Failed to find ServerSSHPasswordConnection in context.")
	}
	return resource, nil
}
