package tools

import (
	"strconv"
	"strings"
	"time"
)

func S2i64(s string) (int64, error) {
	if s == "" {
		return int64(0), nil
	}
	return strconv.ParseInt(s, 10, 64)
}

func S2i(s string) (int, error) {
	if s == "" {
		return 0, nil
	}
	return strconv.Atoi(s)
}

func I642s(i int64) string {
	return strconv.FormatInt(i, 10)
}

func U322s(u uint32) string {
	return strconv.FormatUint(uint64(u), 10)
}

func U642s(u uint64) string {
	return strconv.FormatUint(u, 10)
}

//根据时间字符串的后缀，将时间字符串转换成int64的毫秒
const (
	Millisecond = 1
	Second      = 1000 * Millisecond
	Minute      = 60 * Second
	Hour        = 60 * Minute
	Day         = 24 * Hour
)

func CoversionTimeBySuffix2Millisecond(timeString string) int64 {
	var num int64
	var unit string
	for index, value := range strings.Split(timeString, "") {
		if n, err := strconv.ParseInt(value, 10, 64); err == nil {
			num = num*10 + n
		} else {
			unit += string(timeString[index])
		}
	}
	switch unit {
	case "":
		return num
	case "ms":
		return num
	case "s":
		return num * Second
	case "m":
		return num * Minute
	case "h":
		return num * Hour
	case "d":
		return num * Day
	default:
		return 0
	}
}

func CoversionInt642TimeDuration(timeStamp int64) time.Duration {
	return time.Duration(timeStamp) * time.Millisecond
}

func CoversionTimeBySuffix2TimeDuration(timeString string) time.Duration {
	return CoversionInt642TimeDuration(CoversionTimeBySuffix2Millisecond(timeString))
}

func ConverTimeStamp2Time(timeStamp int64) time.Time {
	return time.Unix(timeStamp/1000, timeStamp%1000)
}
