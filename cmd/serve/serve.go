package serve

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/lpisces/bootstrap/cmd/serve/mvc/controller"
	"gopkg.in/urfave/cli.v1"
	"log"
	"net"
)



func Run(c *cli.Context) {

	e := echo.New()

	// public
	e.Static("/public", "public")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", controller.Welcome(c))

	// Start server
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", c.String("bind"), c.String("port")))
	if err != nil {
		e.Logger.Fatal(l)
	}

	e.Listener = l
	e.Logger.Fatal(e.Start(""))

	if c.Bool("debug") {
		log.Printf("http server started on %s:%s", c.String("bind"), c.String("port"))
	}
}


