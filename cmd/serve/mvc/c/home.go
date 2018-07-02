package c

import (
	"github.com/labstack/echo"
	"net/http"
)

var (
	store string
)

func GetHome(c echo.Context) error {
	name := c.QueryParam("name")
	if name == "" {
		name = store
	} else {
		store = name
	}
	return c.String(http.StatusOK, name)
	//return c.Render(http.StatusOK, "home", "home")
}
