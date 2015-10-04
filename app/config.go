package app

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/stevenjack/cig/Godeps/_workspace/src/github.com/mitchellh/go-homedir"
	"github.com/stevenjack/cig/Godeps/_workspace/src/gopkg.in/yaml.v2"
)

func Config(configPath string) (map[string]string, error) {
        repoList := make(map[string]string)
        homeDir, err := homedir.Dir()

        if configPath != "" {
                homeDir = configPath
	}

	if err != nil {
		return nil, errors.New("Couldn't determine home directory")
	}

        path := filepath.Join(homeDir, ".cig.yaml")
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Can't find config '%s'", path))
	}

        err = yaml.Unmarshal([]byte(data), &repoList)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Problem parsing '%s', please check documentation", path))
	}

        return repoList, nil
}
