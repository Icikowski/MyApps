package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	cliv2 "github.com/urfave/cli/v2"
	"icikowski.pl/myapps/common"
	"icikowski.pl/myapps/config"
)

var updateFlags = []cliv2.Flag{
	&cliv2.BoolFlag{
		Name:    "all",
		Aliases: []string{"a"},
		Usage:   "update all installed applications",
	},
}

func update(ctx *cliv2.Context) error {
	repos := config.GetRepositories()
	deployments := config.GetDeployments()

	effectiveApps := []string{}
	if ctx.Bool("all") {
		for _, deployment := range deployments {
			effectiveApps = append(effectiveApps, deployment.String())
		}
	} else {
		args := ctx.Args()
		if !args.Present() {
			return common.ExitWithErrMsg("specify at least one application to update")
		}

		for _, appName := range args.Slice() {
			if !strings.Contains(appName, "/") {
				appName = config.GetConfiguration().DefaultRepository + "/" + appName
			}
			effectiveApps = append(effectiveApps, appName)
		}
	}

	fmt.Printf("Processing %s application(s)...\n", common.FmtHeader("%d", len(effectiveApps)))

	errorsOcurred := false
	for _, appFullName := range effectiveApps {
		fmt.Printf("\nProcessing %s...\n", color.BlueString(appFullName))
		splittedAppFullName := strings.Split(appFullName, "/")
		repoName, appName := splittedAppFullName[0], splittedAppFullName[1]

		common.BoxInfo.Print("Checking for deployment presence")
		deployment, ok := deployments.Find(repoName, appName)
		if !ok {
			common.BoxWarn.Print("Application is not installed")
			continue
		}

		common.BoxInfo.Print("Checking for repository presence")
		repo, ok := repos.FindByName(repoName)
		if !ok {
			common.BoxErr.Print("Repository not found")
			errorsOcurred = true
			continue
		}

		common.BoxInfo.Print("Checking for application presence")
		app, ok := repo.Contents.FindByName(appName)
		if !ok {
			common.BoxErr.Print("Application not found")
			errorsOcurred = true
			continue
		}

		common.BoxInfo.Print("Checking for updates")
		currentVersion, err := app.GetCurrentVersion()
		if err != nil {
			common.BoxErr.Print(fmt.Sprintf("Cannot check current version: %s", err.Error()))
			errorsOcurred = true
			continue
		}
		latestVersion, err := app.GetLatestVersion()
		if err != nil {
			common.BoxErr.Print(fmt.Sprintf("Cannot check latest version: %s", err.Error()))
			errorsOcurred = true
			continue
		}
		if currentVersion.LessThan(latestVersion) {
			common.BoxErr.Print("Newer version is available, updating application")
			if err := app.Update(); err != nil {
				common.BoxErr.Print(fmt.Sprintf("Failed to update application: %s", err.Error()))
				errorsOcurred = true
				continue
			}

			deployment.UpdatedOn = time.Now()
			deployments = deployments.Update(deployment)
			common.BoxSuccess.Print("Application updated successfully")
		} else {
			common.BoxSuccess.Print("Application is already up to date")
		}
	}
	config.SetDeployments(deployments)

	if errorsOcurred {
		return common.ExitWithWarnMsg("some applications failed to update")
	}
	return nil
}
