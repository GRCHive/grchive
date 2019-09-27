package main

import (
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
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
	r.GET(createGetStartedUrl(), renderGettingStartedPage)
	r.GET(createContactUsUrl(), renderContactUsPage)
	r.GET(createHomePageUrl(), renderHomePage)
	r.GET(createLoginUrl(), renderLoginPage)
	r.GET(createLearnMoreUrl(), renderLearnMorePage)

	// TODO: Configurable port?
	r.Run(":8080")
}
