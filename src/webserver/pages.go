package main

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"net/http"
)

func renderGettingStartedPage(c *gin.Context) {
	c.HTML(http.StatusOK, "gettingStarted.tmpl", core.LoadGlobalProps())
}

func renderContactUsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "contactUs.tmpl", core.LoadGlobalProps())
}

func renderHomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", core.LoadGlobalProps())
}

func renderLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", core.LoadGlobalProps())
}

func renderLearnMorePage(c *gin.Context) {
	c.HTML(http.StatusOK, "learnMore.tmpl", core.LoadGlobalProps())
}
