package webcore

import (
	"github.com/gorilla/mux"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"net/http"
)

func LoggedRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		core.Info(
			"START: Remote: ", r.RemoteAddr,
			" URL: ", r.URL,
			" Method: ", r.Method)
		next.ServeHTTP(w, r)
	})
}

func ObtainUserSessionInContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, newR, err := FindValidUserSession(w, r)
		if err != nil && session == nil {
			core.Info("Error in finding valid user session: " + err.Error())
			next.ServeHTTP(w, r)
			return
		}

		data, err := ExtractParsedDataFromSession(session)
		if err != nil {
			core.Info("Error in parsing user session: " + err.Error())
			next.ServeHTTP(w, r)
			return
		}

		context := AddSessionParsedDataToContext(data, newR.Context())
		newR = newR.WithContext(context)
		next.ServeHTTP(w, newR)
	})
}

func CreateAuthenticatedRequestMiddleware(failure http.HandlerFunc) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := FindSessionInContext(r.Context())
			if err != nil {
				core.Info("No user session: " + err.Error())
				failure.ServeHTTP(w, r)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}

func ObtainOrganizationInfoFromRequestInContextMiddleware(next http.Handler) http.Handler {
	// If we can't find the organization we should direct to the dashboard home page.
	// Note that this runs under the assumption that we won't ever have the case where
	// the dashboard home page directs to an invalid org...
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		org, err := GetOrganizationFromRequestUrl(r)
		if err != nil {
			core.Info("Bad organization: " + err.Error())
			http.Redirect(w, r, MustGetRouteUrl(DashboardHomeRouteName), http.StatusTemporaryRedirect)
			return
		}

		ctx := AddOrganizationInfoToContext(org, r.Context())
		newR := r.WithContext(ctx)
		next.ServeHTTP(w, newR)
	})
}

func CreateVerifyUserHasAccessToOrganizationMiddleware(failure http.HandlerFunc) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			organization, err := FindOrganizationInContext(r.Context())
			if err != nil {
				core.Info("No current organization: " + err.Error())
				failure.ServeHTTP(w, r)
				return
			}

			currentData, err := FindSessionParsedDataInContext(r.Context())
			if err != nil {
				core.Info("No current user: " + err.Error())
				failure.ServeHTTP(w, r)
				return
			}

			if currentData.Org.OktaGroupId != organization.OktaGroupId {
				core.Info("Unauthenticated user: " + currentData.CurrentUser.Email)
				failure.ServeHTTP(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func CreateVerifyUserHasAccessToUserMiddleware(failure http.HandlerFunc) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check that the user stored in the session matches the email specified
			// in the URL. We will probably need coarser access permissions later so that
			// you can view the profile of a co-worker (for example).
			email, err := GetUserEmailFromRequestUrl(r)
			if err != nil {
				core.Info("No user email: " + err.Error())
				failure.ServeHTTP(w, r)
				return
			}

			currentData, err := FindSessionParsedDataInContext(r.Context())
			if err != nil {
				core.Info("No current user: " + err.Error())
				failure.ServeHTTP(w, r)
				return
			}

			if currentData.CurrentUser.Email != email {
				core.Info("Unauthenticated user: " + currentData.CurrentUser.Email)
				failure.ServeHTTP(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
