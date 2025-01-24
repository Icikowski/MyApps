package repos

import "git.sr.ht/~icikowski/myapps/types"

// RepositoryLoader is used to load repositories from a given source
type RepositoryLoader interface {
	LoadRepository(source string) (types.Repository, bool)
}
