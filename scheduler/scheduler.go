package scheduler

import "github.com/robfig/cron/v3"

type Scheduler struct {
	*cron.Cron
}

func New() *Scheduler {
	c := cron.New(
		cron.WithSeconds(),
		cron.WithChain(
			cron.SkipIfStillRunning(cron.DefaultLogger),
			cron.Recover(cron.DefaultLogger),
		),
	)

	return &Scheduler{c}
}
