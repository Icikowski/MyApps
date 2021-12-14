package cli

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	cliv2 "github.com/urfave/cli/v2"
	"icikowski.pl/myapps/common"
	"icikowski.pl/myapps/config"
)

var uninstallFlags = []cliv2.Flag{}

func uninstall(ctx *cliv2.Context) error {
	args := ctx.Args()
	if !args.Present() {
		return common.ExitWithErrMsg("specify at least one application to uninstall")
	}

	defaultRepo := config.GetConfiguration().DefaultRepository
	effectiveApps := []string{}
	for _, appName := range args.Slice() {
		if !strings.Contains(appName, "/") {
			appName = defaultRepo + "/" + appName
		}
		effectiveApps = append(effectiveApps, appName)
	}

	repos := config.GetRepositories()
	deployments := config.GetDeployments()

	fmt.Printf("Processing %s application(s)...\n", common.FmtHeader("%d", len(effectiveApps)))

	errorsOcurred := false
	for _, appFullName := range effectiveApps {
		fmt.Printf("\nProcessing %s...\n", color.BlueString(appFullName))
		splittedAppFullName := strings.Split(appFullName, "/")
		repoName, appName := splittedAppFullName[0], splittedAppFullName[1]

		common.BoxInfo.Print("Checking for deployment presence")
		if !deployments.Exists(repoName, appName) {
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

		common.BoxInfo.Print("Uninstalling application")
		if err := app.Uninstall(); err != nil {
			common.BoxErr.Print(fmt.Sprintf("Failed to uninstall application: %s", err.Error()))
			errorsOcurred = true
			continue
		}

		deployments = deployments.Delete(repoName, appName)
		common.BoxSuccess.Print("Application uninstalled successfully")
	}
	config.SetDeployments(deployments)

	fmt.Println()
	if errorsOcurred {
		return common.ExitWithErrMsg("some applications were not uninstalled")
	}
	return nil
}
