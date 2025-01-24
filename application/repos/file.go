package repos

import (
	"io/ioutil"

	"git.sr.ht/~icikowski/myapps/common"
	"git.sr.ht/~icikowski/myapps/types"
	"gopkg.in/yaml.v3"
)

// FileRepositoryLoader loads applications repository from file
type FileRepositoryLoader struct{}

// LoadRepository implements RepositoryLoader
func (*FileRepositoryLoader) LoadRepository(source string) (types.Repository, bool) {
	var repository types.Repository

	contents, err := ioutil.ReadFile(source)
	if err != nil {
		common.PrintErrorWhileMsg("reading file", source, err)
		return repository, false
	}

	if err := yaml.Unmarshal(contents, &repository); err != nil {
		common.PrintErrorWhileMsg("parsing file", source, err)
		return repository, false
	}

	return repository, true
}

var _ RepositoryLoader = &FileRepositoryLoader{}
