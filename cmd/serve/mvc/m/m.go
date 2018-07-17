package m

import (
	//"github.com/labstack/gommon/log"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/jordan-wright/email"
	"github.com/lpisces/bootstrap/cmd/serve"
	"golang.org/x/crypto/bcrypt"
	"net/smtp"
	"net/textproto"
)

var (
	DB *gorm.DB
)

type MailData struct {
	To      []string
	From    string
	Cc      []string
	Bcc     []string
	Subject string
	Text    []byte
	HTML    []byte
}

// SendMail
func SendMail(data *MailData) (err error) {
	conf := serve.Conf.SMTP
	e := &email.Email{
		To:      data.To,
		From:    data.From,
		Subject: data.Subject,
		Text:    data.Text,
		HTML:    data.HTML,
		Headers: textproto.MIMEHeader{},
	}
	return e.Send(
		fmt.Sprintf("%s:%s", conf.Hostname, conf.Port),
		smtp.PlainAuth("", conf.Username, conf.Password, conf.Hostname))
}

func Migrate() (err error) {
	//config := serve.Conf
	db, err := GetDB() //gorm.Open(config.DB.Driver, config.DB.DataSource)
	defer db.Close()

	if err != nil {
		return err
	}

	db.AutoMigrate(&User{}, &Token{})
	return
}

func GetDB() (*gorm.DB, error) {
	config := serve.Conf
	if DB != nil {
		return DB, nil
	}
	DB, err := gorm.Open(config.DB.Driver, config.DB.DataSource)
	if config.Mode != "production" {
		DB.LogMode(true)
	}
	return DB, err
}

func Crypt(str string) (hash string, err error) {
	//config := serve.Conf
	hashByte, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.MinCost)
	hash = string(hashByte)
	return
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
