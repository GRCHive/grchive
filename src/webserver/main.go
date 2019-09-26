package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
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

	// Dynamic(?) content that needs to be served by Go.
	r.LoadHTMLGlob("src/webserver/templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", loadGlobalProps())
	})

	// TODO: Configurable port?
	r.Run(":8080")
}
