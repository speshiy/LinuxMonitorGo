package cron

import (
	"github.com/robfig/cron"
	"github.com/speshiy/LinuxMonitorGo/_monitoring/controllers/cmonitor"
)

//InitCron init periodic deals
func InitCron() {
	c := cron.New()

	c.AddFunc("@every 10s", func() {
		cmonitor.PostRAM()
	})

	c.Start()
}
