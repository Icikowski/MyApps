package repos

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"icikowski.pl/myapps/common"
	"icikowski.pl/myapps/types"
)

// GitHubRepositoryLoader loads applications repository from GitHub repository
type GitHubRepositoryLoader struct {
	config types.GitHubRepositoryLoaderConfig
	target *WebRepositoryLoader
}

// LoadRepository implements RepositoryLoader
func (l *GitHubRepositoryLoader) LoadRepository(source string) (types.Repository, bool) {
	var repository types.Repository

	// Format: <user>/<repo>[@<branch]/<filename>
	re := regexp.MustCompile(`^(?P<user>[^/]+)/(?P<repo>[^@/]+)(@(?P<branch>[^/]+))?/(?P<filename>.+)$`)
	matches := re.FindStringSubmatch(source)

	if len(matches) == 0 {
		common.PrintErrorWhileMsg("parsing source", source, fmt.Errorf("invalid format"))
		return repository, false
	}

	user := matches[re.SubexpIndex("user")]
	repo := matches[re.SubexpIndex("repo")]
	branch := matches[re.SubexpIndex("branch")]
	filename := matches[re.SubexpIndex("filename")]

	if len(user) == 0 || len(repo) == 0 || len(filename) == 0 {
		common.PrintErrorWhileMsg("parsing source", source, fmt.Errorf("invalid format"))
		return repository, false
	}

	if len(branch) == 0 {
		if l.config.DetectDefaultBranch {
			url := fmt.Sprintf("https://api.github.com/repos/%s/%s", user, repo)

			req, _ := http.NewRequest(http.MethodGet, url, nil)
			req.Header.Set("Accept", "application/vnd.github.v3+json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil || resp.StatusCode != http.StatusOK {
				if err == nil {
					err = fmt.Errorf("got status code %d", resp.StatusCode)
				}
				common.PrintErrorWhileMsg("detecting default branch of", fmt.Sprintf("%s/%s", user, repo), err)
				return repository, false
			}

			var out map[string]any
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				common.PrintErrorWhileMsg("detecting default branch of", fmt.Sprintf("%s/%s", user, repo), err)
				return repository, false
			}

			default_branch := out["default_branch"]
			if default_branch == nil {
				common.PrintErrorWhileMsg("detecting default branch of", fmt.Sprintf("%s/%s", user, repo), fmt.Errorf("repository (probably) does not exist"))
				return repository, false
			}
			branch = out["default_branch"].(string)
		} else {
			branch = l.config.DefaultBranchName
		}
	}

	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", user, repo, branch, filename)
	return l.target.LoadRepository(url)
}

var _ RepositoryLoader = &GitHubRepositoryLoader{}
