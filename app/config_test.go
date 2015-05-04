package app

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stevenjack/cig/app"
)

func TestMapReturnedFromCigYaml(t *testing.T) {
	repoList, err := config("cig.yaml")

	if err != nil {
		t.Error(fmt.Sprintf("return of app.Config should be of type 'map[string]string, got: %s", repoList))
	}
}

func TestErrorRaisedWhenConfigDoesNotExist(t *testing.T) {
	_, err := config("cig_doesnt_exist.yaml")

	if err == nil {
		t.Error("app.Config should raise error with invalid path")
	}
}

func TestMalformedYamlRasiesError(t *testing.T) {
	_, err := config("invalid_cig.yaml")

	if err == nil {
		t.Error("app.Config should raise error with invalid yaml")
	}
}

func config(filename string) (map[string]string, error) {
	cwd, _ := os.Getwd()
	path := filepath.Join(cwd, "..", "test", "fixtures", filename)
	return app.Config(path)
}
