package utils

import (
	"fmt"
	"testing"
	"time"
)

func timeSub(t1, t2 time.Time) int {
	t1 = t1.UTC().Truncate(24 * time.Hour)
	t2 = t2.UTC().Truncate(24 * time.Hour)
	return int(t1.Sub(t2).Hours() / 24)
}

func TestTime(t *testing.T) {
	layout := "2006-01-02 15:04:05"

	// just one second
	t1, _ := time.Parse(layout, "2007-01-02 23:59:59")
	t2, _ := time.Parse(layout, "2007-01-03 00:00:00")
	if timeSub(t2, t1) != 1 {
		panic("one second but different day should return 1")
	}

	// just one day
	t1, _ = time.Parse(layout, "2007-01-02 23:59:59")
	t2, _ = time.Parse(layout, "2007-01-03 23:59:59")
	if timeSub(t2, t1) != 1 {
		panic("just one day should return 1")
	}

	t1, _ = time.Parse(layout, "2017-09-01 10:00:00")
	t2, _ = time.Parse(layout, "2017-09-02 11:00:00")
	if timeSub(t2, t1) != 1 {
		panic("just one day should return 1")
	}

	// more than one day
	t1, _ = time.Parse(layout, "2007-01-02 23:59:59")
	t2, _ = time.Parse(layout, "2007-01-04 00:00:00")
	if timeSub(t2, t1) != 2 {
		panic("just one day should return 2")
	}
	// just 3 day
	t1, _ = time.Parse(layout, "2007-01-02 00:00:00")
	t2, _ = time.Parse(layout, "2007-01-05 00:00:00")
	if timeSub(t2, t1) != 3 {
		panic("just 3 day should return 3")
	}

	// different month
	t1, _ = time.Parse(layout, "2007-01-02 00:00:00")
	t2, _ = time.Parse(layout, "2007-02-02 00:00:00")
	if timeSub(t2, t1) != 31 {
		fmt.Println(timeSub(t2, t1))
		panic("just one month:31 days should return 31")
	}

	// 29 days in 2mth
	t1, _ = time.Parse(layout, "2000-02-01 00:00:00")
	t2, _ = time.Parse(layout, "2000-03-01 00:00:00")
	if timeSub(t2, t1) != 29 {
		fmt.Println(timeSub(t2, t1))
		panic("just one month:29 days should return 29")
	}

	fmt.Println(KeepYear("20210529000000", 1, 0, -1))
	fmt.Println(KeepYear("20300428000000", 9, 11, 0))
	sed, err := StrTimeToTTL("20200630152500")
	fmt.Println(sed)
	fmt.Println(err)
}
