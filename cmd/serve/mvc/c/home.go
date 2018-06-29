package c

import (
	"github.com/labstack/echo"
	"net/http"
)

func HomeHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "hello", "World")
}
