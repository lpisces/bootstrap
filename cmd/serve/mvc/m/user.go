package m

import (
	"fmt"
	//"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/lpisces/bootstrap/cmd/serve"
	valid "gopkg.in/asaskevich/govalidator.v4"
	"time"
)

type User struct {
	//gorm.Model
	ID              uint `gorm:"primary_key"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Name            string `gorm:"size:255;not null;unique" valid:"length(6|15),optional" form:"name"`
	Email           string `gorm:"size:255;not null;unique" valid:"required~请输入Email,email~Email格式不正确" form:"email"`
	Avatar          string `gorm:"size:1024;not null;" valid:"optional,alphanum" form:"avatar"`
	PasswrodDigest  string `grom:"size:1024;not null;" valid:"-"`
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
	defer db.Close()

	uu := &User{}
	if !(db.Where("email = ?", u.Email).First(uu).RecordNotFound()) {
		return fmt.Errorf("%s 已经存在", u.Email)
	}

	hash, err := Crypt(u.Password) //bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	if err != nil {
		return
	}

	u.PasswrodDigest = string(hash)
	db.Create(u)
	return
}

// Exist
func (u *User) Exist() (exist bool, err error) {

	db, err := GetDB()
	if err != nil {
		return
	}
	defer db.Close()

	uu := &User{}
	if db.Where("email = ?", u.Email).First(uu).RecordNotFound() {
		exist = false
		return
	}
	exist = true
	return
}

// Auth
func (u *User) Auth() (bool, error) {
	exist, err := u.Exist()
	if err != nil {
		return false, err
	}
	if !exist {
		return false, nil
	}

	hash, err := Crypt(u.Password) //bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	if err != nil {
		return false, err
	}

	db, err := GetDB()
	if err != nil {
		return false, err
	}
	defer db.Close()

	if db.Where("password_digest = ?", hash).First(&User{}).RecordNotFound() {
		return false, nil
	}

	return true, nil
}

// SignIn
func (u *User) SignIn(c echo.Context) (err error) {
	exist, err := u.Exist()
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("user not exists")
	}

	siteConfig := serve.Conf.Site
	sess, err := session.Get(siteConfig.SessionName, c)
	if err != nil {
		return err
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	sess.Values["uid"] = u.ID
	sess.Save(c.Request(), c.Response())
	return
}

// SignOut
func (u *User) SignOut(c echo.Context) (err error) {
	siteConfig := serve.Conf.Site
	sess, err := session.Get(siteConfig.SessionName, c)
	if err != nil {
		return err
	}

	delete(sess.Values, "uid")
	sess.Save(c.Request(), c.Response())
	return
}
