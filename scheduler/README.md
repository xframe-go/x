- EverySeconds(secs ...int64) - 每 N 秒执行一次（默认 1 秒）
- EveryMinutes(mins ...int64) - 每 N 分钟执行一次（默认 1 分钟）
- EveryHours(hours ...int64) - 每 N 小时执行一次（默认 1 小时）
- EveryDuration(interval) - 按指定时间间隔执行
- DailyAt(hour, minute, second) - 每天指定时间执行


# 使用示例：
c := cron.New()

c.Schedule(EverySeconds(5), job)      // 每 5 秒
c.Schedule(EveryMinutes(10), job)     // 每 10 分钟
c.Schedule(EveryHours(2), job)        // 每 2 小时
c.Schedule(EveryDuration(30*time.Second), job)  // 每 30 秒
c.Schedule(DailyAt(9, 30, 0), job)    // 每天 9:30:00

c.Schedule(EveryMinutesAt(5), job)   // 每5分钟的0秒执行（23:00, 23:05, 23:10...）
c.Schedule(EveryMinutesAt(10), job)  // 每10分钟的0秒执行（23:00, 23:10, 23:20...）
c.Schedule(EveryMinutesAt(15), job)  // 每15分钟的0秒执行（23:00, 23:15, 23:30...）