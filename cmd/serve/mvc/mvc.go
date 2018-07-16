package mvc

import (
	"fmt"
	"github.com/gobuffalo/packr"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/lpisces/bootstrap/cmd/serve"
	"github.com/lpisces/bootstrap/cmd/serve/mvc/m"
	//"github.com/lpisces/bootstrap/cmd/serve/mw"
	"gopkg.in/urfave/cli.v1"
	"html/template"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Run run app
func Run(c *cli.Context) (err error) {

	// init config
	if err = InitConfig(c); err != nil {
		return
	}

	// init db
	if err = InitDB(); err != nil {
		return
	}

	// start server
	if err = startSrv(); err != nil {
		return
	}

	return
}

// serve start web server
func startSrv() (err error) {
	// get config
	config := serve.Conf

	// init echo instance
	e, err := InitEcho()
	if err != nil {
		return err
	}

	// set template render, use packr if not development mode
	t, err := InitTemplate("cmd/serve/mvc/v")
	if err != nil {
		return
	}
	e.Renderer = t

	// static file
	if err := InitStaticServer(e, "public"); err != nil {
		return err
	}

	// Routes
	Route(e)

	e.Logger.Infof("http server started on %s:%s in %s mode", config.Srv.Host, config.Srv.Port, config.Mode)
	e.Logger.Fatal(e.Start(""))
	return
}

// InitEcho create echo instance
func InitEcho() (e *echo.Echo, err error) {
	// get config
	config := serve.Conf

	// new echo instance
	e = echo.New()

	// log
	if serve.Debug {
		e.Logger.SetLevel(log.DEBUG)
	} else {
		e.Logger.SetLevel(log.ERROR)
	}

	// Middleware
	e.Use(middleware.Logger())
	if !serve.Debug {
		e.Use(middleware.Recover())
	}
	e.Use(middleware.Gzip())
	//e.Use(middleware.CSRF())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(config.Secret.Session))))
	//e.Use(mw.CasbinAuth)

	// Start server
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", config.Srv.Host, config.Srv.Port))
	if err != nil {
		return e, err
	}

	e.Listener = l
	e.HideBanner = true

	return
}

// InitConfig init config related var
func InitConfig(c *cli.Context) (err error) {
	// Load default config
	config := serve.Conf

	// ini config file override default config
	configFile := c.String("config")
	if configFile != "" {
		if err := config.Load(configFile); err != nil {
			return err
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

	return
}

// InitDB run db migration
func InitDB() (err error) {
	// migrate db
	if err = m.Migrate(); err != nil {
		return
	}
	return
}

// InitTemplate
func InitTemplate(path string) (t *Template, err error) {
	// set template render, use packr if not development mode
	templates := template.New("")
	viewPath := path
	box := packr.NewBox("./v")

	if _, err := os.Stat(viewPath); err == nil {
		err = filepath.Walk(viewPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				relativePath := strings.Replace(path, viewPath+"/", "", 1)
				if serve.Embed {
					_, err := templates.Parse(box.String(relativePath))
					if err != nil {
						return err
					}
				} else {
					_, err := templates.ParseGlob(path)
					if err != nil {
						return err
					}
				}
			}
			return nil
		})
		if err != nil {
			//e.Logger.Fatal(err)
			return nil, err
		}
	}

	t = &Template{
		templates: templates,
	}
	return
}

func InitStaticServer(e *echo.Echo, path string) (err error) {
	// static file serve, use packr if not development
	if !serve.Embed {
		e.Static("/public", path)
		e.File("/favicon.ico", path+"/favicon.ico")
	} else {
		if _, err := os.Stat(path); err == nil {
			publicBox := packr.NewBox("../../../" + path)
			assetHandler := http.FileServer(publicBox)
			e.GET("/public/*", echo.WrapHandler(http.StripPrefix("/public/", assetHandler)))
			e.GET("/favicon.ico", echo.WrapHandler(http.StripPrefix("/", assetHandler)))
		} else {
			return err
		}
	}
	return
}
