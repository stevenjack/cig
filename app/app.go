package app

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/stevenjack/cig/output"
	"github.com/stevenjack/cig/repo"
)

func Handle(repoList map[string]string, projectTypeToCheck string, filter string, output_channel chan output.Payload) {
	var wg sync.WaitGroup

	for projectType, path := range repoList {
		if projectTypeToCheck == "" || projectTypeToCheck == projectType {
			output_channel <- output.Print(fmt.Sprintf("\nChecking '%s' (%s) repos...", projectType, path))

			visit := func(visitedPath string, info os.FileInfo, err error) error {
				matched, _ := regexp.MatchString(filter, visitedPath)
				if info.IsDir() && (filter == "" || matched) {
					wg.Add(1)
					go repo.Check(path, visitedPath, output_channel, &wg)
				}
				return nil
			}

			err := filepath.Walk(path, visit)
			if err != nil {
				output_channel <- output.FatalError(err.Error())
			}
		}

		wg.Wait()
	}
	wg.Wait()
}
