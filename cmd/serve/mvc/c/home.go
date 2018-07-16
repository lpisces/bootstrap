package c

import (
	"github.com/labstack/echo"
	"net/http"
)

var (
	store string
)

func GetHome(c echo.Context) error {
	return c.Render(http.StatusOK, "home", "home")
}
