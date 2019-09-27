package rest

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/b3h47pte/audit-stuff/core"
)

func RegisterPaths(r *gin.Engine) {
	r.POST(core.CreateGetStartedUrl(), postGettingStartedInterest)
}
