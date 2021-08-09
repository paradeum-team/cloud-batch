package clean

import (
	"cloud-batch/configs"
	"cloud-batch/pkg/utils"
	"fmt"
	"github.com/kataras/golog"
	"io/ioutil"
	"os"
	"time"
)

func DeleteDaysAgoFile(path string) error {
	// MaxAgeDay 天前
	daysAgo := utils.GetAfterTimeByYear(0, 0, -configs.LogConfig.MaxAgeDay)

	dirList, e := ioutil.ReadDir(path)
	if e != nil {
		return e
	}

	for _, fileInfo := range dirList {
		golog.Debugf("name: %s , time: %s", fileInfo.Name(), fileInfo.ModTime().Format(time.RFC3339))
		// 目录暂不删除
		if fileInfo.IsDir() {
			continue
		}
		if fileInfo.ModTime().Before(daysAgo) {
			filePath := fmt.Sprintf("%s/%s", path, fileInfo.Name())
			err := os.RemoveAll(filePath)
			if err != nil {
				return err
			}

			golog.Infof("Delete file %s !", filePath)
		}
	}

	return nil
}
