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
	db, err := gorm.Open("mysql", "root:"+settings.DBRP+"@tcp("+host+":3306)/")
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
func CreateDatabase(host string, databaseName string) error {
	var err error
	var db *gorm.DB
	db, err = gorm.Open("mysql", "root:"+settings.DBRP+"@tcp("+host+":3306)/")
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

//DropDatabase drop database
func DropDatabase(host string, databaseName string, username string) error {
	db, err := gorm.Open("mysql", "root:"+settings.DBRP+"@tcp("+host+":3306)/")
	if err != nil {
		return err
	}
	defer db.Close()

	sql := "DROP DATABASE IF EXISTS `" + databaseName + "`"
	r := db.Exec(sql)
	if r.Error != nil {
		return r.Error
	}

	sql = "DROP USER IF EXISTS '" + username + "'@`%`"
	r = db.Exec(sql)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

//CreateUserOnHost создаёт пользователя на указанном хосте
func CreateUserOnHost(host string, databaseName string, username string, password string) error {
	var err error
	var db *gorm.DB
	db, err = gorm.Open("mysql", "rtu:"+settings.DBRTUP+"@tcp("+settings.DetermineHost(host)+":3306)/")
	if err != nil {
		return err
	}
	defer db.Close()

	//host 192.168.%
	err = db.Exec("CREATE USER IF NOT EXISTS '" + username + "'@'192.168.%' IDENTIFIED BY '" + password + "'").Error
	if err != nil {
		return err
	}

	err = db.Exec("GRANT EXECUTE, CREATE, ALTER, DROP, INDEX, LOCK TABLES, SELECT, INSERT, UPDATE, DELETE ON `" + databaseName + "`.* TO '" + username + "'@'192.168.%'").Error
	if err != nil {
		return err
	}

	//host 127.0.0.1
	err = db.Exec("CREATE USER IF NOT EXISTS '" + username + "'@'127.0.0.1' IDENTIFIED BY '" + password + "'").Error
	if err != nil {
		return err
	}

	err = db.Exec("GRANT EXECUTE, CREATE, ALTER, DROP, INDEX, LOCK TABLES, SELECT, INSERT, UPDATE, DELETE ON `" + databaseName + "`.* TO '" + username + "'@'127.0.0.1'").Error
	if err != nil {
		return err
	}

	return nil
}

//CreateDatabase if not exist
func CreateDatabaseOnHost(host string, databaseName string) error {
	var err error
	var db *gorm.DB
	db, err = gorm.Open("mysql", "rtu:"+settings.DBRTUP+"@tcp("+settings.DetermineHost(host)+":3306)/")
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

//DropDatabase drop database
func DropDatabaseOnHost(host string, databaseName string, username string) error {
	db, err := gorm.Open("mysql", "rtu:"+settings.DBRTUP+"@tcp("+settings.DetermineHost(host)+":3306)/")
	if err != nil {
		return err
	}
	defer db.Close()

	sql := "DROP DATABASE IF EXISTS `" + databaseName + "`"
	r := db.Exec(sql)
	if r.Error != nil {
		return r.Error
	}

	sql = "DROP USER IF EXISTS '" + username + "'@`127.0.0.1`"
	r = db.Exec(sql)
	if r.Error != nil {
		return r.Error
	}

	sql = "DROP USER IF EXISTS '" + username + "'@`192.168.%`"
	r = db.Exec(sql)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func tryOpenDatabase(host string, databaseName string, username string, password string, location string) (*gorm.DB, error) {
	if location == "" {
		location = "UTC"
	} else {
		location = strings.Replace(location, "/", "%2F", -1)
	}
	return gorm.Open("mysql", ""+username+":"+password+"@tcp("+host+":3306)/"+databaseName+"?charset=utf8&parseTime=True&loc="+location)
}

//OpenDatabase opens database
func OpenDatabase(databaseName string, username string, password string, location string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	var host = "localhost"

	db, err = tryOpenDatabase(host, databaseName, username, password, location)
	if err != nil {
		//if err on user than try to create
		if strings.Contains(err.Error(), "1045") {
			err = CreateDatabase(host, databaseName)
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

		//if err unknown database than create Main DB
		if strings.Contains(err.Error(), "1049") {
			err = CreateDatabase(host, databaseName)
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
