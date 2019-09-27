package rest

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"net/http"
)

type tGettingStartedInterest struct {
	Name  string `form:"name" binding:"required"`
	Email string `form:"email" binding:"required"`
}

func postGettingStartedInterest(c *gin.Context) {
	// Retrieve the client's name and email from the input form.
	var data tGettingStartedInterest
	if err := c.ShouldBind(&data); err != nil {
		core.Warning("Failed to bind data.")
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if data.Name == "" || data.Email == "" {
		core.Warning("Empty name or email.")
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	// Save name and email to the database.

	c.JSON(http.StatusOK, gin.H{})
}
