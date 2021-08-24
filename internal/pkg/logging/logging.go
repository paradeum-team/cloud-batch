package logging

import (
	"cloud-batch/configs"
	"cloud-batch/internal/pkg/clean"
	"fmt"
	"github.com/gogf/gf/os/gfile"
	"github.com/kataras/golog"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/robfig/cron"
	"log"
	"time"
)

var Logger *golog.Logger

// set default golog
func init() {
	var err error
	if err != nil {
		log.Fatalf("%+v", err)
	}
	// init logger
	if Logger == nil {
		Logger = golog.New()
		initGolog(Logger, configs.LogConfig.AppName)
	}
	// 定时清理 log 任务
	go Cron()
}

// new golog
func New() *golog.Logger {
	return initGolog(golog.New(), configs.LogConfig.AppName)
}

func initGolog(log *golog.Logger, logPrefix string) *golog.Logger {
	log.SetLevel(configs.LogConfig.Level)
	log.SetTimeFormat(configs.LogConfig.TimeFormat)

	log.SetPrefix(fmt.Sprintf("%s ", configs.Server.HostName))
	if configs.LogConfig.IsOutPutFile == false {
		return log
	}
	w := GetRotatelogsWriter(logPrefix)
	log.SetOutput(w)
	return log
}

func Cron() {

	golog.Info("Clean Starting...")
	var err error
	// 启动时执行一次
	// 清理上传异常临时文件
	//err := clean.DeleteDaysAgoFile(fmt.Sprintf("%s/%s", configs.Server.RuntimeRootPath, configs.Server.MediaSavePath))
	//if err != nil {
	//	golog.Errorf("DeleteDaysAgoFile err: %+v", err)
	//}
	// 清理日志文件
	err = clean.DeleteDaysAgoFile(fmt.Sprintf("%s/%s", configs.Server.RuntimeRootPath, configs.LogConfig.Path))
	if err != nil {
		golog.Errorf("DeleteDaysAgoFile err: %+v", err)
	}

	c := cron.New()
	c.AddFunc("0 0 0 */1 * *", func() {
		golog.Info("Clean ...")
		// 定时清理上传异常临时文件
		//err = clean.DeleteDaysAgoFile(fmt.Sprintf("%s/%s", configs.Server.RuntimeRootPath, configs.Server.MediaSavePath))
		//if err != nil {
		//	golog.Warnf("DeleteDaysAgoFile err: %v", err)
		//}
		// 定时清理日志文件
		err = clean.DeleteDaysAgoFile(fmt.Sprintf("%s/%s", configs.Server.RuntimeRootPath, configs.LogConfig.Path))
		if err != nil {
			golog.Warnf("DeleteDaysAgoFile err: %v", err)
		}
	})

	c.Start()

	select {}
}

func GetRotatelogsWriter(prefix string) *rotatelogs.RotateLogs {
	logDir := fmt.Sprintf("%s/%s", configs.Server.RuntimeRootPath, configs.LogConfig.Path)
	ok := gfile.Exists(logDir)
	if !ok {
		err := gfile.Mkdir(logDir)
		if err != nil {
			log.Fatal(err)
		}
	}
	ok = gfile.IsWritable(logDir)
	if !ok {
		log.Fatalf("No permission to write %s", logDir)
	}

	logInfoPath := fmt.Sprintf("%s/%s.%s", logDir, prefix, configs.Server.HostName) + ".%Y%m%d.log"
	w, err := rotatelogs.New(
		logInfoPath,
		rotatelogs.WithMaxAge(time.Duration(configs.LogConfig.MaxAgeDay*24)*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour))
	if err != nil {
		log.Fatalf("%+v", err)
	}
	return w
}
