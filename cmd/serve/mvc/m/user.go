package m

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	valid "gopkg.in/asaskevich/govalidator.v4"
)

type User struct {
	gorm.Model
	Name            string `gorm:"size:255;not null;unique" valid:"length(6|15),optional" form:"name"`
	Email           string `gorm:"size:255;not null;unique" valid:"required,email" form:"email"`
	Avatar          string `gorm:"size:1024;not null;" valid:"optional,alphanum" form:"avatar"`
	PasswrodDigest  string `grom:"size:1024;not null;" valid:"optional,alphanum"`
	Password        string `gorm:"-" valid:"required,length(6|15)" form:"password"`
	PasswordConfirm string `gorm:"-" valid:"required,length(6|15)" form:"password_confirm"`
}

// Validate
func (u *User) Validate() (ok bool, errs map[string]string) {
	ok, err := valid.ValidateStruct(u)
	if !ok {
		errs = valid.ErrorsByField(err)
	}
	return
}
