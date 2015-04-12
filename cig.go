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
	"os/user"
	"regexp"
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

		for _, v := range paths {
			files, _ := ioutil.ReadDir(v.(string))
			for i, f := range files {
				if f.IsDir() {
					go checkRepo(v.(string)+"/"+f.Name(), channel, i, len(files))
				}
			}
		}

		for {
			entry := <-channel
			if entry == "exit" {
				//				os.Exit(0)
			} else {
				fmt.Printf(entry)
			}
		}

	}

	app.Run(os.Args)
}

func checkRepo(path string, channel chan string, current int, size int) {
	repoPath := flag.String("repo"+path, path, "path to the git repository")
	flag.Parse()
	repo, err := git.OpenRepository(*repoPath)

	opts := &git.StatusOptions{}
	opts.Show = git.StatusShowIndexAndWorkdir
	opts.Flags = git.StatusOptIncludeUntracked | git.StatusOptRenamesHeadToIndex | git.StatusOptSortCaseSensitively

	if err == nil {
		statusList, err := repo.StatusList(opts)
		check(err)

		entryCount, err := statusList.EntryCount()
		check(err)

		currentBranch, err := repo.Head()
		r := regexp.MustCompile("refs/heads/([/a-z-0-9_]+)")
		branch := r.FindStringSubmatch(currentBranch.Name())[1]

		_, ref, err := repo.RevparseExt(branch)
		//		check(err)
		_, ref_two, err := repo.RevparseExt(fmt.Sprintf("origin/%v", branch))
		//		check(err)

		if ((ref != nil && ref_two != nil) && ref.Target().String() != ref_two.Target().String()) || entryCount > 0 {
			channel <- fmt.Sprintf("\n%v (%v)\n", path, branch)
		}

		if ref != nil && ref_two != nil {
			if ref.Target().String() != ref_two.Target().String() {
				channel <- color.RedString("Push to master needed\n")
			}
		}

		if entryCount > 0 {
			channel <- color.RedString(fmt.Sprintf("%v file(s) changed/staged\n", entryCount))
		}

		if (current + 1) == size {
			channel <- "exit"
		}
	}
}
