package cli

import (
	"fmt"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	cliv2 "github.com/urfave/cli/v2"
	"icikowski.pl/myapps/common"
	"icikowski.pl/myapps/config"
	"icikowski.pl/myapps/types"
)

const (
	signUpdate = "+"
	signError  = "!"
)

var conditionalFirstColumnFormatter = func(s string, i ...interface{}) string {
	switch string(i[0].(string)[0]) {
	case signUpdate:
		return common.FmtFirstCol(s, i...)
	case signError:
		return color.New(color.FgHiRed, color.Bold).Sprintf(s, i...)
	default:
		return fmt.Sprintf(s, i...)
	}
}

func list(ctx *cliv2.Context) error {
	deployments := config.GetDeployments()
	repos := config.GetRepositories()

	fmt.Printf("Found %s installed application(s)\n", common.FmtHeader("%d", len(deployments)))

	if len(deployments) == 0 {
		return nil
	}

	fmt.Println()

	tbl := table.New("", "Application", "Current version", "Latest version", "Details")
	tbl.
		WithHeaderFormatter(common.FmtHeader).
		WithFirstColumnFormatter(conditionalFirstColumnFormatter)

	tableLock, waitGroup := sync.Mutex{}, sync.WaitGroup{}

	bar := common.NewProgressBar(len(deployments))
	for _, deployment := range deployments {
		waitGroup.Add(1)
		deployment := deployment

		go func(deployment types.Deployment) {
			sign, currentVersion, latestVersion, details := signError, "N/A", "N/A", "-"

			if repo, ok := repos.FindByName(deployment.Repository); ok {
				if app, ok := repo.Contents.FindByName(deployment.Application); ok {
					errs := []error{}

					currentVersionObj, err := app.GetCurrentVersion()
					if err != nil {
						errs = append(errs, err)
					} else {
						currentVersion = currentVersionObj.String()
					}

					latestVersionObj, err := app.GetLatestVersion()
					if err != nil {
						errs = append(errs, err)
					} else {
						latestVersion = latestVersionObj.String()
					}

					if len(errs) != 0 {
						strErrs := []string{}
						for _, err := range errs {
							strErrs = append(strErrs, err.Error())
						}

						details = fmt.Sprintf("Cannot proceed with version check: %s",
							strings.Join(strErrs, "; "),
						)
					} else if currentVersionObj.LessThan(latestVersionObj) {
						sign = signUpdate
						details = "Update available"
					} else {
						sign = ""
					}
				} else {
					details = "Application not found in repository"
				}
			} else {
				details = "Repository not found"
			}

			tableLock.Lock()
			tbl.AddRow(
				sign, deployment.String(),
				currentVersion, latestVersion, details,
			)
			tableLock.Unlock()

			bar.Increment()
			waitGroup.Done()
		}(deployment)
	}
	waitGroup.Wait()
	bar.Finish()
	tbl.Print()
	return nil
}
