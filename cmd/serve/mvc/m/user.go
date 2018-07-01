package m

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"golang.org/x/crypto/bcrypt"
	valid "gopkg.in/asaskevich/govalidator.v4"
)

type User struct {
	gorm.Model
	Name            string `gorm:"size:255;not null;unique" valid:"length(6|15),optional" form:"name"`
	Email           string `gorm:"size:255;not null;unique" valid:"required~请输入Email,email~Email格式不正确" form:"email"`
	Avatar          string `gorm:"size:1024;not null;" valid:"optional,alphanum" form:"avatar"`
	PasswrodDigest  string `grom:"size:1024;not null;" valid:"optional,alphanum"`
	Password        string `gorm:"-" valid:"required~请输入密码,length(6|15)~密码长度必须在6到15位之间" form:"password"`
	PasswordConfirm string `gorm:"-" valid:"required~请确认密码" form:"password_confirm"`
}

// Validate
func (u *User) Validate() (ok bool, errs map[string]string) {
	ok, err := valid.ValidateStruct(u)
	errs = make(map[string]string)
	if !ok {
		errs = valid.ErrorsByField(err)
	}
	if u.Password != u.PasswordConfirm {
		errs["PasswordConfirm"] = "两次输入密码不一致"
		ok = false
	}
	return
}

// Create
func (u *User) Create() (err error) {
	db, err := GetDB()
	if err != nil {
		return
	}

	uu := &User{}
	if !db.Where("email = ?", u.Email).First(uu).RecordNotFound() {
		return fmt.Errorf("%s 已经存在", u.Email)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	if err != nil {
		return
	}

	u.PasswrodDigest = string(hash)
	db.Create(u)
	return
}
