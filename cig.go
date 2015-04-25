package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/stevenjack/cig/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/stevenjack/cig/Godeps/_workspace/src/github.com/mitchellh/go-homedir"
	"github.com/stevenjack/cig/Godeps/_workspace/src/gopkg.in/yaml.v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "cig"
	app.Usage = "cig (Can I go?) checks all your git repos to see if they're in the state you want them to be"
	app.Version = "0.1.1"

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

	app.Action = func(c *cli.Context) {
		project_type := c.String("type")

		var channel = make(chan string)
		go output(channel)

		paths := make(map[interface{}]interface{})
		home_dir, err := homedir.Dir()
		path := filepath.Join(home_dir, ".cig.yaml")

		data, err := ioutil.ReadFile(path)

		if err != nil {
			channel <- error_output(fmt.Sprintf("Can't find config '%s'\n", path))
			os.Exit(-1)
		}

		err = yaml.Unmarshal([]byte(data), &paths)

		if err != nil {
			channel <- error_output(fmt.Sprintf("Problem parsing '%s', please check documentation", path))
			os.Exit(-1)
		}
		check(err)

		var wg sync.WaitGroup

		for k, v := range paths {
			if project_type == "" || project_type == k {
				fmt.Printf("\nChecking '%s' (%s) repos...\n", k, v)

				visit := func(path string, info os.FileInfo, err error) error {
					filter := c.String("filter")

					matched, _ := regexp.MatchString(filter, path)
					if info.IsDir() && (filter == "" || matched) {
						wg.Add(1)
						go checkRepo(v.(string), path, channel, &wg)
					}
					return nil
				}

				err := filepath.Walk(v.(string), visit)
				if err != nil {
					log.Fatal(err)
				}
			}

			wg.Wait()
		}

		wg.Wait()
	}
	app.Run(os.Args)
}
