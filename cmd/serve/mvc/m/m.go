package m

import (
	//"github.com/labstack/gommon/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/lpisces/bootstrap/cmd/serve"
)

var (
	DB *gorm.DB
)

func Migrate() (err error) {
	db, err := GetDB() //gorm.Open(config.DB.Driver, config.DB.DataSource)
	defer db.Close()

	if err != nil {
		return err
	}

	return
}

func GetDB() (*gorm.DB, error) {
	config := serve.Conf
	if DB != nil {
		return DB, nil
	}
	DB, err := gorm.Open(config.DB.Driver, config.DB.DataSource)
	if serve.Debug {
		DB.LogMode(true)
	}
	return DB, err
}
