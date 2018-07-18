package m

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	//"github.com/labstack/gommon/log"
	"bytes"
	"github.com/lpisces/bootstrap/cmd/serve"
	valid "gopkg.in/asaskevich/govalidator.v4"
	"html/template"
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
	PasswordDigest  string `grom:"size:1024;not null;column:password_digest;" valid:"-"`
	Password        string `gorm:"-" valid:"required~请输入密码,length(6|15)~密码长度必须在6到15位之间" form:"password"`
	PasswordConfirm string `gorm:"-" valid:"required~请确认密码" form:"password_confirm"`
	Admin           bool   `gorm:"not null;default:false;" valid:"optional"`
	Activated       bool   `gorm:"not null;default:false;" valid:"optional"`
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

	hash, err := HashPassword(u.Password) //bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	if err != nil {
		return
	}

	u.Name = u.Email
	u.PasswordDigest = hash
	db.Create(u)

	NewToken(TokenTypeActivate, u)
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

// Load
func (u *User) Load() (err error) {

	db, err := GetDB()
	if err != nil {
		return
	}
	defer db.Close()

	if db.Where("email = ?", u.Email).First(u).RecordNotFound() {
		return fmt.Errorf("not exists")
	}
	return
}

// Auth
func (u *User) Auth() (bool, error) {

	db, err := GetDB()
	if err != nil {
		return false, err
	}
	defer db.Close()

	if db.Where(&User{Email: u.Email}).First(u).RecordNotFound() {
		return false, nil
	}

	if !CheckPasswordHash(u.Password, u.PasswordDigest) {
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

// LastActivateToken
func (u *User) LastToken(tokenType TokenType) (t *Token, err error) {

	db, err := GetDB()
	if err != nil {
		return t, err
	}
	defer db.Close()

	t = &Token{}
	if db.Where("user_id = ?", u.ID).Where("type = ?", tokenType).Last(t).RecordNotFound() {
		err = fmt.Errorf("no activate token found")
	}
	return
}

// IsAdmin
func (u *User) IsAdmin() (ok bool) {
	return u.Admin || u.ID == 1
}

// SendActivateMail
func (u *User) SendActivateMail() (err error) {

	tmpl := template.Must(template.ParseGlob("cmd/serve/mvc/v/mail/activate_mail.html"))
	token, err := u.LastToken(TokenTypeActivate)
	if err != nil {
		return err
	}

	data := struct {
		User     *User
		Token    *Token
		Title    string
		SiteName string
		BaseURL  string
	}{
		User:     u,
		Token:    token,
		Title:    serve.Conf.Site.Name + "-" + "注册激活邮件",
		SiteName: serve.Conf.Site.Name,
		BaseURL:  serve.Conf.Site.BaseURL,
	}

	var mailHTML bytes.Buffer
	if err := tmpl.Execute(&mailHTML, data); err != nil {
		return err
	}

	config := serve.Conf.SMTP
	mail := &MailData{
		To:      []string{u.Email},
		From:    fmt.Sprintf("%s <%s>", config.FromName, config.FromAddr),
		Subject: fmt.Sprintf("[%s]注册激活邮件", serve.Conf.Site.Name),
		//Text:    []byte("activate mail"),
		HTML: mailHTML.Bytes(),
	}
	return SendMail(mail)
}

// SendForgetPasswordMail
func (u *User) SendForgetPasswordMail() (err error) {
	tmpl := template.Must(template.ParseGlob("cmd/serve/mvc/v/mail/forget_password_mail.html"))
	token, err := u.LastToken(TokenTypeResetPassword)
	if err != nil {
		return err
	}

	data := struct {
		User     *User
		Token    *Token
		Title    string
		SiteName string
		BaseURL  string
	}{
		User:     u,
		Token:    token,
		Title:    serve.Conf.Site.Name + "-" + "重置密码邮件",
		SiteName: serve.Conf.Site.Name,
		BaseURL:  serve.Conf.Site.BaseURL,
	}

	var mailHTML bytes.Buffer
	if err := tmpl.Execute(&mailHTML, data); err != nil {
		return err
	}

	config := serve.Conf.SMTP
	mail := &MailData{
		To:      []string{u.Email},
		From:    fmt.Sprintf("%s <%s>", config.FromName, config.FromAddr),
		Subject: fmt.Sprintf("[%s]重置密码邮件", serve.Conf.Site.Name),
		//Text:    []byte("activate mail"),
		HTML: mailHTML.Bytes(),
	}
	return SendMail(mail)
}

// Activate
func (u *User) Activate() (err error) {
	db, err := GetDB()
	if err != nil {
		return err
	}
	defer db.Close()

	u.Activated = true

	db.Save(u)
	return
}

// Save
func (u *User) Save() (err error) {
	db, err := GetDB()
	if err != nil {
		return err
	}
	defer db.Close()

	hash, err := HashPassword(u.Password) //bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	if err != nil {
		return
	}

	u.PasswordDigest = hash
	db.Save(u)
	return
}
