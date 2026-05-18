package x

import (
	"github.com/xframe-go/x/scheduler"
)

func RegisterScheduler(fn func(schedule *scheduler.Scheduler)) {
	cron := scheduler.New()

	fn(cron)

	cron.Start()

	AddShutdownListener(func() {
		ctx := cron.Stop()
		<-ctx.Done()
		Logger().Info("shutdown scheduler")
	})
}
