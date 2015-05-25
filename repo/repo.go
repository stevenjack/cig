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

func Check(root string, path string, output_channel chan output.Payload, wg *sync.WaitGroup) {
	exists, err := Exists(filepath.Join(path, ".git"))

	if err != nil {
		return
	}

	if exists {
		modified_files := exec.Command("git", "status", "--porcelain")
		modified_files.Dir = path

		count_out, err := modified_files.Output()
		modified_lines := strings.Split(string(count_out), "\n")
		modified := len(modified_lines) - 1

		if err != nil {
			output_channel <- output.Error(err.Error())
		}

		changes := []string{}

		if modified > 0 && modified_lines[0] != "" {
			changes = append(changes, output.ApplyColour(fmt.Sprintf(" M(%d)", modified), "red"))
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
			changes = append(changes, output.ApplyColour(" P", "blue"))
		}

		if len(changes) > 0 {
			var buffer bytes.Buffer

			repo_name := strings.Replace(path, fmt.Sprintf("%s%s", root, string(os.PathSeparator)), "", -1)

			buffer.WriteString(fmt.Sprintf("- %s (%s)", repo_name, branch_name))
			for _, change := range changes {
				buffer.WriteString(change)
			}
			output_channel <- output.Print(buffer.String())
		}

	}
	wg.Done()
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
