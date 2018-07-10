package c

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/lpisces/bootstrap/cmd/serve"
	"github.com/lpisces/bootstrap/cmd/serve/mvc/m"
)

// GetSession
func GetSession(c echo.Context) (sess *sessions.Session, err error) {
	siteConfig := serve.Conf.Site
	return session.Get(siteConfig.SessionName, c)
}

// IsLogin
func IsLogin(c echo.Context) (ok bool, err error) {
	sess, err := GetSession(c)
	if err != nil {
		return false, err
	}

	_, ok = sess.Values["uid"]
	return
}

// CurrentUser
func CurrentUser(c echo.Context) (user *m.User, err error) {

	sess, err := GetSession(c)
	if err != nil {
		return user, err
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
