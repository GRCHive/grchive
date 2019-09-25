package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func main() {
	r := gin.Default()

	dir, _ := os.Getwd()
	fmt.Println(dir)

	r.LoadHTMLGlob("src/webserver/templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})

	// TODO: Configurable port?
	r.Run(":8080")
}
