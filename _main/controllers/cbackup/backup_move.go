package cbackup

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/speshiy/LinuxMonitorGo/_main/models/mbackup"
	"github.com/speshiy/LinuxMonitorGo/_main/models/muser"
	"github.com/speshiy/LinuxMonitorGo/database"
	"github.com/speshiy/LinuxMonitorGo/settings"
)

//DatabaseMoveFromServerToServer перемещает базу с одного сервера на другой
func DatabaseMoveFromServerToServer(c *gin.Context) {
	var err error

	var userConnections []muser.UserConnectionData
	var incomeData struct {
		DestHost        string `json:"DestHost"`
		DatabaseIDStart uint   `json:"DatabaseIDStart"`
		DatabaseIDEnd   uint   `json:"DatabaseIDEnd"`
	}

	if err = c.ShouldBindJSON(&incomeData); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	err = muser.GetConnectionsDataByIDsRange(c, incomeData.DatabaseIDStart, incomeData.DatabaseIDEnd, &userConnections, nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	for _, uc := range userConnections {
		// Если база уже находится на целевом хосте, то пропускаем её перенос
		if uc.DBHost == settings.DetermineHost(incomeData.DestHost) {
			continue
		}

		DBHostOld := uc.DBHost
		DBHostNew := incomeData.DestHost

		// Бэкапируем базу
		err = DoBackup(c, uc.DBName, nil, nil, false, false)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
			return
		}

		// Создаём базу на целевом хосте
		err = database.CreateDatabaseOnHost(DBHostNew, uc.DBName)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
			return
		}

		// Создаём пользователя на новом хосте
		err = database.CreateUserOnHost(DBHostNew, uc.DBName, uc.DBUser, uc.DBPassword)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
			return
		}

		var backup mbackup.Backup
		backup.DBName = uc.DBName
		err = backup.GetLastByDBName(c, nil)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
			return
		}

		// Переносим базу на новый хост
		uc.DBHost = DBHostNew
		err = doBackupRestore(&backup, &uc)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
			return
		}

		// Удаляем базу пользователя на старом хосте
		uc.DBHost = DBHostOld
		err = database.DropDatabaseOnHost(uc.DBHost, uc.DBName, uc.DBUser)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
			return
		}

		// Обновляем хоста в таблице
		uc.DBHost = settings.DetermineHost(DBHostNew)
		err = uc.PutDBHost(c, nil)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
			return
		}

	}

	c.JSON(http.StatusOK, gin.H{"status": "true", "message": "Базы данных были успешно перенесены"})
}
