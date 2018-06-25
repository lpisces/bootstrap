package controller

import (
	"github.com/labstack/echo"
	"net/http"
)

func Hello(c echo.Context) error {
	//var user model.User
	//log.Print(user)
	return c.String(http.StatusOK, "Hello, World!")
}
