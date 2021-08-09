package v1

import (
	"cloud-batch/configs"
	"cloud-batch/internal/pkg/app"
	"cloud-batch/internal/pkg/e"
	"github.com/gin-gonic/gin"
)

// GetVersion
// @Summary Get Version
// @Tags Version
// @Produce  json
// @Success 200 {object} app.ResponseString
// @Failure 500 {object} app.ResponseString
// @Router /version [get]
func GetVersion(c *gin.Context) {
	appG := app.Gin{C: c}

	appG.ResponseI18nMsg(e.StatusOK, nil, nil, nil, gin.H{
		"version": configs.Server.Version,
	})
}
