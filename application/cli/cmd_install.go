package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	cliv2 "github.com/urfave/cli/v2"
	"icikowski.pl/myapps/config"
	"icikowski.pl/myapps/types"
)

func install(ctx *cliv2.Context) error {
	args := ctx.Args()
	if !args.Present() {
		return exitErrMsg("specify at least one application to install")
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
		if deployments.Exists(repoName, appName) {
			printBox(warningBox, "Application is already installed")
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

		printBox(infoBox, "Installing application")
		if err := app.Install(); err != nil {
			printBox(errorBox, fmt.Sprintf("Failed to install application: %s", err.Error()))
			errorsOcurred = true
			continue
		}

		deployments = deployments.Add(types.Deployment{
			Repository:  repoName,
			Application: appName,
			InstalledOn: time.Now(),
			UpdatedOn:   time.Now(),
		})
		printBox(successBox, "Application installed successfully")
	}
	config.SetDeployments(deployments)

	fmt.Println()
	if errorsOcurred {
		return exitWarnMsg("some applications were not installed")
	}
	return nil
}
