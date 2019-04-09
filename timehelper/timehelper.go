package timehelper

import (
	"fmt"
	"time"
)

type format struct {
	DateTime string
	Date     string
	Time     string
}

var Format *format

func init() {
	Format = &format{}
	Format.DateTime = "2006-01-02 15:04:05"
	Format.Date = "2006-01-02"
	Format.Time = "15:04:05"
}

//取当前时间的unix timestamp
func Timestamp() int64 {
	return time.Now().Unix()
}

//取标准时间
func UtcDate() string {
	t := time.Now().UTC()
	format := "%d-%0.2d-%0.2dT%0.2d:%0.2d:%0.2d.012Z"
	return fmt.Sprintf(format, t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

//格式化unix timestamp
func Date(timestamp int64, args ...string) string {
	layout := Format.DateTime
	if args != nil {
		layout = args[0]
	}
	return time.Unix(timestamp, 0).Format(layout)
}

//格式化当前时间
func Today(args ...string) string {
	layout := Format.DateTime
	if args != nil {
		layout = args[0]
	}
	return time.Now().Format(layout)
}

//将日期转换为unix timestamp
func ToTimestamp(date, format string) int64 {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t, _ := time.ParseInLocation(format, date, loc)
	return t.Unix()
}

//取每月的天数
func MonthDays(year, month int) int {
	days := 0
	switch month {
	case 4:
		fallthrough
	case 6:
		fallthrough
	case 9:
		fallthrough
	case 11:
		days = 30
	case 2:
		if ((year%4) == 0 && (year%100) != 0) || (year%400) == 0 {
			days = 29
		} else {
			days = 28
		}
	default:
		days = 31
	}

	return days
}

//将日期加上一个时间段
func DateAddDuration(day string, duration time.Duration, args ...string) string {
	format := Format.Date
	if len(args) > 0 {
		format = args[0]
	}

	t, _ := time.Parse(format, day)
	return t.Add(duration).Format(format)
}

//将日期加上多少月
func AddMonths(day string, months int, args ...string) string {
	return Add(day, 0, months, 0, args...)
}

//将日期加上多少天
func AddDays(day string, days int, args ...string) string {
	return Add(day, 0, 0, days, args...)
}

//将日期加上多少年
func AddYears(day string, years int, args ...string) string {
	return Add(day, years, 0, 0, args...)
}

//将日期加上年月日
func Add(day string, years, months, days int, args ...string) string {
	format := Format.Date
	if len(args) > 0 {
		format = args[0]
	}

	t, _ := time.Parse(format, day)
	return t.AddDate(years, months, days).Format(format)
}

//计算两个日期间隔的天数
func DateDays(startDate, endDate, format string) int {
	startTime, _ := time.Parse(format, startDate)
	endTime, _ := time.Parse(format, endDate)
	return int(endTime.Sub(startTime).Hours()/24) + 1
}

//计算unix timestamp当天的起始unix timestamp
func DayStartTimestamp(timestamp int64) int64 {
	s := time.Unix(timestamp, 0).Format(Format.Date)
	return ToTimestamp(s, Format.Date)
}