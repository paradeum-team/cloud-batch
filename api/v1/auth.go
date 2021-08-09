package v1

import (
	"cloud-batch/internal/models"
	"cloud-batch/internal/pkg/app"
	"cloud-batch/internal/pkg/e"
	"cloud-batch/internal/service"
	"github.com/gin-gonic/gin"
)

// Login
// @Summary post Auth
// @Tags Auth
// @Produce  json
// @Param auth body models.Auth true "auth"
// @Success 200 {object} app.ResponseString
// @Failure 400 {object} app.ResponseString
// @Failure 500 {object} app.ResponseString
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		auth = models.Auth{}
	)

	err := c.BindJSON(&auth)
	if err != nil {
		code := app.VerdictJsonErr(err)
		appG.ResponseI18nMsgSimple(code, err, nil)
		return
	}

	token, code, err := service.Login(auth)
	if err != nil {
		appG.ResponseI18nMsg(code, err, nil, nil, nil)
		return
	}

	appG.ResponseI18nMsg(e.StatusOK, err, nil, nil, gin.H{
		"token": token,
	})
}

// UpdateAuth
// @Summary Update password
// @Tags Auth
// @Produce  json
// @Security ApiKeyAuth
// @Param updateAuth body models.UpdateAuth true "updateAuth"
// @Success 200 {object} app.ResponseString
// @Failure 400,401 {object} app.ResponseString
// @Failure 500 {object} app.ResponseString
// @Router /auth/passwd [put]
func UpdateAuth(c *gin.Context) {

	appG := app.Gin{C: c}
	updateAuth := models.UpdateAuth{}
	err := c.BindJSON(&updateAuth)
	if err != nil {
		code := app.VerdictJsonErr(err)
		appG.ResponseI18nMsgSimple(code, err, nil)
		return
	}

	code, err := service.UpdateAuth(updateAuth)
	if err != nil {
		appG.ResponseI18nMsg(code, err, nil, nil, nil)
		return
	}

	appG.ResponseI18nMsg(e.StatusOK, err, nil, nil, nil)
}
