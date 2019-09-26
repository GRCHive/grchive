package main

import (
	"github.com/gin-gonic/gin"
)

func loadGlobalProps() gin.H {
	return gin.H{
		"companyName": "Audit Stuff",
	}
}
