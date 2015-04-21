package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"

	"github.com/stevenjack/cig/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/stevenjack/cig/Godeps/_workspace/src/github.com/fatih/color"
	"github.com/stevenjack/cig/Godeps/_workspace/src/github.com/mitchellh/go-homedir"
	"github.com/stevenjack/cig/Godeps/_workspace/src/gopkg.in/yaml.v2"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "cig"
	app.Usage = "cig (Can I go?) checks all your git repos to see if they're in the state you want them to be"
	app.Version = "0.1.0"

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
		path := fmt.Sprintf("%s%s.cig.yaml\n", home_dir, string(os.PathSeparator))

		data, err := ioutil.ReadFile(path)

		if err != nil {
			channel <- color.RedString(fmt.Sprintf("Can't find config '%s'", path))
			os.Exit(-1)
		}

		err = yaml.Unmarshal([]byte(data), &paths)
		if err != nil {
			channel <- color.RedString(fmt.Sprintf("Problem parsing '%s', please check documentation", path))
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

func output(channel chan string) {
	for {
		entry := <-channel
		fmt.Printf(entry)
	}
}

func checkRepo(root string, path string, channel chan string, wg *sync.WaitGroup) {
	exists, err := exists(fmt.Sprintf("%s%s.git", path, string(os.PathSeparator)))

	if exists {
		modified_files := exec.Command("git", "status", "-s")
		modified_files.Dir = path

		count_out, _ := modified_files.Output()
		modified_lines := strings.Split(string(count_out), "\n")
		modified := len(modified_lines) - 1

		if err != nil {
			println(err.Error())
			return
		}

		changes := []string{}

		if modified > 0 && modified_lines[0] != "" {
			changes = append(changes, print_output(fmt.Sprintf(" M(%d)", modified), "red"))
		}

		branch := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
		branch.Dir = path
		bstdout, _ := branch.Output()
		branch_name := strings.TrimSpace(string(bstdout[:]))

		local := exec.Command("git", "rev-parse", branch_name)
		local.Dir = path
		lstdout, _ := local.Output()
		local_ref := strings.TrimSpace(string(lstdout[:]))

		remote := exec.Command("git", "rev-parse", fmt.Sprintf("origin/%s", branch_name))
		remote.Dir = path
		rstdout, err := remote.Output()
		remote_ref := strings.TrimSpace(string(rstdout[:]))

		if err == nil && remote_ref != local_ref {
			changes = append(changes, print_output(" P", "blue"))
		}

		if len(changes) > 0 {
			var buffer bytes.Buffer

			repo_name := strings.Replace(path, fmt.Sprintf("%s%s", root, string(os.PathSeparator)), "", -1)

			buffer.WriteString(fmt.Sprintf("- %s (%s)", repo_name, branch_name))
			for _, change := range changes {
				buffer.WriteString(change)
			}
			channel <- buffer.String() + "\n"
		}

	}
	wg.Done()
}

func print_output(message string, output_type string) string {
	if runtime.GOOS != "windows" {
		switch output_type {
		case "red":
			return color.RedString(message)
		case "blue":
			return color.BlueString(message)
		}

	}
	return message
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
