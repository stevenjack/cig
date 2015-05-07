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

type Repo struct {
	Path    string
	Exists  bool
	Branch  string
	Changes []string
	Synced  bool
}

func NewRepo(string path) *Repo {
	path := filepath.Join(path, ".git")
	var exists bool

	_, err := os.Stat(r.Path)
	if err == nil {
		exists = true
	}
	if os.IsNotExist(err) {
		exists = false
	}

	repo := Repo{path, exists}
	return &repo
}

func (r *Repo) SetChanges() error {

	modifiedFiles := exec.Command("git", "status", "--porcelain")
	modifiedFiles.Dir = r.Path

	countOut, _ := modifiedFiles.Output()
	modifiedLines := strings.Split(string(countOut), "\n")
	modified := len(modifiedLines) - 1

	if err != nil {
		return err.Error()
	}

	if modified > 0 && modifiedLines[0] != "" {
		r.Changes = append(Changes, output.ApplyColour(fmt.Sprintf(" M(%d)", modified), "red"))
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

}

func Check(root string, path string, outputChannel chan output.Payload, wg *sync.WaitGroup) {
	repo := Repo{}
	repo.SetPath(path)

	if !repo.Exists {
		return
	}

	modifiedFiles := exec.Command("git", "status", "--porcelain")
	modifiedFiles.Dir = path

	countOut, _ := modifiedFiles.Output()
	modifiedLines := strings.Split(string(countOut), "\n")
	modified := len(modifiedLines) - 1

	if err != nil {
		outputChannel <- output.Error(err.Error())
	}

	changes := []string{}

	if modified > 0 && modifiedLines[0] != "" {
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
		outputChannel <- output.Print(buffer.String())

	}
	wg.Done()
}
