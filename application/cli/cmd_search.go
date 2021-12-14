package cli

import (
	"errors"
	"fmt"
	"sync"

	"github.com/rodaine/table"
	cliv2 "github.com/urfave/cli/v2"
	"icikowski.pl/myapps/common"
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
				common.PrintErrorWhileMsg("loading repository", repoName, errors.New("no such repository decalred"))
				errorOcurred = true
				continue
			}

			reposToSearchIn = append(reposToSearchIn, repo)
		}

		if errorOcurred {
			return common.ExitWithErrMsg("some repositories were not found")
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
				return common.ExitWithErrMsg("at least one application name must be specified")
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

	fmt.Printf("Found %s match(es)\n", common.FmtHeader("%d", len(foundApps)))

	if len(foundApps) != 0 {
		fmt.Println()
		tbl := table.New("Application", "Version", "Description")
		tbl.
			WithHeaderFormatter(common.FmtHeader).
			WithFirstColumnFormatter(common.FmtFirstCol)

		tableLock, waitGroup := sync.Mutex{}, sync.WaitGroup{}

		bar := common.NewProgressBar(len(foundApps))
		waitGroup.Add(len(foundApps))
		for _, app := range foundApps {
			app := app

			go func(app applicationWrapper) {
				ver, err := app.application.GetLatestVersion()

				var displayVersion string
				if err != nil {
					displayVersion = err.Error()
				} else {
					displayVersion = ver.String()
				}

				tableLock.Lock()
				tbl.AddRow(
					fmt.Sprintf("%s/%s", app.repoName, app.application.Name),
					displayVersion, app.application.Description,
				)
				tableLock.Unlock()

				bar.Increment()
				waitGroup.Done()
			}(app)
		}
		waitGroup.Wait()
		bar.Finish()
		tbl.Print()
	}
	return nil
}
