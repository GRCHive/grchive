package main

import (
	"archive/zip"
	"github.com/gorilla/mux"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/gcloud_api"
	"gitlab.com/grchive/grchive/mail_api"
	"gitlab.com/grchive/grchive/okta_api"
	"gitlab.com/grchive/grchive/render"
	"gitlab.com/grchive/grchive/rest"
	"gitlab.com/grchive/grchive/vault_api"
	"gitlab.com/grchive/grchive/webcore"
	"gitlab.com/grchive/grchive/websocket"
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/zipfs"
	"net/http"
	"os"
	"time"
)

type ZipFile struct {
	rsc  vfs.ReadSeekCloser
	name string
	zip  vfs.FileSystem
}

func (f ZipFile) Read(p []byte) (n int, err error) {
	return f.rsc.Read(p)
}

func (f ZipFile) Close() error {
	return f.rsc.Close()
}

func (f ZipFile) Seek(offset int64, whence int) (int64, error) {
	return f.rsc.Seek(offset, whence)
}

func (f ZipFile) Readdir(count int) ([]os.FileInfo, error) {
	files, err := f.zip.ReadDir(f.name)
	if err != nil {
		return nil, err
	}
	return files[:count], nil
}

func (f ZipFile) Stat() (os.FileInfo, error) {
	return f.zip.Stat(f.name)
}

type ZipFS struct {
	Prefix string
	Zip    vfs.FileSystem
}

func (z ZipFS) Open(name string) (http.File, error) {
	fullName := z.Prefix + name
	rsc, err := z.Zip.Open(fullName)
	if err != nil {
		core.Info(err.Error())
		return nil, err
	}

	return ZipFile{
		rsc:  rsc,
		name: fullName,
		zip:  z.Zip,
	}, nil
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	core.Debug("Core Init")
	core.Init()

	core.Debug("Database Init")
	database.Init()

	core.Debug("Database Init Listeners")
	database.InitListeners(map[string]database.ListenHandler{
		database.NotifyChannelControlOwner:       onNotifyControlOwnerChange,
		database.NotifyChannelDocRequestAssignee: onNotifyDocRequestAssigneeChange,
		database.NotifyChannelDocRequestStatus:   onNotifyDocRequestStatusChange,
		database.NotifyChannelSqlRequestAssignee: onNotifySqlRequestAssigneeChange,
		database.NotifyChannelSqlRequestStatus:   onNotifySqlRequestApprovalChange,
	})

	core.Debug("Render Init")
	render.RegisterTemplates()

	core.Debug("Webcore Init")
	webcore.InitializeWebcore()

	core.Debug("Mail Init")
	mail.InitializeMailAPI(core.EnvConfig.Mail.Provider, core.EnvConfig.Mail.Key)

	core.Debug("Okta Init")
	okta.InitializeOktaAPI(okta.OktaConfig{
		ApiKey:    core.EnvConfig.Okta.ApiKey,
		ApiDomain: core.EnvConfig.Okta.BaseUrl,
	})

	core.Debug("Vault Init")
	vault.Initialize(vault.VaultConfig{
		Url:      core.EnvConfig.Vault.Url,
		Username: core.EnvConfig.Vault.Username,
		Password: core.EnvConfig.Vault.Password,
	}, core.EnvConfig.Tls.Config())

	core.Debug("GCloud Init")
	gcloud.DefaultGCloudApi.InitFromJson(core.EnvConfig.Gcloud.AuthFilename)

	core.Debug("RabbitMQ Init")
	webcore.DefaultRabbitMQ.Connect(
		*core.EnvConfig.RabbitMQ,
		webcore.QueueConfig{
			NotificationConsume: true,
		},
		core.EnvConfig.Tls)
	defer webcore.DefaultRabbitMQ.Cleanup()

	baser := mux.NewRouter().StrictSlash(true)
	baser.HandleFunc("/healthz", healthCheck)
	baser.Use(webcore.HTTPRedirectStatusCodes)

	r := baser.NewRoute().Subrouter()
	r.Use(webcore.LoggedRequestMiddleware)

	staticRouter := r.PathPrefix("/static").Subrouter()

	// Static assets that can eventually be served by Nginx.
	z, err := zip.OpenReader("src/core/jsui/corejsui.zip")
	if err != nil {
		core.Error("Failed to open corejsui.zip: " + err.Error())
	}
	defer z.Close()

	zfs := zipfs.New(z, "corejsui")
	staticRouter.PathPrefix("/corejsui/").Handler(
		http.StripPrefix(
			"/static/corejsui/",
			http.FileServer(ZipFS{
				Prefix: "/dist",
				Zip:    zfs,
			})))

	staticRouter.PathPrefix("/assets/").Handler(
		http.StripPrefix(
			"/static/assets/",
			http.FileServer(http.Dir("src/core/jsui/assets"))))

	dynamicRouter := r.PathPrefix("/").Subrouter()

	// Dynamic(?) content that needs to be served by Go.
	pageRouter := dynamicRouter.Methods("GET").Subrouter()
	pageRouter.Use(webcore.ObtainUserSessionInContextMiddleware)
	pageRouter.Use(webcore.GrantCSRFMiddleware)
	pageRouter.HandleFunc(core.GetStartedUrl, render.RenderGettingStartedPage)
	pageRouter.HandleFunc(core.ContactUsUrl, render.RenderContactUsPage)
	pageRouter.HandleFunc(core.HomePageUrl, render.RenderHomePage).Name(webcore.LandingPageRouteName)
	pageRouter.HandleFunc(core.LoginUrl, render.RenderLoginPage).Name(webcore.LoginRouteName)
	pageRouter.HandleFunc(core.RegisterUrl, render.RenderRegisterPage).Name(webcore.RegisterRouteName)
	pageRouter.HandleFunc(core.LearnMoreUrl, render.RenderLearnMorePage)
	createDashboardSubrouter(pageRouter)

	rest.RegisterPaths(dynamicRouter)
	websocket.RegisterPaths(dynamicRouter)

	// Should be last?
	webcore.RegisterRouter(r)

	// Custom 404
	r.NotFoundHandler = http.HandlerFunc(render.Render404)

	SetupRabbitMQInterface()

	// TODO: Configurable port?
	srv := &http.Server{
		Handler:      baser,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err = srv.ListenAndServe(); err != http.ErrServerClosed {
		core.Error(err.Error())
	}
}
