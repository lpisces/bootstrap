package mvc

import (
	"github.com/labstack/echo"
	"github.com/lpisces/bootstrap/cmd/serve/mvc/c"
)

func Route(e *echo.Echo) {

	// home
	e.GET("/", c.HomeHandler)

	// hello
	e.GET("/hello", c.HelloHandler)

	// register
	e.GET("/register", c.GetRegister)
	e.POST("/register", c.PostRegister)

	// login
	e.GET("/login", c.GetLogin)

}
