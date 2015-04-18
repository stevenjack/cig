package main

import (
	"bytes"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
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
		path := fmt.Sprintf("%s/.cig.yaml", home_dir)

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
	exists, err := exists(fmt.Sprintf("%v/.git", path))

	if exists {
		modified_files := exec.Command("git", "status", "-s")
		modified_files.Dir = path
		count := exec.Command("wc", "-l")
		count.Dir = path

		stdout, _ := modified_files.StdoutPipe()
		modified_files.Start()
		count.Stdin = stdout

		count_out, _ := count.Output()

		if err != nil {
			println(err.Error())
			return
		}

		modified, _ := strconv.ParseInt(strings.TrimSpace(string(count_out[:])), 0, 64)

		changes := []string{}

		if modified > 0 {
			changes = append(changes, color.RedString(fmt.Sprintf(" M(%d)", modified)))
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
			changes = append(changes, color.BlueString(" P"))
		}

		if len(changes) > 0 {
			var buffer bytes.Buffer

			repo_name := strings.Replace(path, root+"/", "", -1)

			buffer.WriteString(fmt.Sprintf("- %s (%s)", repo_name, branch_name))
			for _, change := range changes {
				buffer.WriteString(change)
			}
			channel <- buffer.String() + "\n"
		}

	}
	wg.Done()
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
