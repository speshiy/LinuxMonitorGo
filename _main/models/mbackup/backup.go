package mbackup

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/speshiy/LinuxMonitorGo/settings"
)

//Backup structure
type Backup struct {
	gorm.Model
	Host              string  `gorm:"column:host;type:varchar(250);not null;default:''" json:"Host"`
	DBName            string  `gorm:"column:db_name;type:varchar(250);not null;default:''" json:"DBName"`
	UserName          string  `gorm:"column:user_name;type:varchar(250);not null;default:''" json:"UserName"`
	BackupName        string  `gorm:"column:backup_name;type:varchar(2000);not null;default:''" json:"BackupName"`
	BackupPath        string  `gorm:"column:backup_path;type:varchar(2000);not null;default:''" json:"BackupPath"`
	FileSize          float32 `gorm:"column:file_size;" json:"FileSize"`
	IsDeletedDatabase bool    `gorm:"column:is_deleted_database;default: 0" json:"IsDeletedDatabase"`
}

//TableName return new table name for User struct
func (Backup) TableName() string {
	return "s_backup"
}

//GetBackups return backup by DBName
func GetBackups(c *gin.Context, DB *gorm.DB, backups *[]Backup) error {
	if DB == nil {
		DB = c.MustGet("DB").(*gorm.DB)
	}

	r := DB.Order("created_at desc").Find(&backups)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

//GetBackupsByDBName return backup by DBName
func GetBackupsByDBName(c *gin.Context, dbName string, b *[]Backup, DB *gorm.DB) error {
	if DB == nil {
		DB = c.MustGet("DB").(*gorm.DB)
	}
	var r *gorm.DB

	if dbName == "" {
		r = DB.Limit(5).Where("db_name IN (?, ?)", "main", "client").Order("created_at desc").Find(&b)
	} else {
		r = DB.Limit(5).Where("db_name = ?", dbName).Order("created_at desc").Find(&b)
	}

	if r.Error != nil {
		return r.Error
	}

	return nil
}

//GetByID return backup by id
func (b *Backup) GetByID(c *gin.Context, DB *gorm.DB) error {
	if DB == nil {
		DB = c.MustGet("DB").(*gorm.DB)
	}

	r := DB.First(&b)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

//GetLastByDBName return backup by name
func (b *Backup) GetLastByDBName(c *gin.Context, DB *gorm.DB) error {
	if DB == nil {
		DB = c.MustGet("DB").(*gorm.DB)
	}

	r := DB.Order("id DESC").Where("db_name = ?", b.DBName).First(&b)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

//Post backup
func (b *Backup) Post(c *gin.Context, DB *gorm.DB) error {
	if DB == nil {
		DB = c.MustGet("DB").(*gorm.DB)
	}

	r := DB.Create(&b)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

//Put backup
func (b *Backup) Put(c *gin.Context, DB *gorm.DB) error {
	if DB == nil {
		DB = c.MustGet("DB").(*gorm.DB)
	}

	r := DB.Model(&b).Updates(&b)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

//Delete backup
func (b *Backup) Delete(c *gin.Context, DB *gorm.DB) error {
	if DB == nil {
		DB = c.MustGet("DB").(*gorm.DB)
	}

	r := DB.Delete(&b)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

//IsNeedBackup check that this DB is need to backup
func IsNeedBackup(c *gin.Context, dbName string, DB *gorm.DB) (bool, error) {
	if DB == nil {
		DB = c.MustGet("DB").(*gorm.DB)
	}
	var err error
	var sql string
	var diffHours uint

	sql = `SELECT 
				HOUR(TIMEDIFF(t.created_at, UTC_TIMESTAMP())) as DiffHours 
			FROM 
			(SELECT 
					b.db_name, 
					max(b.created_at) as created_at 
			FROM s_backup as b 
			WHERE b.db_name = ? 
			GROUP BY 
				b.db_name ) as t 
			ORDER BY 
			t.created_at`

	row := DB.Raw(sql, dbName).Select("DiffHours").Row()
	err = row.Scan(&diffHours)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return true, nil
		}
		return false, err
	}

	// Если это системная база, то её бэкапируем каждые 24 часа
	// if dbName == "main" || dbName == "client" {
	// 	if diffHours <= 36 {
	// 		return true, nil
	// 	}
	// }

	// Проверяем истёкло ли время запрета бэкапа
	backupIntervalHours, _ := strconv.Atoi(settings.BackupIntervalHours)
	if diffHours < uint(backupIntervalHours) {
		return false, nil
	}

	return true, nil
}
