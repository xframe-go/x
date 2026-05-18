package cron

import "github.com/robfig/cron/v3"

func New() {
	c := cron.New(cron.WithSeconds())

	c.Run()
}
