package app

import (
	"cloud-batch/internal/pkg/logging"
	"github.com/astaxie/beego/validation"
)

func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		logging.Logger.Info(err.Key, err.Message)
	}

	return
}
