package cli

import (
	"errors"
	"fmt"

	"github.com/rodaine/table"
	cliv2 "github.com/urfave/cli/v2"
	"icikowski.pl/myapps/config"
	"icikowski.pl/myapps/types"
)

var searchFlags = []cliv2.Flag{
	&cliv2.StringSliceFlag{
		Name:    "repo",
		Aliases: []string{"r"},
		Usage:   "repo to search in (may be specified more than once)",
	},
	&cliv2.BoolFlag{
		Name:    "all",
		Aliases: []string{"a"},
		Usage:   "search all applications in specified repo(s)",
	},
}

type applicationWrapper struct {
	application types.Application
	repoName    string
}

func search(ctx *cliv2.Context) error {
	allRepos := config.GetRepositories()

	var reposToSearchIn types.Repositories
	if ctx.IsSet("repo") {
		reposToSearchIn = types.Repositories{}

		errorOcurred := false
		for _, repoName := range ctx.StringSlice("repo") {
			repo, ok := allRepos.FindByName(repoName)
			if !ok {
				printErrorWhileMsg("loading repository", repoName, errors.New("no such repository decalred"))
				errorOcurred = true
				continue
			}

			reposToSearchIn = append(reposToSearchIn, repo)
		}

		if errorOcurred {
			return exitErrMsg("some repositories were not found")
		}
	} else {
		reposToSearchIn = allRepos
	}

	foundApps := []applicationWrapper{}
	for _, repo := range reposToSearchIn {
		if ctx.Bool("all") {
			for _, app := range repo.Contents {
				foundApps = append(foundApps, applicationWrapper{
					application: app,
					repoName:    repo.Name,
				})
			}
		} else {
			args := ctx.Args()

			if !args.Present() {
				return exitErrMsg("at least one application name must be specified")
			}

			for _, appName := range args.Slice() {
				if apps, ok := repo.Contents.FindByNameLike(appName); ok {
					for _, app := range apps {
						foundApps = append(foundApps, applicationWrapper{
							application: app,
							repoName:    repo.Name,
						})
					}
				}
			}
		}
	}

	fmt.Printf("Found %s match(es)\n", headerFormatter("%d", len(foundApps)))

	if len(foundApps) != 0 {
		fmt.Println()
		tbl := table.New("Application", "Version", "Description")
		tbl.
			WithHeaderFormatter(headerFormatter).
			WithFirstColumnFormatter(firstColumnFormatter)

		bar := progressBar.Start(len(foundApps))
		for _, app := range foundApps {
			bar.Increment()
			ver, err := app.application.GetLatestVersion()

			var displayVersion string
			if err != nil {
				displayVersion = err.Error()
			} else {
				displayVersion = ver.String()
			}

			tbl.AddRow(
				fmt.Sprintf("%s/%s", app.repoName, app.application.Name),
				displayVersion, app.application.Description,
			)
		}
		finishProgressBar(bar)
		tbl.Print()
	}
	return nil
}
