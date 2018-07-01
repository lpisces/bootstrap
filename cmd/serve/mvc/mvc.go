package mvc

import (
	"fmt"
	"github.com/gobuffalo/packr"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/lpisces/bootstrap/cmd/serve"
	"github.com/lpisces/bootstrap/cmd/serve/mvc/m"
	"gopkg.in/urfave/cli.v1"
	"html/template"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// serve start web server
func startSrv() (err error) {

	// migrate db
	if err := m.Migrate(); err != nil {
		log.Fatal(err)
	}

	// get config
	config := serve.Conf

	// new echo instance
	e := echo.New()

	// set template render
	viewPath := "cmd/serve/mvc/v"
	box := packr.NewBox("./v")
	templates := template.New("")

	if _, err := os.Stat(viewPath); err == nil {
		err = filepath.Walk(viewPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				relativePath := strings.Replace(path, viewPath+"/", "", 1)
				//log.Info(box.String(relativePath))
				_, err := templates.Parse(box.String(relativePath))
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			e.Logger.Fatal(err)
		}
	}

	t := &Template{
		templates: templates,
	}
	e.Renderer = t

	// public
	//e.Static("/public", "public")
	//staticPath := "public"

	if _, err := os.Stat("public"); err == nil {
		publicBox := packr.NewBox("../../../public")
		assetHandler := http.FileServer(publicBox)
		e.GET("/public/*", echo.WrapHandler(http.StripPrefix("/public/", assetHandler)))
	}

	// Middleware
	e.Use(middleware.Logger())
	//e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	//e.Use(middleware.CSRF())

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
