package utils

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/gommon/log"
)

var Conn *gorm.DB

func GetConn(conf Config) (db *gorm.DB, err error){
	log.Print(conf.Database)
	db, err = gorm.Open(conf.Database.Dialect, conf.Database.Args)
	return
}
