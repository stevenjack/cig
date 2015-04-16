package main

import (
	"flag"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	"github.com/libgit2/git2go"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"sync"
  "strings"
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
		usr, _ := user.Current()
		dir := usr.HomeDir
		path := dir + "/cig.yaml"

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
	repoPath := flag.String("repo"+path, path, "path to the git repository")
	flag.Parse()
	repo, err := git.OpenRepository(*repoPath)

	opts := &git.StatusOptions{}
	opts.Show = git.StatusShowIndexAndWorkdir
	opts.Flags = git.StatusOptIncludeUntracked | git.StatusOptRenamesHeadToIndex | git.StatusOptSortCaseSensitively

	exists, err := exists(fmt.Sprintf("%v/.git", path))

	if exists {

		modified_files := exec.Command("git", "ls-files")
    modified_files.Dir = path

		stdout, err := modified_files.Output()

		if err != nil {
			println(err.Error())
			return
		}

    split_strings := strings.Split(fmt.Sprintf("%s", stdout), "\n")

    channel <- fmt.Sprintf("Modified files: %s\n", split_strings)
	}

	if err == nil {
		statusList, err := repo.StatusList(opts)
		check(err)

		entryCount, err := statusList.EntryCount()
		check(err)

		currentBranch, err := repo.Head()
		r := regexp.MustCompile("refs/heads/([/a-z-0-9_]+)")
		branch := r.FindStringSubmatch(currentBranch.Name())[1]

		_, local, err := repo.RevparseExt(branch)
		_, remote, err := repo.RevparseExt(fmt.Sprintf("origin/%v", branch))

		changes := []string{}

		if local != nil && remote != nil && local.Target().String() != remote.Target().String() {
			changes = append(changes, color.BlueString(" P"))
		}

		if entryCount > 0 {
			changes = append(changes, color.RedString(fmt.Sprintf(" M(%v)", entryCount)))
		}

		if len(changes) > 0 {
			channel <- fmt.Sprintf("- %v (%v)", path, branch)
			for _, change := range changes {
				channel <- change
			}
			channel <- "\n"
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
