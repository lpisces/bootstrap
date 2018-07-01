package m

import (
	//"github.com/labstack/gommon/log"
	"github.com/lpisces/bootstrap/cmd/serve"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	DB *gorm.DB
)

func Migrate() (err error) {
	//config := serve.Conf
	db, err := GetDB() //gorm.Open(config.DB.Driver, config.DB.DataSource)
	defer db.Close()

	if err != nil {
		return err
	}

	db.AutoMigrate(&User{})
	return
}

func GetDB() (*gorm.DB, error) {
	config := serve.Conf
	if DB != nil {
		return DB, nil
	}
	DB, err := gorm.Open(config.DB.Driver, config.DB.DataSource)
	return DB, err
}
