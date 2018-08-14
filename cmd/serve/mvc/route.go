package mvc

import (
	"github.com/dchest/captcha"
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

	// activate
	e.GET("/activate", c.GetActivate)

	// login
	e.GET("/login", c.GetLogin)
	e.POST("/login", c.PostLogin)

	// logout
	e.GET("/logout", c.GetLogout)

	// forget_password
	e.GET("/forget_password", c.GetForgetPassword)
	e.POST("/forget_password", c.PostForgetPassword)

	// reset password
	e.GET("/reset_password", c.GetResetPassword)
	e.POST("/reset_password", c.PostResetPassword)

	// captcha
	e.GET("/captcha/:id", func(c echo.Context) error {
		captcha.Server(captcha.StdWidth/2, captcha.StdHeight/2).ServeHTTP(c.Response(), c.Request())
		return nil
	})
}
