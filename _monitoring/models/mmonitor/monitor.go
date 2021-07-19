package mmonitor

import (
	"github.com/jinzhu/gorm"
)

//Monitor stucture
type Monitor struct {
	gorm.Model
	IP                 string `json:"IP"`
	Server             string `json:"Server"`
	RAM                uint64 `json:"RAM"`
	RAMUsed            uint64 `json:"RAMUsed"`
	Swap               uint64 `json:"Swap"`
	SwapUsed           uint64 `json:"SwapUsed"`
	Sort               uint   `json:"Sort"`
	IsBackupServer     bool   `json:"IsBackupServer"`
	IsBackupInProgress bool   `json:"IsBackupInProgress"`
	BackupingDBName    string `json:"BackupingDBName"`
}
