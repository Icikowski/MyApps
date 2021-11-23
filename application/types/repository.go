package types

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/hashicorp/go-version"
	"github.com/rodaine/table"
)

// Repository represents the applications repository.
type Repository struct {
	Name        string       `json:"name" yaml:"name"`
	Description string       `json:"description" yaml:"description"`
	Maintainer  string       `json:"maintainer" yaml:"maintainer"`
	Contents    Applications `json:"contents" yaml:"contents"`
}

// Update retrieves the newest version of each application in the repository.
func (repo Repository) NewestVersions() (map[string]*version.Version, map[string]error) {
	states := map[string]*version.Version{}
	errors := map[string]error{}

	for _, app := range repo.Contents {
		version, err := app.GetLatestVersion()
		if err != nil {
			errors[app.Name] = err
			continue
		}

		states[app.Name] = version
	}

	if len(errors) > 0 {
		return states, errors
	}
	return states, nil
}

// Print writes pretty printed repository details to stdout.
func (repo Repository) Print() {
	fmt.Printf("%s %s\n%s %s\n%s %s\n%s %d\n",
		headerFormatter("Repository name:"), color.BlueString(repo.Name),
		headerFormatter("Description:    "), repo.Description,
		headerFormatter("Maintainer:     "), repo.Maintainer,
		headerFormatter("Apps count:     "), len(repo.Contents),
	)

	if len(repo.Contents) > 0 {
		fmt.Println("\n")

		tbl := table.New("Application name", "Description")
		tbl.
			WithHeaderFormatter(headerFormatter).
			WithFirstColumnFormatter(firstColumnFormatter)

		for _, app := range repo.Contents {
			tbl.AddRow(app.Name, app.Description)
		}

		tbl.Print()
	}
}

// Repositories is a list of repositories.
type Repositories []Repository

// FindByName returns the repository with the given name.
func (repos Repositories) FindByName(name string) (Repository, bool) {
	for _, repo := range repos {
		if repo.Name == name {
			return repo, true
		}
	}
	return Repository{}, false
}

// Print writes pretty printed repositories list to stdout.
func (repos Repositories) Print(defaultRepositoryName string) {
	fmt.Printf("%s %d\nDefault repository is marked with %s sign.\n\n",
		headerFormatter("Number of repositories:"), len(repos),
		firstColumnFormatter("*"),
	)

	tbl := table.New("", "Name", "Description", "Maintainer", "Apps count")
	tbl.
		WithHeaderFormatter(headerFormatter).
		WithFirstColumnFormatter(firstColumnFormatter)

	sign := map[bool]string{true: "*", false: ""}

	for _, repo := range repos {
		tbl.AddRow(
			sign[repo.Name == defaultRepositoryName],
			repo.Name, repo.Description,
			repo.Maintainer, len(repo.Contents),
		)
	}

	tbl.Print()
}
