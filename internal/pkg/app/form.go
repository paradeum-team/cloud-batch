package app

import (
	"cloud-batch/internal/pkg/e"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// BindAndValid binds and validates data
func BindAndValid(c *gin.Context, form interface{}) (e.ErrorCode, error) {
	appG := Gin{C: c}
	err := c.Bind(form)
	if err != nil {
		return e.MalformedJSON, err
	}

	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		return e.InternalError, err
	}
	if !check {
		return e.InvalidParameter, appG.MargeErrorsToError(valid.Errors)
	}

	return e.StatusOK, nil
}
