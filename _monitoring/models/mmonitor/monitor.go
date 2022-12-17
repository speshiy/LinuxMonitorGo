package mmonitor

import (
	"time"

	"github.com/jinzhu/gorm"
)

//Monitor stucture
type Monitor struct {
	gorm.Model
	IP               string `gorm:"column:ip" json:"IP"`
	Server           string `gorm:"column:server" json:"Server"`
	Disk             uint64 `gorm:"column:disk;default:null" json:"Disk"`
	DiskUsed         uint64 `gorm:"column:disk_used;default:null" json:"DiskUsed"`
	RAM              uint64 `gorm:"column:ram;default:null" json:"RAM"`
	RAMUsed          uint64 `gorm:"column:ram_used;default:null" json:"RAMUsed"`
	Swap             uint64 `gorm:"column:swap;default:null" json:"Swap"`
	SwapUsed         uint64 `gorm:"column:swap_used;default:null" json:"SwapUsed"`
	Sort             uint   `gorm:"column:sort;default:null" json:"Sort"`
	IsDatabaseServer bool   `gorm:"column:is_database_server;default:0" json:"IsDatabaseServer"`
	IsBackupServer   bool   `gorm:"column:is_backup_server;default:0" json:"IsBackupServer"`
	IsAlive          bool   `gorm:"-" json:"IsAlive"`
}

//TableName return new table name
func (Monitor) TableName() string {
	return "sys_monitoring"
}

//GetMonitors into DB
func GetMonitors(DBMonitor *gorm.DB, monitors *[]Monitor) error {
	r := DBMonitor.Order("sort").Find(&monitors)
	if r.Error != nil && !gorm.IsRecordNotFoundError(r.Error) {
		return r.Error
	}

	// Если сервер не обновлял о себе инфу более 3 минут, то он мёртв
	for idx, monitor := range *monitors {
		(*monitors)[idx].IsAlive = time.Since(monitor.UpdatedAt).Minutes() <= 3
	}

	return nil
}

//GetByServer into DB
func (m *Monitor) GetByServer(DBMonitor *gorm.DB) error {
	r := DBMonitor.Where("server = ?", m.Server).First(&m)
	if r.Error != nil && !gorm.IsRecordNotFoundError(r.Error) {
		return r.Error
	}

	return nil
}

//Post into DB
func (m *Monitor) Post(DBMonitor *gorm.DB) error {
	r := DBMonitor.Create(&m)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

//Put into DB
func (m *Monitor) Put(DBMonitor *gorm.DB) error {
	r := DBMonitor.Model(&m).Where("id = ?", m.ID).Updates(map[string]interface{}{
		"ip":                 m.IP,
		"server":             m.Server,
		"disk":               m.Disk,
		"disk_used":          m.DiskUsed,
		"ram":                m.RAM,
		"ram_used":           m.RAMUsed,
		"swap":               m.Swap,
		"swap_used":          m.SwapUsed,
		"is_backup_server":   m.IsBackupServer,
		"is_database_server": m.IsDatabaseServer,
	})
	if r.Error != nil {
		return r.Error
	}

	return nil
}
