package c

import (
	"github.com/labstack/echo"
	//"github.com/labstack/gommon/log"
	"github.com/lpisces/bootstrap/cmd/serve"
	"github.com/lpisces/bootstrap/cmd/serve/mvc/m"
	"net/http"
)

type (
	RegFlash struct {
		Title    string
		SiteName string
		User     *m.User
		Error    map[string]string
	}

	LoginFlash struct {
		Title    string
		SiteName string
		User     *m.User
		Error    map[string]string
	}

	ForgetPasswordFlash struct {
		Title    string
		SiteName string
		User     *m.User
		Error    map[string]string
	}
)

// GetRegister register page
func GetRegister(c echo.Context) (err error) {

	IfLoginedRedirectTo(c, "/")

	// prepare data
	data := RegFlash{
		Title:    serve.Conf.Site.Name + "-" + "注册",
		SiteName: serve.Conf.Site.Name,
		User:     &m.User{},
		Error:    make(map[string]string),
	}

	// session flashes
	sess, err := GetSession(c)
	if flashes := sess.Flashes(); len(flashes) > 0 {
		data = flashes[0].(RegFlash)
	}
	sess.Save(c.Request(), c.Response())

	return c.Render(http.StatusOK, "register", data)
}

// PostRegister handle register request
func PostRegister(c echo.Context) (err error) {

	IfLoginedRedirectTo(c, "/")

	// prepare data
	data := RegFlash{
		Title:    serve.Conf.Site.Name + "-" + "注册",
		SiteName: serve.Conf.Site.Name,
		User:     &m.User{},
		Error:    make(map[string]string),
	}

	user := new(m.User)
	if err = c.Bind(user); err != nil {
		return
	}
	data.User = user

	sess, err := GetSession(c)
	if err != nil {
		return err
	}

	// validate
	if ok, errs := user.Validate(); !ok {
		data.Error = errs
		sess.AddFlash(data)
		sess.Save(c.Request(), c.Response())
		return c.Redirect(http.StatusFound, "/register")
	}

	// create user account
	if err := user.Create(); err != nil {
		data.Error = map[string]string{"Email": err.Error()}
		sess.AddFlash(data)
		sess.Save(c.Request(), c.Response())
		return c.Redirect(http.StatusFound, "/register")
	}

	return c.Redirect(http.StatusFound, "/login")
}

// GetLogout handle logout request
func GetLogout(c echo.Context) (err error) {

	IfNotLoginedRedirectTo(c, "/")

	// get logined user
	user, err := CurrentUser(c)
	if err != nil {
		return err
	}

	// do signout action
	if err = user.SignOut(c); err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, "/login")
}

// GetLogin login page
func GetLogin(c echo.Context) (err error) {

	IfLoginedRedirectTo(c, "/")

	data := LoginFlash{
		Title:    serve.Conf.Site.Name + "-" + "登录",
		SiteName: serve.Conf.Site.Name,
		Error:    map[string]string{},
		User:     &m.User{},
	}

	// session flashes
	sess, err := GetSession(c)
	if flashes := sess.Flashes(); len(flashes) > 0 {
		data = flashes[0].(LoginFlash)
	}
	sess.Save(c.Request(), c.Response())

	return c.Render(http.StatusOK, "login", data)
}

// PostLogin handle login request
func PostLogin(c echo.Context) (err error) {

	data := LoginFlash{
		Title:    serve.Conf.Site.Name + "-" + "登录",
		SiteName: serve.Conf.Site.Name,
		Error:    map[string]string{},
		User:     &m.User{},
	}

	user := new(m.User)
	if err = c.Bind(user); err != nil {
		return
	}
	data.User = user

	// get session
	sess, err := GetSession(c)
	if err != nil {
		return err
	}

	ok, err := user.Auth()
	if err != nil {
		return err
	}

	if !ok {
		data.Error = map[string]string{"Password": "邮箱或密码错误"}
		sess.AddFlash(data)
		sess.Save(c.Request(), c.Response())
		return c.Redirect(http.StatusFound, "/login")
	}

	if err = user.SignIn(c); err != nil {
		return
	}

	return c.Redirect(http.StatusFound, "/")
}

// GetForgetPassword
func GetForgetPassword(c echo.Context) (err error) {

	IfLoginedRedirectTo(c, "/")

	data := ForgetPasswordFlash{
		Title:    serve.Conf.Site.Name + "-" + "忘记密码",
		SiteName: serve.Conf.Site.Name,
		Error:    map[string]string{},
		User:     &m.User{},
	}

	// get session
	sess, err := GetSession(c)
	if err != nil {
		return err
	}

	if flashes := sess.Flashes(); len(flashes) > 0 {
		data = flashes[0].(ForgetPasswordFlash)
	}
	sess.Save(c.Request(), c.Response())

	return c.Render(http.StatusOK, "forget_password", data)
}

// PostForgetPassword
func PostForgetPassword(c echo.Context) (err error) {

	IfLoginedRedirectTo(c, "/")

	data := ForgetPasswordFlash{
		Title:    serve.Conf.Site.Name + "-" + "忘记密码",
		SiteName: serve.Conf.Site.Name,
		Error:    map[string]string{},
		User:     &m.User{},
	}

	if err = c.Bind(data.User); err != nil {
		return
	}

	ok, err := data.User.Exist()
	if err != nil {
		return err
	}

	// get session
	sess, err := GetSession(c)
	if err != nil {
		return err
	}

	if !ok {
		data.Error["Email"] = "Email不存在"
		sess.AddFlash(data)
		sess.Save(c.Request(), c.Response())
		return c.Redirect(http.StatusFound, "/forget_password")
	}

	return c.Render(http.StatusOK, "forget_password_ok", data)
}
