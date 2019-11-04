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

func CreateObtainOrganizationInfoFromUserInContextMiddleware(failure http.HandlerFunc) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			currentData, err := FindSessionParsedDataInContext(r.Context())
			if err != nil {
				core.Info("No current user: " + err.Error())
				failure.ServeHTTP(w, r)
				return
			}

			ctx := AddOrganizationInfoToContext(currentData.Org, r.Context())
			newR := r.WithContext(ctx)
			next.ServeHTTP(w, newR)
		})
	}
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

func GrantAPIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userParsedData, err := FindSessionParsedDataInContext(r.Context())
		if err != nil {
			core.Info("Failed to add find user for API key: " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = RefreshGrantAPIKey(userParsedData.CurrentUser.Id, w, r)
		if err != nil {
			core.Info("Failed to grant API key: " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GrantCSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := AddCSRFTokenToRequest(w, r); err != nil {
			core.Info("Failed to add CSRF: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func CreateVerifyCSRFMiddleware(failure http.HandlerFunc) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := GetCSRFToken(r)
			if err != nil {
				core.Info("Failed to get CSRF token: " + err.Error())
				failure.ServeHTTP(w, r)
				return
			}

			ok, err := VerifyCSRFToken(token, r)
			if !ok || err != nil {
				core.Info("Failed to verify CSRF token: " + core.ErrorString(err))
				failure.ServeHTTP(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Redirect non-GET 301 and 302 to 308 and 307 respectively. The 307 and 308
// status codes will pass the original body and method to the redirected path
// which is desirable behavior for us whereas 301 and 302 will generally always
// redirect to a GET.
// See:
// 	301: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/301
//  302: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/302
//  307: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/307
//  308: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/308
func HTTPRedirectStatusCodes(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			wr := RedirectResponseWriter{w}
			next.ServeHTTP(wr, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
