package main

import (
	"github.com/gorilla/mux"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"gitlab.com/b3h47pte/audit-stuff/mail_api"
	"gitlab.com/b3h47pte/audit-stuff/okta_api"
	"gitlab.com/b3h47pte/audit-stuff/render"
	"gitlab.com/b3h47pte/audit-stuff/rest"
	"gitlab.com/b3h47pte/audit-stuff/vault_api"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
	"gitlab.com/b3h47pte/audit-stuff/websocket"
	"net/http"
	"os"
	"time"
)

func main() {
	core.Init()
	database.Init()
	render.RegisterTemplates()
	webcore.InitializeWebcore()
	mail.InitializeMailAPI(core.EnvConfig.Mail.Provider, core.EnvConfig.Mail.Key)
	okta.InitializeOktaAPI(okta.OktaConfig{
		ApiKey:    core.EnvConfig.Okta.ApiKey,
		ApiDomain: core.EnvConfig.Okta.BaseUrl,
	})
	vault.Initialize(vault.VaultConfig{
		Url:   core.EnvConfig.Vault.Url,
		Token: core.EnvConfig.Vault.Token,
	})

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
	pageRouter := dynamicRouter.Methods("GET").Subrouter()
	pageRouter.Use(webcore.ObtainUserSessionInContextMiddleware)
	pageRouter.Use(webcore.GrantCSRFMiddleware)
	pageRouter.HandleFunc(core.GetStartedUrl, render.RenderGettingStartedPage).Name(string(webcore.GettingStartedRouteName))
	pageRouter.HandleFunc(core.ContactUsUrl, render.RenderContactUsPage).Name(string(webcore.ContactUsRouteName))
	pageRouter.HandleFunc(core.HomePageUrl, render.RenderHomePage).Name(string(webcore.LandingPageRouteName))
	pageRouter.HandleFunc(core.LoginUrl, render.RenderLoginPage).Name(string(webcore.LoginRouteName))
	pageRouter.HandleFunc(core.RegisterUrl, render.RenderRegisterPage).Name(string(webcore.RegisterRouteName))
	pageRouter.HandleFunc(core.LearnMoreUrl, render.RenderLearnMorePage).Name(string(webcore.LearnMoreRouteName))
	createDashboardSubrouter(pageRouter)

	rest.RegisterPaths(dynamicRouter)
	websocket.RegisterPaths(dynamicRouter)

	// Should be last?
	webcore.RegisterRouter(r)

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
