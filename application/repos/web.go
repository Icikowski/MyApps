package repos

import (
	"fmt"
	"net/http"

	"gopkg.in/yaml.v3"
	"icikowski.pl/myapps/common"
	"icikowski.pl/myapps/types"
)

// WebRepositoryLoader loads applications repository from given URL
type WebRepositoryLoader struct{}

// LoadRepository implements RepositoryLoader
func (*WebRepositoryLoader) LoadRepository(source string) (types.Repository, bool) {
	var repository types.Repository

	req, _ := http.NewRequest(http.MethodGet, source, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		if err == nil {
			err = fmt.Errorf("got status code %d", resp.StatusCode)
		}
		common.PrintErrorWhileMsg("loading source", source, err)
		return repository, false
	}

	if err := yaml.NewDecoder(resp.Body).Decode(&repository); err != nil {
		common.PrintErrorWhileMsg("parsing source", source, err)
		return repository, false
	}

	return repository, true
}

var _ RepositoryLoader = &WebRepositoryLoader{}

// WebRepositoryIntermediateLoader is a loader that acts as a proxy for
// WebRepositoryLoader by loading repositories with custom prefix
type WebRepositoryIntermediateLoader struct {
	prefix string
	target *WebRepositoryLoader
}

// LoadRepository implements RepositoryLoader
func (l *WebRepositoryIntermediateLoader) LoadRepository(source string) (types.Repository, bool) {
	return l.target.LoadRepository(l.prefix + ":" + source)
}

var _ RepositoryLoader = &WebRepositoryIntermediateLoader{}
