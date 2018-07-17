package c

import (
	"encoding/gob"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/lpisces/bootstrap/cmd/serve"
	"github.com/lpisces/bootstrap/cmd/serve/mvc/m"
	"net/http"
)

func init() {
	gob.Register(RegFlash{})
	gob.Register(LoginFlash{})
	gob.Register(ForgetPasswordFlash{})
}

// GetSession
func GetSession(c echo.Context) (sess *sessions.Session, err error) {
	siteConfig := serve.Conf.Site
	sess, err = session.Get(siteConfig.SessionName, c)
	if err != nil {
		return
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	return
}

// IfLoginedRedirectTo
func IfLoginedRedirectTo(c echo.Context, path string) error {
	// logined user redirect
	if IsLogined(c) {
		return c.Redirect(http.StatusFound, path)
	}
	return nil
}

// IfNotLoginedRedirectTo
func IfNotLoginedRedirectTo(c echo.Context, path string) error {
	if !IsLogined(c) {
		return c.Redirect(http.StatusFound, path)
	}
	return nil
}

// IsLogined
func IsLogined(c echo.Context) (ok bool) {
	_, err := CurrentUser(c)
	if err != nil {
		return
	}
	return true
}

// CurrentUser
func CurrentUser(c echo.Context) (user *m.User, err error) {

	sess, err := GetSession(c)
	if err != nil {
		return user, err
	}

	if sess == nil {
		err = fmt.Errorf("no session")
		return
	}

	uid, ok := sess.Values["uid"]
	if !ok {
		err = fmt.Errorf("uid not found")
		return
	}

	db, err := m.GetDB()
	if err != nil {
		return
	}
	defer db.Close()

	user = &m.User{}
	if db.First(user, uid).RecordNotFound() {
		err = fmt.Errorf("user not found")
		return
	}
	return
}
