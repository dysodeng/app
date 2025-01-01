package helper

import "time"

// IsToday 是否为今天
func IsToday(t time.Time) bool {
	now := time.Now().Local()
	return t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day()
}

// IsSameYear 是否为同一年
func IsSameYear(t time.Time) bool {
	now := time.Now().Local()
	return t.Year() == now.Year()
}

// LastDayOfMonth 获取月份最后一天日期
func LastDayOfMonth(t time.Time) time.Time {
	t = t.Local()
	year := t.Year()
	month := int(t.Month())
	lastDay := time.Date(year, time.Month(month+1), 0, 23, 59, 59, 0, time.Local).Add(-time.Nanosecond)
	return lastDay
}

// MonthDays 获取月份天数
func MonthDays(t time.Time) int {
	year := t.Year()
	month := int(t.Month())
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 4, 6, 9, 11:
		return 30
	case 2:
		if ((year%4) == 0 && (year%100) != 0) || (year%400) == 0 {
			return 29
		}
		return 28
	default:
		panic("无效的月份")
	}
}

// WeekDay 获取星期几
func WeekDay(t time.Time) uint8 {
	week := t.Weekday()
	switch week {
	case time.Sunday:
		return 7
	case time.Monday:
		return 1
	case time.Tuesday:
		return 2
	case time.Wednesday:
		return 3
	case time.Thursday:
		return 4
	case time.Friday:
		return 5
	case time.Saturday:
		return 6
	}
	return 0
}

// WeekChineseDay 获取星期几中文名称
func WeekChineseDay(t time.Time) string {
	week := WeekDay(t)
	switch week {
	case 1:
		return "周一"
	case 2:
		return "周二"
	case 3:
		return "周三"
	case 4:
		return "周四"
	case 5:
		return "周五"
	case 6:
		return "周六"
	case 7:
		return "周日"
	}
	return ""
}

func WeekChinese(t time.Time) string {
	weekday := t.Weekday()
	var cnWeekday string

	switch weekday {
	case time.Monday:
		cnWeekday = "星期一"
	case time.Tuesday:
		cnWeekday = "星期二"
	case time.Wednesday:
		cnWeekday = "星期三"
	case time.Thursday:
		cnWeekday = "星期四"
	case time.Friday:
		cnWeekday = "星期五"
	case time.Saturday:
		cnWeekday = "星期六"
	case time.Sunday:
		cnWeekday = "星期日"
	}

	return cnWeekday
}

// DurationTime 消耗时间计算
// 用法 defer DurationTime(time.Now()) 可计算函数执行时间
func DurationTime(pre time.Time) time.Duration {
	return time.Since(pre)
}

// ResolveTime 将整数转换为时分秒
// @param int seconds 秒数
// @return int hour 小时数
// @return int minute 分钟数
// @return int second 秒数
func ResolveTime(seconds int) (hour, minute, second int) {
	hour = seconds / 3600
	minute = (seconds - hour*3600) / 60
	second = seconds - hour*3600 - minute*60
	return
}
