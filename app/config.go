package app

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/stevenjack/cig/Godeps/_workspace/src/gopkg.in/yaml.v2"
)

func Config(path string) (map[string]string, error) {
	repo_list := make(map[string]string)
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
