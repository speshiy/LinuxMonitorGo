package cbackup

import (
	"compress/flate"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mholt/archiver"
	"github.com/speshiy/LinuxMonitorGo/_main/models/mbackup"
	"github.com/speshiy/LinuxMonitorGo/_main/models/muser"
	"github.com/speshiy/LinuxMonitorGo/database"
	"github.com/speshiy/LinuxMonitorGo/settings"
)

//Table structure
type Table struct {
	Name string `gorm:"-" json:"Name"`
}

var dumpsPath = "/tuvis/mysql/dumps"

//BackupByDBName start backuping DB
func BackupByDBName(c *gin.Context) {
	var err error
	var DBMain database.DBUserConnect

	dbName := c.Query("db")

	err = DBMain.Open(settings.DBHostDump, "main", "rtu", settings.DBRTUP, "UTC")
	if err != nil {
		return
	}
	defer DBMain.Close()

	err = DoBackup(c, dbName, DBMain.Connect, nil, false, false)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "true", "message": "Дампирование успешно завершено!"})
}

//DoBackup do backup
func DoBackup(c *gin.Context, dbName string, DBMain *gorm.DB, DB *gorm.DB, checkBackupInterval bool, isDeletedDatabase bool) error {
	var err error
	var backupPath string
	var backupName string
	var userConnections []muser.UserConnectionData
	var backup mbackup.Backup

	err = os.MkdirAll(settings.ResourcesPath+dumpsPath, 0777)
	if err != nil {
		return err
	}

	if dbName != "system" {
		err = muser.GetConnectionsDataByDBName(c, dbName, &userConnections, DBMain)
		if err != nil {
			return err
		}
	} else if dbName == "system" {
		dbName = ""
	}

	if dbName == "" {
		//Append main DB in LIST
		userConnections = append(userConnections, muser.UserConnectionData{
			DBHost:     settings.DBHostDump,
			DBName:     "main",
			DBPassword: "",
			DBUser:     "main",
		})

		//Append client DB in LIST
		userConnections = append(userConnections, muser.UserConnectionData{
			DBHost:     settings.DBHostDump,
			DBName:     "client",
			DBPassword: "",
			DBUser:     "client",
		})
	} else if dbName == "client" {
		userConnections = append(userConnections, muser.UserConnectionData{
			DBHost:     settings.DBHostDump,
			DBName:     "client",
			DBPassword: "",
			DBUser:     "client",
		})
	} else if dbName == "main" {
		userConnections = append(userConnections, muser.UserConnectionData{
			DBHost:     settings.DBHostDump,
			DBName:     "main",
			DBPassword: "",
			DBUser:     "main",
		})
	}

	zip := archiver.Zip{
		CompressionLevel:       flate.BestCompression,
		MkdirAll:               true,
		SelectiveCompression:   true,
		ContinueOnError:        false,
		OverwriteExisting:      true,
		ImplicitTopLevelFolder: false,
	}

	for _, uc := range userConnections {
		if uc.DBHost == "127.0.0.1" || uc.DBHost == "localhost" {
			uc.DBHost = settings.DBHostDump
		}

		if checkBackupInterval {
			isNeedBackup, err := mbackup.IsNeedBackup(c, uc.DBName, DB)
			if err != nil {
				return err
			}

			if !isNeedBackup {
				continue
			}
		}

		settings.BackupingDBName = uc.DBName

		backupPath, err = backupAndZip(&zip, uc)
		if err != nil {
			return err
		}

		backup = mbackup.Backup{}
		backup.DBName = uc.DBName
		backup.UserName = uc.DBUser
		backup.Host = settings.DetermineHost(uc.DBHost)
		backup.BackupPath = backupPath
		backupName = backupPath[strings.LastIndex(backupPath, "/")+1:]
		backup.BackupName = backupName
		backup.IsDeletedDatabase = isDeletedDatabase

		fi, err := os.Stat(backup.BackupPath)
		if err != nil {
			log.Println(err.Error())
		}
		backup.FileSize = float32(fi.Size()) / 1024

		err = backup.Post(c, DB)
		if err != nil {
			return err
		}

		log.Println(backup.DBName, "was dumped, wait next dump in 1 seconds...")

		// Ждём 10 секунд
		time.Sleep(time.Second * 1)
	}

	settings.BackupingDBName = ""

	return nil
}

func backupAndZip(zip *archiver.Zip, uc muser.UserConnectionData) (string, error) {
	var zipFileName string
	var DBUser database.DBUserConnect
	var err error

	username := "rtu"
	password := settings.DBRTUP
	hostname := settings.DetermineHost(uc.DBHost)
	port := "3306"
	dbname := uc.DBName

	dumpDir := settings.ResourcesPath + dumpsPath
	dirName := dumpDir + "/" + time.Now().Format("2006_01_02_15_04_05_") + dbname
	os.MkdirAll(dirName, os.ModeDir)

	err = DBUser.Open(hostname, dbname, username, password, "UTC")
	if err != nil {
		return "", err
	}
	defer DBUser.Close()

	// Получаем список таблиц по БД
	tables := []Table{}
	sql := `SELECT 
						table_name as name
					FROM 
						information_schema.tables 
					WHERE 
						table_schema NOT IN ('information_schema','mysql') AND table_schema=?
					ORDER BY
						table_name`

	r := DBUser.Connect.Raw(sql, dbname).Scan(&tables)
	if r.Error != nil {
		return "", err
	}

	for _, table := range tables {

		cmd := exec.Command("mysqldump", "-P"+port, "-h"+hostname, "-u"+username, "-p"+password, "--default-character-set=utf8", "--skip-add-locks", "--quick", "-e", dbname, table.Name)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}

		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}

		bytes, err := ioutil.ReadAll(stdout)
		if err != nil {
			log.Fatal(err)
		}

		cmd.Wait()

		sqlFileName := dirName + "/" + table.Name + ".sql"

		err = ioutil.WriteFile(sqlFileName, bytes, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	// ZIP sql file
	zipFileName = dirName + ".zip"
	err = zip.Archive([]string{dirName}, zipFileName)
	if err != nil {
		fmt.Println("Error zipping:", err)
		return "", err
	}

	//Remove unzipped dump
	os.RemoveAll(dirName)

	return zipFileName, nil
}
