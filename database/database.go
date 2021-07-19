package database

import (
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/speshiy/LinuxMonitorGo/settings"
)

//CreateUser create user
func CreateUser(host string, databaseName string, username string, password string) error {
	db, err := gorm.Open("mysql", settings.DatabaseMonitorUser+":"+settings.DatabaseMonitorPassword+"@tcp("+host+":3306)/")
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Exec("CREATE USER IF NOT EXISTS '" + username + "'@'%' IDENTIFIED BY '" + password + "'").Error
	if err != nil {
		return err
	}

	err = db.Exec("GRANT ALL PRIVILEGES ON " + databaseName + ".* TO '" + username + "'@'%' WITH GRANT OPTION").Error
	if err != nil {
		return err
	}

	return nil
}

//CreateDatabase if not exist
func CreateDatabase(host string, databaseName string, username string, password string) error {
	db, err := gorm.Open("mysql", username+":"+password+"@tcp("+host+":3306)/")
	if err != nil {
		return err
	}

	err = db.Exec("CREATE DATABASE IF NOT EXISTS " + databaseName + " CHARACTER SET utf8 COLLATE utf8_general_ci").Error
	if err != nil {
		return err
	}

	db.Close()

	return err
}

func tryOpenDatabase(host string, databaseName string, username string, password string, location string) (*gorm.DB, error) {
	if location == "" {
		location = "UTC"
	} else {
		location = strings.Replace(location, "/", "%2F", -1)
	}
	return gorm.Open("mysql", ""+username+":"+password+"@tcp("+host+":3306)/"+databaseName+"?charset=utf8&parseTime=True&loc="+location)
}

//OpenDatabase open database
func OpenDatabase(host string, databaseName string, username string, password string, location string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	db, err = tryOpenDatabase(host, databaseName, username, password, location)
	if err != nil {
		if strings.Contains(err.Error(), "1045") {
			err = CreateDatabase(host, databaseName, username, password)
			if err != nil {
				return nil, err
			}

			err = CreateUser(host, databaseName, username, password)
			if err != nil {
				return nil, err
			}

			db, err = tryOpenDatabase(host, databaseName, username, password, location)
			if err != nil {
				return nil, err
			}

			return db, nil
		}

		if strings.Contains(err.Error(), "1049") {
			err = CreateDatabase(host, databaseName, username, password)
			if err != nil {
				return nil, err
			}

			db, err = tryOpenDatabase(host, databaseName, username, password, location)
			if err != nil {
				return nil, err
			}

			return db, nil
		}

		log.Print("Connection to DB in service FAILED.", err.Error())
		return nil, err
	}

	return db, nil
}

//ReturnConnectionToDBMonitoring open connection to DB
func ReturnConnectionToDBMonitoring() (*gorm.DB, error) {
	var err error
	var DB *gorm.DB

	DB, err = OpenDatabase(settings.DatabaseMonitorHost, settings.DatabaseMonitorName, settings.DatabaseMonitorUser, settings.DatabaseMonitorPassword, "UTC")
	if err != nil {
		return nil, err
	}

	return DB, nil
}
