package repo

import (
	"bytes"
	"errors"
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

func Check(root string, path string, wg *sync.WaitGroup) (*Repo, error) {
	repo := NewRepo(path)

	if !repo.Exists {
		return nil, errors.New("No git repository found")
	}

	countOut = executeCommand(path, "git", "status", "--porcelain")
	modifiedLines := strings.Split(string(countOut), "\n")
	modified := len(modifiedLines) - 1

	if err != nil {
		return nil, output.Error(err.Error())
	}

	if modified > 0 && modifiedLines[0] != "" {
		repo.Changes = append(repo.Changes, output.ApplyColour(fmt.Sprintf(" M(%d)", modified), "red"))
	}

	repo.Branch = executeCommand(path, "git", "rev-parse", "--abbrev-ref", "HEAD")
	local_ref = executeCommand(path, "git", "rev-parse", repo.Branch)
	remote_ref := executeCommand(path, "git", "rev-parse", fmt.Sprintf("origin/%s", repo.Branch))

	if err == nil && remote_ref != local_ref {
		repo.Changes = append(repo.Changes, output.ApplyColour(" P", "blue"))
	}

	wg.Done()
	return repo, nil
}

func executeCommand(dir, string, commands ...string) (error, string) {
	cmd := exec.Command(commands)
	cmd.Dir = dir
	stdout, err := cmd.Output()

	if err != nil {
		errors.New(fmt.Sprintf("Error running command: %s", err.Error()))
	}

	return strings.TrimSpace(string(stdout[:]))
}
