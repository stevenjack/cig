package main

import (
	"bytes"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
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

	app.Action = func(c *cli.Context) {
		paths := make(map[interface{}]interface{})
		home_dir, err := homedir.Dir()
		path := fmt.Sprintf("%s/.cig.yaml", home_dir)

		data, err := ioutil.ReadFile(path)
		check(err)

		err = yaml.Unmarshal([]byte(data), &paths)
		check(err)

		var channel = make(chan string)
		var wg sync.WaitGroup

		go output(channel)

		for k, v := range paths {
			files, _ := ioutil.ReadDir(v.(string))
			fmt.Printf("\nChecking '%s' repos...\n", k)
			for _, f := range files {
				if f.IsDir() {
					wg.Add(1)
					go checkRepo(v.(string)+"/"+f.Name(), channel, &wg)
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

func checkRepo(path string, channel chan string, wg *sync.WaitGroup) {
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

			buffer.WriteString(fmt.Sprintf("- %s (%s)", path, branch_name))
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
