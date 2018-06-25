package serve

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/urfave/cli.v1"
	"net"
	"github.com/lpisces/bootstrap/utils"
	"github.com/lpisces/bootstrap/cmd/serve/mvc/controller"
	"github.com/labstack/gommon/log"
	"github.com/jinzhu/gorm"
)

func Run(c *cli.Context) {

	e := echo.New()

	// public
	e.Static("/public", "public")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Load config
	var config utils.Config
	if err := config.Load(c); err != nil {
		e.Logger.Fatal(err)
	}

	// InitDB
	if conn, err := utils.GetConn(config); err != nil {
		e.Logger.Fatal(err)
	} else {
		utils.Conn = conn
		//defer utils.Conn.Close()
	}
	log.Print(utils.Conn)

	type User struct {
		gorm.Model
		Name string
		Email string
		PasswordDigest string
		Avatar string
	}
	utils.Conn.AutoMigrate(&User{})

	// Routes
	controller.Route(e)

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


