package c

import (
	"github.com/labstack/echo"
	"net/http"
)

func GetHome(c echo.Context) error {
	return c.String(http.StatusOK, "Hello!")
}
