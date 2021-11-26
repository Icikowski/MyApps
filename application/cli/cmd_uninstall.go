package cli

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	cliv2 "github.com/urfave/cli/v2"
	"icikowski.pl/myapps/config"
)

var uninstallFlags = []cliv2.Flag{}

func uninstall(ctx *cliv2.Context) error {
	args := ctx.Args()
	if !args.Present() {
		return exitErrMsg("specify at least one application to uninstall")
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

	fmt.Printf("Processing %s application(s)...\n", headerFormatter("%d", len(effectiveApps)))

	errorsOcurred := false
	for _, appFullName := range effectiveApps {
		fmt.Printf("\nProcessing %s...\n", color.BlueString(appFullName))
		splittedAppFullName := strings.Split(appFullName, "/")
		repoName, appName := splittedAppFullName[0], splittedAppFullName[1]

		printBox(infoBox, "Checking for deployment presence")
		if !deployments.Exists(repoName, appName) {
			printBox(warningBox, "Application is not installed")
			continue
		}

		printBox(infoBox, "Checking for repository presence")
		repo, ok := repos.FindByName(repoName)
		if !ok {
			printBox(errorBox, "Repository not found")
			errorsOcurred = true
			continue
		}

		printBox(infoBox, "Checking for application presence")
		app, ok := repo.Contents.FindByName(appName)
		if !ok {
			printBox(errorBox, "Application not found")
			errorsOcurred = true
			continue
		}

		printBox(infoBox, "Uninstalling application")
		if err := app.Uninstall(); err != nil {
			printBox(errorBox, fmt.Sprintf("Failed to uninstall application: %s", err.Error()))
			errorsOcurred = true
			continue
		}

		deployments = deployments.Delete(repoName, appName)
		printBox(successBox, "Application uninstalled successfully")
	}
	config.SetDeployments(deployments)

	fmt.Println()
	if errorsOcurred {
		return exitWarnMsg("some applications were not uninstalled")
	}
	return nil
}
