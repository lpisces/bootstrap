package controller

import (
	"github.com/labstack/echo"
	"net/http"
	"gopkg.in/urfave/cli.v1"
)

func Hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func Welcome(ctx *cli.Context)(fn func(c echo.Context)error) {
	return func (c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	}
}