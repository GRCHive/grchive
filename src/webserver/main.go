package main

import (
	"github.com/gorilla/mux"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"gitlab.com/b3h47pte/audit-stuff/rest"
	"net/http"
	"os"
	"time"
)

func loggedRequestHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		core.Info(
			"Remote: ", r.RemoteAddr,
			" URL: ", r.URL,
			" Method: ", r.Method)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	database.Init()
	registerTemplates()

	r := mux.NewRouter()

	// Static assets that can eventually be served by Nginx.
	_, err := os.Stat("src/core/jsui/dist-smap")
	if os.IsNotExist(err) {
		r.PathPrefix("/static/corejsui/").Handler(
			http.StripPrefix(
				"/static/corejsui/",
				http.FileServer(http.Dir("src/core/jsui/dist-nosmap"))))
	} else {
		r.PathPrefix("/static/corejsui/").Handler(
			http.StripPrefix(
				"/static/corejsui/",
				http.FileServer(http.Dir("src/core/jsui/dist-smap"))))
	}
	r.PathPrefix("/static/assets/").Handler(
		http.StripPrefix(
			"/static/assets/",
			http.FileServer(http.Dir("src/core/jsui/assets"))))

	// Dynamic(?) content that needs to be served by Go.
	r.HandleFunc(core.GetStartedUrl, renderGettingStartedPage).Methods("GET")
	r.HandleFunc(core.ContactUsUrl, renderContactUsPage).Methods("GET")
	r.HandleFunc(core.HomePageUrl, renderHomePage).Methods("GET")
	r.HandleFunc(core.LoginUrl, renderLoginPage).Methods("GET")
	r.HandleFunc(core.LearnMoreUrl, renderLearnMorePage).Methods("GET")
	rest.RegisterPaths(r)

	//// TODO: Configurable port?
	srv := &http.Server{
		Handler:      loggedRequestHandler(r),
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	if err = srv.ListenAndServe(); err != nil {
		core.Error(err.Error())
	}
}
