package cron

import (
	"github.com/robfig/cron"
	"github.com/speshiy/LinuxMonitorGo/_main/controllers/cbackup"
	"github.com/speshiy/LinuxMonitorGo/_monitoring/controllers/cmonitor"
)

//InitCron init periodic backups
func InitCron() {
	c := cron.New()

	// Мониторинг
	c.AddFunc("@every 5s", func() {
		cmonitor.PostRAM()
	})

	c.AddFunc("@every 15m", func() {
		cbackup.BackupCron()
	})

	c.AddFunc("@every 24h", func() {
		cbackup.BackupDelete()
	})
	c.Start()

	go cbackup.BackupCron()
	go cbackup.BackupDelete()
}
