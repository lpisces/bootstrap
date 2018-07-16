package mvc

import (
	"github.com/labstack/echo"
	"github.com/lpisces/bootstrap/cmd/serve/mvc/c"
)

func Route(e *echo.Echo) {

	// home
	e.GET("/", c.GetHome)

	// hello
	e.GET("/admin", c.GetAdmin)

	// register
	e.GET("/register", c.GetRegister)
	e.POST("/register", c.PostRegister)

	// login
	e.GET("/login", c.GetLogin)
	e.POST("/login", c.PostLogin)

	// logout
	e.GET("/logout", c.GetLogout)

	// forget_password
	e.GET("forget_password", c.GetForgetPassword)
	e.POST("forget_password", c.PostForgetPassword)

}
