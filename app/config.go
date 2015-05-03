package app

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/stevenjack/cig/Godeps/_workspace/src/github.com/mitchellh/go-homedir"
	"github.com/stevenjack/cig/Godeps/_workspace/src/gopkg.in/yaml.v2"
)

func Config() (map[string]string, error) {
	repo_list := make(map[string]string)
	home_dir, err := homedir.Dir()

	if err != nil {
		return nil, errors.New("Couldn't determine home directory")
	}

	path := filepath.Join(home_dir, ".cig.yaml")
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Can't find config '%s'", path))
	}

	err = yaml.Unmarshal([]byte(data), &repo_list)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Problem parsing '%s', please check documentation", path))
	}

	return repo_list, nil
}
