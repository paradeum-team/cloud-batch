package utils

import (
	"fmt"
	"github.com/kataras/golog"
	"runtime"
	"time"
)

var bfsTimeTemplate = "20060102150405"
var timeDay = "20060102"
var NanoTimeStamp = "20060102150405.999999999"

//获取格式 20060102 日期
func GetDayStr() string {
	return time.Now().Format(timeDay)
}

func NowNanoTimeStamp() string {
	return time.Now().Format(NanoTimeStamp)
}

func ConvertBfsTimeToTime(bfsTime string) (time.Time, error) {
	timeStamp, err := time.ParseInLocation(bfsTimeTemplate, bfsTime, time.Local)
	if err != nil {
		return time.Time{}, err
	}
	return timeStamp, nil
}

func ConvertRFC3339ToTime(standardTime string) (time.Time, error) {
	timeStamp, err := time.ParseInLocation(time.RFC3339, standardTime, time.Local)
	if err != nil {
		return time.Time{}, err
	}
	return timeStamp, nil
}

func ConvertTimeToBfsTime(_time time.Time) string {
	return _time.Format(bfsTimeTemplate)
}

func ConvertBfsTimeToRFC3339(bfsTime string) (string, error) {
	timeStamp, err := time.ParseInLocation(bfsTimeTemplate, bfsTime, time.Local)
	if err != nil {
		return "", err
	}
	return timeStamp.Format(time.RFC3339), nil
}

func ConvertRFC3339ToBfsTime(standardTime string) (string, error) {
	timeStamp, err := time.ParseInLocation(time.RFC3339, standardTime, time.Local)
	if err != nil {
		return "", err
	}
	return timeStamp.Format(bfsTimeTemplate), nil
}

func StrTimeToTTL(strTime string) (time.Duration, error) {
	timeStamp, err := time.ParseInLocation(bfsTimeTemplate, strTime, time.Local)
	if err != nil {
		return 0, err
	}
	ttl := timeStamp.Sub(time.Now())
	return ttl, nil
}

func IsExpire(expireTime string) bool {
	//1d 1m 1y
	expireTimeStamp, _ := time.ParseInLocation(bfsTimeTemplate, expireTime, time.Local)

	now := time.Now()
	return now.After(expireTimeStamp)
}

func AfterYearDays(year int) int {
	today := time.Now()
	yearLater := today.AddDate(year, 0, -1)
	return int(yearLater.Sub(today).Hours()) / 24
}

func KeepYear(currentExpire string, years, months, days int) (bool, string) {
	// years年 months月 days后 天
	aYearLater := GetAfterTimeByYear(years, months, days)
	if currentExpire != "" {
		// 输入的当时时间格式化
		expireTimeStamp, _ := time.ParseInLocation(bfsTimeTemplate, currentExpire, time.Local)
		// 如果当时时间足够，返回false
		if expireTimeStamp.After(aYearLater) {
			return false, ""
		}
	}
	// 不足返回新时间,固定小时分钟秒 235959
	return true, fmt.Sprintf("%s235959", aYearLater.Format(timeDay))

	//current := expireTimeStamp.UTC().Truncate(24 * time.Hour)
	//now := time.Now().UTC().Truncate(24 * time.Hour)
	//days := int(current.Sub(now).Hours() / 24)
	//if days < 365 {
	//	dd, _ := time.ParseDuration("24h")
	//	expireTime := time.Now().Add(dd * 365)
	//	return true, expireTime.Format(bfsTimeTemplate)
	//}
}

func GetAfterTimestampByYear(years, months, days int) string {
	return GetAfterTimeByYear(years, months, days).Format(bfsTimeTemplate)
}

func GetAfterTimeByYear(years int, months int, days int) time.Time {
	return time.Now().AddDate(years, months, days)
}

// 打印函数名称行数及运行时间
// 注意，是对 TimeCost()返回的函数进行调用，因此调用时需要加两对小括号,
// 参考：https://blog.csdn.net/K346K346/article/details/92673425， https://colobu.com/2018/11/03/get-function-name-in-go/
func TimeCost() func() {
	pc, file, line, _ := runtime.Caller(1)

	start := time.Now()
	return func() {
		tc := time.Since(start)
		golog.Infof("%s:%d %s, time cost = %v\n", file, line, runtime.FuncForPC(pc).Name(), tc)
	}
}
