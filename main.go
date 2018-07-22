package main

import (
	pcli "github.com/lpisces/pcheck/cmd/cli"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "pcheck"
	app.Usage = "package check"

	app.Commands = []cli.Command{
		{
			Name:   "cli",
			Usage:  "execute command line",
			Action: pcli.Run,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "env, e",
					Usage: "set env",
					Value: "dev",
				},
				cli.StringFlag{
					Name:  "config, c",
					Usage: "load config file",
				},
				cli.StringFlag{
					Name:  "source, s",
					Usage: "data source",
					Value: "./package.csv",
				},
				cli.StringFlag{
					Name:  "output, o",
					Usage: "check result output",
					Value: "./result.csv",
				},
				cli.StringFlag{
					Name:  "cache",
					Usage: "cache dir",
					Value: "./cache",
				},
				cli.StringFlag{
					Name:  "key, k",
					Usage: "kuaidi100 key",
				},
				cli.StringFlag{
					Name:  "customer",
					Usage: "kuaidi100 customer id",
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
