package repo

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/stevenjack/cig/output"
)

func Check(root string, path string, outputChannel chan output.Payload, wg *sync.WaitGroup) {
	defer wg.Done()
	exists, err := Exists(filepath.Join(path, ".git"))

	if err != nil {
		return
	}

	if exists {
                modifiedFiles := exec.Command("git", "status", "--porcelain")
                modifiedFiles.Dir = path

                countOut, _ := modifiedFiles.Output()
                modifiedLines := strings.Split(string(countOut), "\n")
                modified := len(modifiedLines) - 1

		if err != nil {
                        outputChannel <- output.FatalError(err.Error())
		}

		changes := []string{}

                if modified > 0 && modifiedLines[0] != "" {
			changes = append(changes, output.ApplyColour(fmt.Sprintf(" M(%d)", modified), "red"))
		}

		branch := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
		branch.Dir = path
		bstdout, _ := branch.Output()
                branchName := strings.TrimSpace(string(bstdout[:]))

                local := exec.Command("git", "rev-parse", branchName)
		local.Dir = path
		lstdout, _ := local.Output()
                localRef := strings.TrimSpace(string(lstdout[:]))

                remote := exec.Command("git", "rev-parse", fmt.Sprintf("origin/%s", branchName))
		remote.Dir = path
		rstdout, err := remote.Output()
                remoteRef := strings.TrimSpace(string(rstdout[:]))

                if err == nil && remoteRef != localRef {
			changes = append(changes, output.ApplyColour(" P", "blue"))
		}

		if len(changes) > 0 {
			var buffer bytes.Buffer

                        repoName := strings.Replace(path, fmt.Sprintf("%s%s", root, string(os.PathSeparator)), "", -1)

                        buffer.WriteString(fmt.Sprintf("- %s (%s)", repoName, branchName))
			for _, change := range changes {
				buffer.WriteString(change)
			}
                        outputChannel <- output.Print(buffer.String())
		}

	}
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
