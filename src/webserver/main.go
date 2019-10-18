package main

import (
	"github.com/gorilla/mux"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"gitlab.com/b3h47pte/audit-stuff/render"
	"gitlab.com/b3h47pte/audit-stuff/rest"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
	"gitlab.com/b3h47pte/audit-stuff/websocket"
	"net/http"
	"os"
	"time"
)

func main() {
	database.Init()
	render.RegisterTemplates()
	webcore.InitializeSessions()

	r := mux.NewRouter().StrictSlash(true)
	r.Use(webcore.HTTPRedirectStatusCodes)

	staticRouter := r.PathPrefix("/static").Subrouter()

	// Static assets that can eventually be served by Nginx.
	_, err := os.Stat("src/core/jsui/dist-smap")
	if os.IsNotExist(err) {
		staticRouter.PathPrefix("/corejsui/").Handler(
			http.StripPrefix(
				"/static/corejsui/",
				http.FileServer(http.Dir("src/core/jsui/dist-nosmap"))))
	} else {
		staticRouter.PathPrefix("/corejsui/").Handler(
			http.StripPrefix(
				"/static/corejsui/",
				http.FileServer(http.Dir("src/core/jsui/dist-smap"))))
	}
	staticRouter.PathPrefix("/assets/").Handler(
		http.StripPrefix(
			"/static/assets/",
			http.FileServer(http.Dir("src/core/jsui/assets"))))

	dynamicRouter := r.PathPrefix("/").Subrouter()

	// Dynamic(?) content that needs to be served by Go.
	dynamicRouter.Use(webcore.ObtainUserSessionInContextMiddleware)
	dynamicRouter.HandleFunc(core.GetStartedUrl, render.RenderGettingStartedPage).Methods("GET").Name(string(webcore.GettingStartedRouteName))
	dynamicRouter.HandleFunc(core.ContactUsUrl, render.RenderContactUsPage).Methods("GET").Name(string(webcore.ContactUsRouteName))
	dynamicRouter.HandleFunc(core.HomePageUrl, render.RenderHomePage).Methods("GET").Name(string(webcore.LandingPageRouteName))
	dynamicRouter.HandleFunc(core.LoginUrl, render.RenderLoginPage).Methods("GET").Name(string(webcore.LoginRouteName))
	dynamicRouter.HandleFunc(core.LearnMoreUrl, render.RenderLearnMorePage).Methods("GET").Name(string(webcore.LearnMoreRouteName))
	rest.RegisterPaths(dynamicRouter)
	websocket.RegisterPaths(dynamicRouter)
	createDashboardSubrouter(dynamicRouter)
	webcore.RegisterRouter(dynamicRouter)

	// Custom 404
	r.NotFoundHandler = http.HandlerFunc(render.Render404)

	// TODO: Configurable port?
	srv := &http.Server{
		Handler:      webcore.LoggedRequestMiddleware(r),
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err = srv.ListenAndServe(); err != http.ErrServerClosed {
		core.Error(err.Error())
	}
}
