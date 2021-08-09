package middleware

import (
	"cloud-batch/configs"
	"cloud-batch/internal/pkg/logging"
	"fmt"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
	"time"
)

/**
日志中间件
*/

// defaultLogFormatter is the default log format function Logger middleware uses.
var defaultLogFormatter = func(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency = param.Latency - param.Latency%time.Second
	}
	return fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)
}

// access 日志输出到文件，并添加自定义字段
func LoggerToFile(appName string) gin.HandlerFunc {
	var out io.Writer
	if configs.LogConfig.IsOutPutFile == true {
		out = logging.GetRotatelogsWriter("access")
	} else {
		out = gin.DefaultWriter
	}

	conf := gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			hostname, err := os.Hostname()
			if err != nil {
				log.Fatal(err)
			}
			// goaccess 日志格式 w3c: [appName] [%^] %h - [%d %t] %m %U %H %s %D %b %^ %u
			// goaccess 分析工具参考：https://www.goaccess.cc/?mod=man
			return fmt.Sprintf("[%s] [%s] %s - [%s] %s %s %s %d %d %d %s \"%s\" %s\n",
				appName,
				hostname,
				param.ClientIP,
				param.TimeStamp.Format("2006-01-02 15:04:05"),
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency.Microseconds(),
				param.BodySize,
				param.Keys["X-Request-ID"],
				param.Request.UserAgent(),
				param.ErrorMessage,
			)
		},
		Output:    out,
		SkipPaths: nil,
	}
	return LoggerWithConfig(conf)
}

// LoggerWithConfig instance a Logger middleware with config.
func LoggerWithConfig(conf gin.LoggerConfig) gin.HandlerFunc {
	formatter := conf.Formatter
	if formatter == nil {
		formatter = defaultLogFormatter
	}

	out := conf.Output
	if out == nil {
		out = gin.DefaultWriter
	}

	notlogged := conf.SkipPaths

	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		requestid := requestid.Get(c)

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			param := gin.LogFormatterParams{
				Request: c.Request,
				Keys:    c.Keys,
			}

			// Stop timer
			param.TimeStamp = time.Now()
			param.Latency = param.TimeStamp.Sub(start)

			param.ClientIP = c.ClientIP()
			param.Method = c.Request.Method
			param.StatusCode = c.Writer.Status()
			param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

			param.BodySize = c.Writer.Size()
			if param.Keys == nil {
				param.Keys = make(map[string]interface{})
			}
			param.Keys["X-Request-ID"] = requestid

			if raw != "" {
				path = path + "?" + raw
			}

			param.Path = path

			fmt.Fprint(out, formatter(param))
		}
	}
}
