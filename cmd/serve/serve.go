package serve

import (
	"gopkg.in/urfave/cli.v1"
	"log"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"net"
	"fmt"
)

func Run(c *cli.Context) {

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)

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

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
