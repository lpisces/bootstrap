package c

import (
	//"github.com/gorilla/sessions"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/lpisces/bootstrap/cmd/serve"
	"github.com/lpisces/bootstrap/cmd/serve/mvc/m"
)

// IsLogin
func IsLogin(c echo.Context) (ok bool, err error) {
	siteConfig := serve.Conf.Site
	sess, err := session.Get(siteConfig.SessionName, c)
	if err != nil {
		return false, err
	}

	_, ok = sess.Values["uid"]
	return
}

// CurrentUser
func CurrentUser(c echo.Context) (user *m.User, err error) {

	user = &m.User{}

	siteConfig := serve.Conf.Site
	sess, err := session.Get(siteConfig.SessionName, c)
	if err != nil {
		return user, err
	}

	uid, ok := sess.Values["uid"]
	if !ok {
		err = fmt.Errorf("not signed in")
		return
	}

	db, err := m.GetDB()
	if err != nil {
		return
	}
	defer db.Close()

	db.First(user, uid)
	return

}
