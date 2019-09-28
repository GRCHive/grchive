package rest

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"net/http"
	"strings"
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
		c.JSON(http.StatusBadRequest, struct{}{})
		return
	}

	data.Name = strings.TrimSpace(data.Name)
	data.Name = strings.TrimSpace(data.Name)

	if data.Name == "" || data.Email == "" {
		core.Warning("Empty name or email.")
		c.JSON(http.StatusBadRequest, struct{}{})
		return
	}

	// Save name and email to the database.
	isDuplicate, err := database.AddNewGettingStartedInterest(data.Name, data.Email)

	// If the error is related to having a duplicate then we should let the user know.
	// Otherwise, our service probably failed somewhere which hopefully got logged.
	if err != nil {
		if isDuplicate {
			core.Warning("Detected duplicate entry.")
			c.JSON(http.StatusBadRequest, struct {
				IsDuplicate bool
			}{
				IsDuplicate: true,
			})
		} else {
			core.Warning("Failed to add getting started interest.")
			c.JSON(http.StatusInternalServerError, struct{}{})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
