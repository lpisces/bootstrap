package c

import (
	"github.com/labstack/echo"
	"net/http"
)

func GetAdmin(c echo.Context) error {
	return c.String(http.StatusOK, "admin")
}
