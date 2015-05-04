package main

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/stevenjack/cig/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/stevenjack/cig/Godeps/_workspace/src/github.com/mitchellh/go-homedir"
	"github.com/stevenjack/cig/app"
	"github.com/stevenjack/cig/output"
)

const version string = "0.1.2"

func main() {
	var outputChannel = make(chan output.Payload)
	go output.Wait(outputChannel)

	cliWrapper := mainApp()
	path, err := configPath()

	if err != nil {
		outputChannel <- output.Error(err.Error())
	}

	repoList, err := app.Config(path)

	if err != nil {
		outputChannel <- output.Error(err.Error())
	}

	cliWrapper.Action = func(context *cli.Context) {
		projectType := context.String("type")
		filter := context.String("filter")
		app.Handle(repoList, projectType, filter, outputChannel)
	}

	cliWrapper.Run(os.Args)
}

func configPath() (string, error) {
	homeDir, err := homedir.Dir()

	if err != nil {
		return "", errors.New("Couldn't determine home directory")
	}

	path := filepath.Join(homeDir, ".cig.yaml")

	return path, nil
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
	}
	return app
}
