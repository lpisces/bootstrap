package mvc

import (
	"github.com/labstack/echo"
	"github.com/lpisces/bootstrap/cmd/serve/mvc/c"
)

func Route(e *echo.Echo) {

	// home
	e.GET("/", c.GetHome)

}
