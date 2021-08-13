package v1

import (
	"cloud-batch/internal/models"
	"cloud-batch/internal/pkg/app"
	"cloud-batch/internal/pkg/e"
	"cloud-batch/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Login
// @Summary post Auth
// @Tags Auth
// @Produce  json
// @Param username formData string true "username"
// @Param password formData string true "password" format(password)
// @Success 200 {object} app.ResponseString
// @Failure 400 {object} app.ResponseString
// @Failure 500 {object} app.ResponseString
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		auth = models.Auth{}
		valid = validator.New()
	)

	auth.Username = c.PostForm("username")
	auth.Password = c.PostForm("password")

	//authJson,err := json.Marshal(auth)
	//if err != nil {
	//	appG.ResponseI18nMsgSimple(e.InvalidParameter, err, nil)
	//	return
	//}

	err := valid.Struct(auth)

	if err != nil {
		appG.ResponseI18nMsgSimple(e.InvalidParameter, err, nil)
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
// @Param username formData string true "username"
// @Param password formData string true "password" format(password)
// @Param old_password formData string true "old password" format(password)
// @Success 200 {object} app.ResponseString
// @Failure 400,401 {object} app.ResponseString
// @Failure 500 {object} app.ResponseString
// @Router /auth/passwd [put]
func UpdateAuth(c *gin.Context) {

	var(
		appG = app.Gin{C: c}
		updateAuth = models.UpdateAuth{}
		valid = validator.New()
	)

	updateAuth.Username = c.PostForm("username")
	updateAuth.Password = c.PostForm("password")
	updateAuth.OldPassword = c.PostForm("old_password")

	err := valid.Struct(updateAuth)
	if err != nil {
		appG.ResponseI18nMsgSimple(e.InvalidParameter, err, nil)
		return
	}

	code, err := service.UpdateAuth(updateAuth)
	if err != nil {
		appG.ResponseI18nMsg(code, err, nil, nil, nil)
		return
	}

	appG.ResponseI18nMsg(e.StatusOK, err, nil, nil, nil)
}
