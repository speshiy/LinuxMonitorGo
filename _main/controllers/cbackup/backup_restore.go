package cbackup

import (
	"compress/flate"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/mholt/archiver"
	"github.com/speshiy/LinuxMonitorGo/_main/models/mbackup"
	"github.com/speshiy/LinuxMonitorGo/_main/models/muser"
	"github.com/speshiy/LinuxMonitorGo/settings"
)

//BackupRestore restore backup by ID
func BackupRestore(c *gin.Context) {
	var err error
	var backup mbackup.Backup
	var userConnection muser.UserConnectionData
	var incomeData struct {
		BackupID uint `json:"BackupID"`
	}

	if err = c.Bind(&incomeData); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	if incomeData.BackupID == 0 {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": "ID дампа неверный"})
		return
	}

	backup.ID = incomeData.BackupID
	err = backup.GetByID(c, nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	if backup.DBName == "main" {
		userConnection.DBName = "main"
		userConnection.DBUser = "rtu"
		userConnection.DBPassword = settings.DBRTUP
		userConnection.DBHost = settings.DBHostDump
	} else if backup.DBName == "client" {
		userConnection.DBName = "client"
		userConnection.DBUser = "rtu"
		userConnection.DBPassword = settings.DBRTUP
		userConnection.DBHost = settings.DBHostDump
	} else {
		userConnection.DBName = backup.DBName
		err = muser.GetConnectionDataByDBName(c, &userConnection, nil)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
			return
		}
	}

	err = DoBackup(c, backup.DBName, nil, nil, false, false)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	err = doBackupRestore(&backup, &userConnection)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	err = backup.Put(c, nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "true", "message": "Дамп был успешно восстановлен"})
}

func doBackupRestore(backup *mbackup.Backup, userConnection *muser.UserConnectionData) error {
	var err error

	if userConnection.DBHost == "127.0.0.1" || userConnection.DBHost == "localhost" {
		userConnection.DBHost = settings.DBHostDump
	}

	zip := archiver.Zip{
		CompressionLevel:       flate.BestCompression,
		MkdirAll:               true,
		SelectiveCompression:   true,
		ContinueOnError:        false,
		OverwriteExisting:      true,
		ImplicitTopLevelFolder: true,
	}

	//Unarchive zip
	dumpDir := settings.ResourcesPath + dumpsPath
	err = zip.Unarchive(backup.BackupPath, dumpDir)
	if err != nil {
		return err
	}

	unzipedDir := string(backup.BackupPath[:len(backup.BackupPath)-4])
	files, err := ioutil.ReadDir(unzipedDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		username := "rtu"
		password := settings.DBRTUP
		hostname := settings.DetermineHost(userConnection.DBHost)
		port := "3306"
		dbname := userConnection.DBName

		filename := unzipedDir + "/" + f.Name()

		cmd := exec.Command("mysql", "-P"+port, "-h"+hostname, "-u"+username, "-p"+password, "-D"+dbname, "-e", "source "+filename)
		if err := cmd.Run(); err != nil {
			log.Println(err)
		}
	}

	//Remove unzipped dump
	os.RemoveAll(unzipedDir)

	return nil
}
