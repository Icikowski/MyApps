package repos

import (
	"regexp"

	"icikowski.pl/myapps/config"
	"icikowski.pl/myapps/types"
)

// RepositoryProcessor manages the repository parsers for different sources
type RepositoryProcessor struct {
	loaders           map[string]RepositoryLoader
	fallbackProcessor string
}

// Load loads the repository using dedicated loader
func (p *RepositoryProcessor) Load(prefixedSource string) (types.Repository, bool) {
	re := regexp.MustCompile(`^((?P<loader>[^:]*?):)?(?P<source>.+)$`)
	match := re.FindStringSubmatch(prefixedSource)
	loader := match[re.SubexpIndex("loader")]
	source := match[re.SubexpIndex("source")]

	processor, ok := p.loaders[loader]
	if !ok {
		processor = p.loaders[p.fallbackProcessor]
		source = prefixedSource
	}

	return processor.LoadRepository(source)
}

// NewRepositoryProcessor creates a new repository processor
func NewRepositoryProcessor() *RepositoryProcessor {
	webRepositoryLoader := &WebRepositoryLoader{}
	httpWebRepositoryIntermediateLoader := &WebRepositoryIntermediateLoader{
		prefix: "http",
		target: webRepositoryLoader,
	}
	httpsWebRepositoryIntermediateLoader := &WebRepositoryIntermediateLoader{
		prefix: "https",
		target: webRepositoryLoader,
	}

	loaders := map[string]RepositoryLoader{
		"file":  &FileRepositoryLoader{},
		"web":   webRepositoryLoader,
		"http":  httpWebRepositoryIntermediateLoader,
		"https": httpsWebRepositoryIntermediateLoader,
		"github": &GitHubRepositoryLoader{
			config: config.GetConfiguration().GitHubLoader,
			target: webRepositoryLoader,
		},
	}

	return &RepositoryProcessor{
		loaders:           loaders,
		fallbackProcessor: "file",
	}
}
