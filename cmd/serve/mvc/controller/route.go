package controller

import (
	"github.com/labstack/echo"
)

func Route(e *echo.Echo) {
	e.GET("/", Hello)
	return
}