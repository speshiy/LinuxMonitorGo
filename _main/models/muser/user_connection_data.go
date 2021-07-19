package muser

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//UserConnectionData structure
type UserConnectionData struct {
	gorm.Model
	DBHost     string `gorm:"column:db_host;"`
	DBName     string `gorm:"column:db_name;"`
	DBUser     string `gorm:"column:db_user;"`
	DBPassword string `gorm:"column:db_password;"`
}

//TableName return new table name for User struct
func (UserConnectionData) TableName() string {
	return "s_users_connection_data_t"
}

//GetConnectionsData return connections data
func GetConnectionsData(c *gin.Context, cd *[]UserConnectionData, DBMain *gorm.DB) error {
	if DBMain == nil {
		DBMain = c.MustGet("DBMain").(*gorm.DB)
	}

	r := DBMain.Find(&cd)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

//GetConnectionsDataByDBName return connections data by dbName as array
func GetConnectionsDataByDBName(c *gin.Context, dbName string, cd *[]UserConnectionData, DBMain *gorm.DB) error {
	if dbName == "" {
		err := GetConnectionsData(c, cd, DBMain)
		if err != nil {
			return err
		}
		return nil
	}

	if DBMain == nil {
		DBMain = c.MustGet("DBMain").(*gorm.DB)
	}

	r := DBMain.Where("db_name = ?", dbName).Find(&cd)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

//GetConnectionDataByDBName return connection data
func GetConnectionDataByDBName(c *gin.Context, cd *UserConnectionData, DBMain *gorm.DB) error {
	if DBMain == nil {
		DBMain = c.MustGet("DBMain").(*gorm.DB)
	}

	r := DBMain.Where("db_name = ?", cd.DBName).Find(&cd)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

//GetConnectionsDataByIDsRange return connections data by dbName as array
func GetConnectionsDataByIDsRange(c *gin.Context, idStart uint, idEnd uint, cd *[]UserConnectionData, DBMain *gorm.DB) error {
	if DBMain == nil {
		DBMain = c.MustGet("DBMain").(*gorm.DB)
	}

	r := DBMain.Where("id >= ? AND id <= ?", idStart, idEnd).Find(&cd)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

//PutDBHost обновляет DBHost
func (cd *UserConnectionData) PutDBHost(c *gin.Context, DBMain *gorm.DB) error {
	if DBMain == nil {
		DBMain = c.MustGet("DBMain").(*gorm.DB)
	}

	r := DBMain.Model(&cd).Where("id = ?", cd.ID).
		Updates(map[string]interface{}{
			"db_host": cd.DBHost,
		})
	if r.Error != nil {
		return r.Error
	}

	return nil
}
