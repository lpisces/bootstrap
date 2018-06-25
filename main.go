package main

import (
	"gopkg.in/urfave/cli.v1"
	"os"
	"log"
	"github.com/lpisces/bootstrap/cmd/serve"
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
			Action: serve.Run,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "debug, d",
					Usage: "show debug info",
				},
				cli.StringFlag{
					Name: "port, p",
					Usage: "listen port",
					Value: "1323",
				},
				cli.StringFlag{
					Name: "bind, b",
					Usage: "bind host",
					Value: "127.0.0.1",
				},
				cli.StringFlag{
					Name: "config, c",
					Usage: "load config file",
					Value: "./config.ini",
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}