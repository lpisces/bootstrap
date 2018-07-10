package c

import (
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/lpisces/bootstrap/cmd/serve"
	"github.com/lpisces/bootstrap/cmd/serve/mvc/m"
	"net/http"
)

// GetRegister register page
func GetRegister(c echo.Context) (err error) {

	if _, err := CurrentUser(c); err == nil {
		return c.Redirect(http.StatusFound, "/")
	}

	type Data struct {
		Title    string
		SiteName string
		Error    map[string]string
		User     *m.User
		Checked  bool
	}

	data := Data{
		Title:    serve.Conf.Site.Name + "-" + "注册",
		Checked:  false,
		SiteName: serve.Conf.Site.Name,
	}

	return c.Render(http.StatusOK, "register", data)
}

// PostRegister handle register request
func PostRegister(c echo.Context) (err error) {

	type Data struct {
		Title    string
		SiteName string
		Error    map[string]string
		User     *m.User
		Checked  bool
	}

	data := Data{
		Title:    serve.Conf.Site.Name + "-" + "注册",
		Checked:  false,
		SiteName: serve.Conf.Site.Name,
	}

	user := new(m.User)
	if err = c.Bind(user); err != nil {
		return
	}

	data.User = user
	if ok, errs := user.Validate(); !ok {
		log.Info(errs)
		data.Error = errs
		data.Checked = true
		return c.Render(http.StatusOK, "register", data)
	}

	if err := user.Create(); err != nil {
		data.Error = map[string]string{"Email": err.Error()}
		data.Checked = true
		log.Info(data.Error)
		return c.Render(http.StatusOK, "register", data)
	}

	return c.Redirect(http.StatusMovedPermanently, "/login")
}

// GetLogin login page
func GetLogin(c echo.Context) (err error) {

	if _, err := CurrentUser(c); err == nil {
		return c.Redirect(http.StatusFound, "/")
	}

	type Data struct {
		Title    string
		SiteName string
		Error    map[string]string
		User     *m.User
		Checked  bool
	}

	data := Data{
		Title:    serve.Conf.Site.Name + "-" + "登录",
		SiteName: serve.Conf.Site.Name,
		Error:    map[string]string{},
		User:     &m.User{},
		Checked:  false,
	}

	return c.Render(http.StatusOK, "login", data)
}

// PostLogin handle login request
func PostLogin(c echo.Context) (err error) {

	type Data struct {
		Title    string
		SiteName string
		Error    map[string]string
		User     *m.User
		Checked  bool
	}

	data := Data{
		Title:    serve.Conf.Site.Name + "-" + "登录",
		SiteName: serve.Conf.Site.Name,
		Error:    map[string]string{},
		User:     &m.User{},
		Checked:  true,
	}

	user := new(m.User)
	if err = c.Bind(user); err != nil {
		return
	}
	data.User = user

	ok, err := user.Auth()
	if err != nil {
		return err
	}

	if !ok {
		data.Error = map[string]string{"Password": "邮箱或密码错误"}
		log.Info(data.Error)
		return c.Render(http.StatusOK, "login", data)
	}

	if err = user.SignIn(c); err != nil {
		return
	}

	return c.Redirect(http.StatusFound, "/")
}
