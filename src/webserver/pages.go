package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func renderGettingStartedPage(c *gin.Context) {
	c.HTML(http.StatusOK, "gettingStarted.tmpl", loadGlobalProps())
}

func renderContactUsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "contactUs.tmpl", loadGlobalProps())
}

func renderHomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", loadGlobalProps())
}

func renderLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", loadGlobalProps())
}
