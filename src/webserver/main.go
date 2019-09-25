package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	// TODO: Configurable port?
	r.Run(":8080")
}
