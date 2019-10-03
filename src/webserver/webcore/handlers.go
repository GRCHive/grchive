package webcore

import (
	"github.com/gorilla/mux"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"net/http"
)

func LoggedRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		core.Info(
			"Remote: ", r.RemoteAddr,
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
		newR = newR.Clone(context)
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
