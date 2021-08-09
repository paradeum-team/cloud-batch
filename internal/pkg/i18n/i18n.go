package i18n

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// 本地化包
var bundle *i18n.Bundle

func init() {
	// 设置默认语言
	bundle = i18n.NewBundle(language.Chinese)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// 加载公共信息语言包文件
	bundle.MustLoadMessageFile("assets/tomls/lang.en.toml")
	bundle.MustLoadMessageFile("assets/tomls/lang.zh-cn.toml")
}

/* 本地化语言
param:
	lang 语言 en|zh-cn
	messageID 语言文件中的 messageID
	errCode: string 错误码
	templateData: i18n 文件中需要变量替换的内容
	pluralCount: 传入 int 或 int64 类型数据 ， 根据数字判断是否返回复数格式 msg
*/
func MustLocalize(accept, messageID string, templateData interface{}, pluralCount interface{}) string {
	localizer := i18n.NewLocalizer(bundle, accept)
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
		PluralCount:  pluralCount,
	})
}
