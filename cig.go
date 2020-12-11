package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/stevenjack/cig/app"
	"github.com/stevenjack/cig/output"
)

const version string = "0.1.5"

func main() {
	var outputChannel = make(chan output.Payload)

	go output.Wait(outputChannel)
	cliWrapper := mainApp()

	cliWrapper.Action = func(context *cli.Context) {
		configPath := context.String("config-path")
		projectType := context.String("type")
		filter := context.String("filter")
		repoList, err := app.Config(configPath)

		if err != nil {
			outputChannel <- output.FatalError(err.Error())
		}

		app.Handle(repoList, projectType, filter, outputChannel)
	}

	cliWrapper.Run(os.Args)
}

func mainApp() *cli.App {
	app := cli.NewApp()
	app.Name = "cig"
	app.Usage = "cig (Can I go?) checks all your git repos to see if they're in the state you want them to be"
	app.Version = version

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "filter, f",
			Value: "",
			Usage: "Filter repos being searched",
		},
		cli.StringFlag{
			Name:  "type, t",
			Value: "",
			Usage: "Filter by type",
		},
		cli.StringFlag{
			Name:  "config-path, cp",
			Value: "",
			Usage: "Path to the cig config (Default ~/.cig.yaml)",
		},
	}
	return app
}
