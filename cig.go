package main

import (
	"os"

	"github.com/stevenjack/cig/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/stevenjack/cig/app"
	"github.com/stevenjack/cig/output"
)

const version string = "0.1.1"

func main() {
	var output_channel = make(chan output.Payload)
	go output.Wait(output_channel)

	cli_wrapper := main_app()
	repo_list, err := app.Config()

	if err != nil {
		output_channel <- output.Error(err.Error())
	}

	cli_wrapper.Action = func(context *cli.Context) {
		project_type := context.String("type")
		filter := context.String("filter")
		app.Handle(repo_list, project_type, filter, output_channel)
	}

	cli_wrapper.Run(os.Args)
}

func main_app() *cli.App {
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
	}
	return app
}
