package app

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/mitchellh/go-homedir"
	"github.com/stevenjack/cig/output"
	"github.com/stevenjack/cig/repo"
)

func Handle(repoList map[string]string, projectTypeToCheck string, filter string, outputChannel chan output.Payload) {
	var wg sync.WaitGroup

	for projectType, path := range repoList {
		if projectTypeToCheck == "" || projectTypeToCheck == projectType {
			outputChannel <- output.Print(fmt.Sprintf("\nChecking '%s' (%s) repos...", projectType, path))
			path = resolvePath(path)

			visit := func(visitedPath string, info os.FileInfo, err error) error {
				if err != nil {
					outputChannel <- output.Error(fmt.Sprintf("- %s", err.Error()))
					return nil
				}
				matched, _ := regexp.MatchString(filter, visitedPath)
				if info.IsDir() && (filter == "" || matched) {
					wg.Add(1)
					go repo.Check(path, visitedPath, outputChannel, &wg)
				}
				return nil
			}

			err := filepath.Walk(path, visit)
			if err != nil {
				outputChannel <- output.Error(err.Error())
			}
		}

		wg.Wait()
	}
	wg.Wait()
}

func resolvePath(path string) string {
	hasTilde := strings.HasPrefix(path, "~")
	if hasTilde {
		homeDir, err := homedir.Dir()
		if err == nil {
			return strings.Replace(path, "~", homeDir, -1)
		}
	}
	return path
}
