package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	r = r.Delims("[[", "]]")

	// Static assets that can eventually be served by Nginx.
	r.Static("/static/corejsui", "src/core/jsui/dist")

	// Dynamic(?) content that needs to be served by Go.
	r.LoadHTMLGlob("src/webserver/templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", loadGlobalProps())
	})

	// TODO: Configurable port?
	r.Run(":8080")
}
