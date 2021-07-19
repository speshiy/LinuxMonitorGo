package cbackup

import (
	"log"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/speshiy/LinuxMonitorGo/_main/models/mbackup"
	"github.com/speshiy/LinuxMonitorGo/database"
	"github.com/speshiy/LinuxMonitorGo/settings"
)

//BackupDelete start backups delete
func BackupDelete() {
	log.Println("Backups delete cron started")
	var err error
	var DB *gorm.DB
	var backups []mbackup.Backup

	DB, err = database.OpenDatabase("backup_tuvis", "root", settings.DBRP, "UTC")
	if err != nil {
		return
	}
	defer DB.Close()

	err = mbackup.GetBackups(nil, DB, &backups)
	if err != nil {
		return
	}

	for _, backup := range backups {

		// Не надо удалять бэкапы баз, которые удалены
		if backup.IsDeletedDatabase {
			continue
		}

		expiresAt := time.Now()
		diff := expiresAt.Sub(backup.CreatedAt).Hours() / 24

		if diff > 6 {
			err = os.Remove(backup.BackupPath)
			if err != nil {
				continue
			}

			err = backup.Delete(nil, DB)
			if err != nil {
				continue
			}
		}
	}

	// Удаляем записи физически из БД
	DB.Exec("DELETE FROM s_backup WHERE deleted_at IS NOT NULL")

	log.Println("Backup delete cron finished")
}
