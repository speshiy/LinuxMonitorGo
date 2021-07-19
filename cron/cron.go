package cron

import (
	"github.com/robfig/cron"
	"github.com/speshiy/LinuxMonitorGo/_monitoring/controllers/cmonitor"
)

//InitCron init periodic backups
func InitCron() {
	c := cron.New()

	// Мониторинг
	c.AddFunc("@every 5s", func() {
		cmonitor.PostRAM()
	})

	c.Start()
}
