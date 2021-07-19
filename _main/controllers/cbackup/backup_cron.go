package cbackup

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/speshiy/LinuxMonitorGo/database"
	"github.com/speshiy/LinuxMonitorGo/settings"
)

//BackupCron doing a backup
func BackupCron() {
	hour := time.Now().UTC().Hour()

	if hour < 18 && hour > 23 {
		log.Println("Dump is delayed because now is", hour, "and dumping is possible in time between 18:00 and 23:00 (UTC)")
		return
	}

	if settings.IsBackupInProgress {
		log.Println("Backup is progress")
		return
	}

	log.Println("Backup cron started")
	var err error
	var DBMain database.DBUserConnect
	var DB *gorm.DB

	err = DBMain.Open(settings.DBHostDump, "main", "rtu", settings.DBRTUP, "UTC")
	if err != nil {
		return
	}
	defer DBMain.Close()

	DB, err = database.OpenDatabase("backup_tuvis", "root", settings.DBRP, "UTC")
	if err != nil {
		return
	}

	defer DB.Close()

	settings.IsBackupInProgress = true

	err = DoBackup(nil, "", DBMain.Connect, DB, true, false)
	if err != nil {
		log.Println("Backup error " + err.Error())
	}

	settings.IsBackupInProgress = false

	log.Println("Backup cron finished")
}
