package app

import (
	"bytes"
	"cloud-batch/configs"
	"cloud-batch/internal/pkg/e"
	"cloud-batch/internal/pkg/i18n"
	"cloud-batch/internal/pkg/logging"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/validation"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/kataras/golog"
	"github.com/pkg/errors"
	"runtime"
	"strings"
)

type Gin struct {
	C      *gin.Context
	logger *golog.Logger
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type ResponseString struct {
	Code      string      `json:"code"`
	Msg       string      `json:"msg"`
	ErrDetail interface{} `json:"errDetail"`
	Data      interface{} `json:"data"`
	RequestId string      `json:"requestId"`
}

func VerdictJsonErr(err error) e.ErrorCode {
	var jsonErr *json.UnmarshalTypeError
	if errors.As(err, &jsonErr) {
		return e.MalformedJSON
	} else {
		return e.InappropriateJSON
	}
}

func (g *Gin) GetLogger() *golog.Logger {
	if g.logger == nil {
		g.logger = logging.New().SetPrefix(fmt.Sprintf("%s requestId: %s, ", configs.Server.HostName, requestid.Get(g.C)))
	}
	return g.logger
}

func (g *Gin) MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		g.GetLogger().Info(err.Key, err.Message)
	}
	return
}

func (g *Gin) MargeErrorsToError(errors []*validation.Error) error {
	var buffer bytes.Buffer
	for i, err := range errors {
		buffer.WriteString(err.Key)
		buffer.WriteString(": ")
		buffer.WriteString(err.Message)
		if i != len(errors)-1 {
			buffer.WriteString(", ")
		}
	}
	return fmt.Errorf("%s", buffer.String())
}

//
//func (g *Gin) Response(httpCode, errCode int, data interface{}) {
//	if errCode != e.SUCCESS && errCode != e.OK {
//		logger.Warnf("Code: %d, Msg: %s, Data: %v", errCode, e.GetMsg(errCode), data)
//	}
//	g.C.JSON(httpCode, Response{
//		Code: errCode,
//		Msg:  e.GetMsg(errCode),
//		Data: data,
//	})
//	return
//}

/*
	简单 response i18n msg
	param:
		errCode: string 错误码
		err: 详细error信息
		data: 返回的数据 object
*/
func (g *Gin) ResponseI18nMsgSimple(errCode e.ErrorCode, err error, data interface{}) {
	g.ResponseI18nMsg(errCode, err, nil, nil, data)
}

/*
	response i18n msg
	param:
		errCode: string 错误码
		err: 详细error信息
		templateData: i18n 文件中需要变量替换的内容
		pluralCount: 传入 int 或 int64 类型数据 ， 根据数字判断是否返回复数格式 msg
		data: 返回的数据 object
*/
func (g *Gin) ResponseI18nMsg(errCode e.ErrorCode, err error, templateData interface{}, pluralCount interface{}, data interface{}) {
	var msg string

	log := g.GetLogger()

	if errCode.GetCode() != e.StatusOK.GetCode() {
		// 根据 errCode 获取 i18n 报错信息
		accept := g.C.Request.Header.Get("Accept-Language")

		if templateData != nil {
			errCode.TemplateData = templateData
		}

		if pluralCount != nil {
			errCode.PluralCount = pluralCount
		}

		msg = i18n.MustLocalize(accept, errCode.GetCode(), errCode.TemplateData, errCode.PluralCount)
	}

	// 获取报错的具体文件、函数名、行数
	_, file, line, _ := runtime.Caller(2)

	// 隐藏内部错误 和 Auth 错误, 提示非内部错误
	if errCode.GetCode() == e.InternalError.GetCode() || errCode.GetCode() == e.AuthError.GetCode() {
		// 内部错误
		log.Errorf("%s:%d, errCode: %s, err: %+v", file, line, errCode.GetCode(), err)
		g.C.JSON(errCode.GetHttpCode(), ResponseString{
			Code:      errCode.GetCode(),
			Msg:       msg,
			Data:      nil,
			RequestId: requestid.Get(g.C),
		})
	} else if errCode.GetCode() == e.StatusOK.GetCode() {
		// 正常无错误
		g.C.JSON(errCode.GetHttpCode(), ResponseString{
			Code:      errCode.GetCode(),
			Msg:       msg,
			Data:      data,
			RequestId: requestid.Get(g.C),
		})
	} else {
		// 其它错误
		log.Infof("%s:%d, errCode: %s, err: %v", file, line, errCode.GetCode(), err)
		var errDetail string
		if err == nil {
			log.Info("errDetail is nil")
		} else if !strings.HasSuffix(errCode.GetCode(), "Limit") || errCode.GetCode() == e.InappropriateJSON.GetCode() {
			// 除 Limit 相关错误 详细内容返回给前端 或 如果错误是 json 值校验不通过，则把详细内容返回给前端
			errDetail = err.Error()
		}

		g.C.JSON(errCode.GetHttpCode(), ResponseString{
			Code:      errCode.GetCode(),
			Msg:       msg,
			ErrDetail: errDetail,
			Data:      data,
			RequestId: requestid.Get(g.C),
		})
	}

	return
}

//func (g *Gin) ResponseMsg(httpCode, errCode int, msg string, data interface{}) {
//	if errCode != e.SUCCESS && errCode != e.OK {
//		logger.Warnf("Code: %d, Msg: %s, Data: %v", errCode, msg, data)
//	}
//	g.C.JSON(httpCode, Response{
//		Code: errCode,
//		Msg:  msg,
//		Data: data,
//	})
//	return
//}
