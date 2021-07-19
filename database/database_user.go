package database

import (
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//DBUserConnect connect
type DBUserConnect struct {
	Host     string
	DBName   string
	User     string
	Password string
	Connect  *gorm.DB
}

//Open Main Database
func (db *DBUserConnect) Open(host string, dbname string, user string, password string, location string) error {
	var err error

	if location == "" {
		location = "UTC"
	} else {
		location = strings.Replace(location, "/", "%2F", -1)
	}

	db.Host = host
	db.DBName = dbname
	db.User = user
	db.Password = password
	db.Connect, err = gorm.Open("mysql", db.User+":"+db.Password+"@tcp("+db.Host+":3306)/"+db.DBName+"?charset=utf8&parseTime=True&loc="+location)
	if err != nil {
		log.Println("Connection to user DB FAILED -", dbname, ".", err.Error())
		return err
	}

	return nil
}

//Close opens Main Database
func (db *DBUserConnect) Close() error {
	db.Connect.Close()
	return nil
}
