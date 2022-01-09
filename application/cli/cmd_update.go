package cli

import (
	"fmt"
	"strings"
	"sync"
	"time"

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

	outdatedApps := []string{}
	conflictApps := map[string]string{}
	outdatedAppsMutex, conflictAppsMutex, waitGroup := sync.Mutex{}, sync.Mutex{}, sync.WaitGroup{}
	waitGroup.Add(len(effectiveApps))

	for _, appFullName := range effectiveApps {
		appFullName := appFullName

		go func() {
			defer waitGroup.Done()

			splittedAppFullName := strings.Split(appFullName, "/")
			repoName, appName := splittedAppFullName[0], splittedAppFullName[1]

			_, ok := deployments.Find(repoName, appName)
			if !ok {
				conflictAppsMutex.Lock()
				conflictApps[appFullName] = "application is not installed"
				conflictAppsMutex.Unlock()
				return
			}

			repo, ok := repos.FindByName(repoName)
			if !ok {
				conflictAppsMutex.Lock()
				conflictApps[appFullName] = "repository not found"
				conflictAppsMutex.Unlock()
				return
			}

			app, ok := repo.Contents.FindByName(appName)
			if !ok {
				conflictAppsMutex.Lock()
				conflictApps[appFullName] = "application not found"
				conflictAppsMutex.Unlock()
				return
			}

			updateAvailable, err := app.IsUpdateAvailable()
			if err != nil {
				conflictAppsMutex.Lock()
				conflictApps[appFullName] = fmt.Sprintf("cannot check update availability: %s", err.Error())
				conflictAppsMutex.Unlock()
				return
			}

			if updateAvailable {
				outdatedAppsMutex.Lock()
				outdatedApps = append(outdatedApps, appFullName)
				outdatedAppsMutex.Unlock()
			}
		}()
	}
	waitGroup.Wait()

	for appFullName, errorMsg := range conflictApps {
		fmt.Printf("%s: %s\n", common.FmtHeader(appFullName), errorMsg)
	}

	if len(outdatedApps) != 0 {
		fmt.Printf("Updating %s application(s)...\n", common.FmtHeader("%d", len(outdatedApps)))
		errorsOcurred := false
		for _, appFullName := range outdatedApps {
			splittedAppFullName := strings.Split(appFullName, "/")
			repoName, appName := splittedAppFullName[0], splittedAppFullName[1]

			repo, _ := repos.FindByName(repoName)
			app, _ := repo.Contents.FindByName(appName)

			if err := app.Update(); err != nil {
				common.BoxErr.Print(fmt.Sprintf(
					"%s cannot be updated due to error: %s",
					common.FmtHeader("%s", app),
					err.Error(),
				))
				errorsOcurred = true
				continue
			} else {
				common.BoxSuccess.Print(fmt.Sprintf(
					"%s has been updated successfully",
					common.FmtHeader("%s", app.Name),
				))
			}

			deployment, _ := deployments.Find(repoName, appName)
			deployment.UpdatedOn = time.Now()
			deployments.Update(deployment)
		}
		config.SetDeployments(deployments)

		if errorsOcurred {
			return common.ExitWithErrMsg("some applications were not updated")
		}
	} else {
		fmt.Printf("There are no applications that require updates.\n")
	}

	return nil
}
