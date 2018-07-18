package main

import (
	"github.com/lpisces/bootstrap/cmd/serve/mvc"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
)

const (
	Embed = false
)

func main() {
	app := cli.NewApp()
	app.Name = "bootstrap"
	app.Usage = "bootstrap for website server development"

	app.Commands = []cli.Command{
		{
			Name:    "serve",
			Aliases: []string{"s"},
			Usage:   "start web server",
			Action:  mvc.Run,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "env, e",
					Usage: "set run env",
					Value: "development",
				},
				cli.StringFlag{
					Name:  "port, p",
					Usage: "listen port",
				},
				cli.StringFlag{
					Name:  "bind, b",
					Usage: "bind host",
				},
				cli.StringFlag{
					Name:  "config, c",
					Usage: "load config file",
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
