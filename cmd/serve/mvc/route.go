package mvc

import (
	"github.com/labstack/echo"
	"github.com/lpisces/bootstrap/cmd/serve/mvc/c"
)

func Route(e *echo.Echo) {

	// home
	e.GET("/", c.GetHome)

	// hello
	e.GET("/hello", c.HelloHandler)

	e.GET("/admin", c.HelloHandler)

	// register
	e.GET("/register", c.GetRegister)
	e.POST("/register", c.PostRegister)

	// login
	e.GET("/login", c.GetLogin)
	e.POST("/login", c.PostLogin)

}
