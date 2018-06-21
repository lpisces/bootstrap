package serve

import (
	"gopkg.in/urfave/cli.v1"
	"log"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net"
	"fmt"
	"github.com/lpisces/bootstrap/cmd/serve/mvc/controller"
)

func Run(c *cli.Context) {

	e := echo.New()

	// public
	e.Static("/public", "public")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", controller.Hello)

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


