package model

import (
	//"github.com/labstack/gommon/log"
	"github.com/lpisces/bootstrap/cmd/serve"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func Migrate() (err error) {
	config := serve.Conf
	db, err := gorm.Open(config.DB.Driver, config.DB.DataSource)
	defer db.Close()

	if err != nil {
		return err
	}

	db.AutoMigrate(&User{})
	return
}

func GetDB() (db *gorm.DB, err error) {
	config := serve.Conf
	return gorm.Open(config.DB.Driver, config.DB.DataSource)
}
