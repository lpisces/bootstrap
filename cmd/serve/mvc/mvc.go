package mvc

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/lpisces/bootstrap/cmd/serve"
	"github.com/lpisces/bootstrap/cmd/serve/mvc/model"
	"gopkg.in/urfave/cli.v1"
	"html/template"
	"net"
)

// serve start web server
func startSrv() (err error) {

	// migrate db
	if err := model.Migrate(); err != nil {
		log.Fatal(err)
	}

	// get config
	config := serve.Conf

	// new echo instance
	e := echo.New()

	// set template render
	t := &Template{
		templates: template.Must(template.ParseGlob("cmd/serve/mvc/view/*.html")),
	}
	e.Renderer = t

	// public
	e.Static("/public", "public")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	Route(e)

	// Start server
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", config.Srv.Host, config.Srv.Port))
	if err != nil {
		return err
	}

	e.Listener = l

	e.HideBanner = true

	if serve.Debug {
		e.Logger.SetLevel(log.DEBUG)
	}

	e.Logger.Infof("http server started on %s:%s in %s model", config.Srv.Host, config.Srv.Port, config.Mode)
	e.Logger.Fatal(e.Start(""))
	return
}

func Run(c *cli.Context) (err error) {

	// Load default config
	config := serve.DefaultConfig()

	// override default config
	configFilePath := c.String("config")
	if configFilePath != "" {
		if err := config.Load(configFilePath); err != nil {
			log.Fatal(err)
		}
	}

	// flag override ini file config
	bind := c.String("bind")
	if bind != "" {
		config.Srv.Host = bind
	}

	port := c.String("port")
	if port != "" {
		config.Srv.Port = port
	}

	env := c.String("env")
	if env != "" {
		config.Mode = env
	}

	serve.Conf = config

	// run mode
	if config.Mode != "production" {
		serve.Debug = true
	}

	// start server
	err = startSrv()
	if err != nil {
		log.Fatal(err)
	}
	return
}
