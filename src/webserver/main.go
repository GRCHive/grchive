package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("src/webserver/templates/*")
	r.Static("/static/corejsui", "src/core/jsui/dist")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})

	// TODO: Configurable port?
	r.Run(":8080")
}
