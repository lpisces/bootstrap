package controller

import (
	"github.com/labstack/echo"
	"net/http"
)

type User struct {
}

// GetRegister register page
func GetRegister(c echo.Context) (err error) {
	return c.Render(http.StatusOK, "register", nil)
}

// PostRegister handle register request
func PostRegister(c echo.Context) (err error) {
	return c.String(http.StatusOK, "OK")
}
