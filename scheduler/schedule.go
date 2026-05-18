package scheduler

import (
	"time"

	"github.com/robfig/cron/v3"
)

// EverySeconds 返回一个每 N 秒执行一次的调度器
func EverySeconds(secs ...int64) cron.Schedule {
	var second int64 = 1
	if len(secs) > 0 {
		second = secs[0]
	}
	return &intervalSchedule{
		interval: time.Duration(second) * time.Second,
	}
}

// EveryMinutes 返回一个每 N 分钟执行一次的调度器
func EveryMinutes(mins ...int64) cron.Schedule {
	var m int64 = 1
	if len(mins) > 0 {
		m = mins[0]
	}
	return &intervalSchedule{
		interval: time.Duration(m) * time.Minute,
	}
}

// EveryHours 返回一个每 N 小时执行一次的调度器
func EveryHours(hours ...int64) cron.Schedule {
	var hour int64 = 1
	if len(hours) > 0 {
		hour = hours[0]
	}
	return &intervalSchedule{
		interval: time.Duration(hour) * time.Hour,
	}
}

// EveryDuration 返回一个按指定时间间隔执行的调度器
func EveryDuration(interval time.Duration) cron.Schedule {
	return &intervalSchedule{interval: interval}
}

// DailyAt 返回一个每天指定时间执行的调度器
// hour: 0-23, minute: 0-59, second: 0-59
func DailyAt(hour, minute, second int) cron.Schedule {
	return &dailyAtSchedule{
		hour:   hour,
		minute: minute,
		second: second,
	}
}

// EveryMinutesAt 返回一个每隔N分钟的第0分执行的调度器
// interval: 间隔分钟数（如5表示每5分钟）
// 执行时间：0, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55
func EveryMinutesAt(interval int) cron.Schedule {
	return &minuteIntervalSchedule{
		interval: interval,
	}
}

// minuteIntervalSchedule 每隔N分钟的第0分执行的调度器
type minuteIntervalSchedule struct {
	interval int
}

func (s *minuteIntervalSchedule) Next(t time.Time) time.Time {
	minute := t.Minute()
	// 计算下一个满足条件的分钟
	remainder := minute % s.interval
	var nextMinute int
	if remainder == 0 && t.Second() == 0 && t.Nanosecond() == 0 {
		// 正好在整点，跳到下一个间隔
		nextMinute = minute + s.interval
	} else {
		nextMinute = minute + (s.interval - remainder)
	}

	if nextMinute >= 60 {
		// 下一个小时的第0分
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour()+1, 0, 0, 0, t.Location())
	}
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), nextMinute, 0, 0, t.Location())
}

// dailyAtSchedule 每天指定时间执行的调度器
type dailyAtSchedule struct {
	hour   int
	minute int
	second int
}

func (s *dailyAtSchedule) Next(t time.Time) time.Time {
	next := time.Date(t.Year(), t.Month(), t.Day(), s.hour, s.minute, s.second, 0, t.Location())
	if next.Before(t) || next.Equal(t) {
		next = next.AddDate(0, 0, 1)
	}
	return next
}

// intervalSchedule 按固定间隔执行的调度器
type intervalSchedule struct {
	interval time.Duration
}

func (s *intervalSchedule) Next(t time.Time) time.Time {
	return t.Add(s.interval)
}
