package m

import (
	valid "github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type User struct {
	gorm.Model
	Name            string `gorm:"size:255;not null;unique" valid:"length(6|15),optional" form:"name"`
	Email           string `gorm:"size:255;not null;unique" valid:"required,email" form:"email"`
	Avatar          string `gorm:"size:1024;not null;" valid:"-" form:"avatar"`
	PasswrodDigest  string `grom:"size:1024;not null;" valid:"required,alphanum"`
	Password        string `gorm:"-" valid:"-" form:"password"`
	PasswordConfirm string `gorm:"-" valid:"-" form:"password_confirm"`
}

// Validate
func (u *User) Validate() (err error) {
	ok, err := valid.ValidateStruct(u)
	if !ok {
		return
	}
	return
}
