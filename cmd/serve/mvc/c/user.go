package c

import (
	"github.com/labstack/echo"
	//"github.com/labstack/gommon/log"
	"github.com/lpisces/bootstrap/cmd/serve"
	"github.com/lpisces/bootstrap/cmd/serve/mvc/m"
	"net/http"
)

type User struct {
}

// GetRegister register page
func GetRegister(c echo.Context) (err error) {
	type Data struct {
		SiteConfig *serve.SiteConfig
	}
	data := Data{SiteConfig: serve.Conf.Site}
	return c.Render(http.StatusOK, "register", data)
}

// PostRegister handle register request
func PostRegister(c echo.Context) (err error) {
	user := new(m.User)
	if err = c.Bind(user); err != nil {
		return
	}
	if err = user.Validate(); err != nil {
		return
	}
	return c.String(http.StatusOK, "OK")
}
