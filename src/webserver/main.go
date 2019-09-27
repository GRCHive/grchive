package main

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"gitlab.com/b3h47pte/audit-stuff/rest"
	"os"
)

func main() {
	database.Init()

	r := gin.Default()
	r = r.Delims("[[", "]]")

	// Static assets that can eventually be served by Nginx.
	_, err := os.Stat("src/core/jsui/dist-smap")
	if os.IsNotExist(err) {
		r.Static("/static/corejsui", "src/core/jsui/dist-nosmap")
	} else {
		r.Static("/static/corejsui", "src/core/jsui/dist-smap")
	}
	r.Static("/static/assets", "src/core/jsui/assets")

	// Dynamic(?) content that needs to be served by Go.
	r.LoadHTMLGlob("src/webserver/templates/*")
	r.GET(core.CreateGetStartedUrl(), renderGettingStartedPage)
	r.GET(core.CreateContactUsUrl(), renderContactUsPage)
	r.GET(core.CreateHomePageUrl(), renderHomePage)
	r.GET(core.CreateLoginUrl(), renderLoginPage)
	r.GET(core.CreateLearnMoreUrl(), renderLearnMorePage)
	rest.RegisterPaths(r)

	// TODO: Configurable port?
	r.Run(":8080")
}
