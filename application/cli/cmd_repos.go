package cli

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
	cliv2 "github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	"icikowski.pl/myapps/common"
	"icikowski.pl/myapps/config"
	"icikowski.pl/myapps/repos"
)

var addReposFlags = []cliv2.Flag{
	&cliv2.BoolFlag{
		Name:    "force",
		Aliases: []string{"f"},
		Usage:   "force add repository if it already exists",
	},
}

func repositoryNameCompletion(ctx *cliv2.Context) {
	for _, repo := range config.GetRepositories() {
		fmt.Println(repo.Name)
	}
}

func addRepos(ctx *cliv2.Context) error {
	args := ctx.Args()
	if !args.Present() {
		return common.ExitWithErrMsg("at least one source must be specified")
	}

	processor := repos.NewRepositoryProcessor()

	fmt.Println("Processing", color.New(color.FgHiWhite, color.Bold).Sprint(args.Len()), "source(s)...")

	errOcurred := false
	for _, source := range args.Slice() {
		repo, ok := processor.Load(source)
		if !ok {
			errOcurred = true
			continue
		}

		target := fmt.Sprintf("%s/%s.yaml", common.PathRepositories, repo.Name)

		if _, err := os.Stat(target); err == nil && !ctx.Bool("force") {
			common.PrintErrorWhileMsg("adding repository", repo.Name, errors.New("repository already exists"))
			errOcurred = true
			continue
		}

		contents, err := yaml.Marshal(repo)
		if err != nil {
			common.PrintErrorWhileMsg("serializing repository", repo.Name, err)
			errOcurred = true
			continue
		}

		if err := ioutil.WriteFile(target, contents, 0644); err != nil {
			common.PrintErrorWhileMsg("storing repository", repo.Name, err)
			errOcurred = true
			continue
		}

		common.PrintSuccessfullyMsg("added repository", repo.Name)
	}

	if errOcurred {
		return common.ExitWithWarnMsg("some repositories were not stored due to errors")
	}
	return nil
}

func listRepos(ctx *cliv2.Context) error {
	config.GetRepositories().Print(config.GetConfiguration().DefaultRepository)
	return nil
}

func showRepo(ctx *cliv2.Context) error {
	if ctx.Args().Len() != 1 {
		return common.ExitWithErrMsg("exactly one repository name must be specified")
	}

	repoName := ctx.Args().First()
	repo, ok := config.GetRepositories().FindByName(repoName)
	if !ok {
		return common.ExitWithErrMsg(fmt.Sprint("repository ", color.BlueString(repoName), " not found"))
	}

	repo.Print()
	return nil
}

func removeRepos(ctx *cliv2.Context) error {
	args := ctx.Args()
	if !args.Present() {
		return common.ExitWithErrMsg("at least one repository must be specified")
	}

	fmt.Println("Processing", color.New(color.FgHiWhite, color.Bold).Sprint(args.Len()), "repo(s)...")

	deployments := config.GetDeployments()

	usedRepos := map[string]struct{}{}
	for _, deployment := range deployments {
		if _, ok := usedRepos[deployment.Repository]; !ok {
			usedRepos[deployment.Repository] = struct{}{}
		}
	}

	repos := config.GetRepositories()
	errOcurred := false
	for _, repoName := range args.Slice() {
		repo, ok := repos.FindByName(repoName)
		if !ok {
			common.PrintErrorWhileMsg("removing repository", repoName, errors.New("repository not found"))
			errOcurred = true
			continue
		}

		if repo.Name == config.GetConfiguration().DefaultRepository {
			common.PrintErrorWhileMsg("removing repository", repoName, errors.New("cannot remove default repository"))
			errOcurred = true
			continue
		}

		if _, ok := usedRepos[repo.Name]; ok {
			common.PrintErrorWhileMsg("removing repository", repoName, errors.New("cannot remove repository because some applications from it are installed"))
			errOcurred = true
			continue
		}

		target := fmt.Sprintf("%s/%s.yaml", common.PathRepositories, repoName)

		if err := os.Remove(target); err != nil {
			common.PrintErrorWhileMsg("removing repository", repoName, err)
			errOcurred = true
			continue
		}

		common.PrintSuccessfullyMsg("removed repository", repoName)
	}

	if errOcurred {
		return common.ExitWithErrMsg("some repositories were not removed due to errors")
	}
	return nil
}

func defaultRepo(ctx *cliv2.Context) error {
	if ctx.Args().Len() != 1 {
		defaultRepo := config.GetConfiguration().DefaultRepository
		fmt.Println("Default repository name is", color.BlueString(defaultRepo))
		return nil
	}

	if err := allChecks(ctx); err != nil {
		return err
	}

	repoName := ctx.Args().First()
	repo, ok := config.GetRepositories().FindByName(repoName)
	if !ok {
		return common.ExitWithErrMsg(fmt.Sprint("repository ", color.BlueString(repoName), " not found"))
	}

	configuration := config.GetConfiguration()
	configuration.DefaultRepository = repo.Name
	config.SetConfiguration(configuration)

	common.PrintSuccessfullyMsg("set default repository to", repo.Name)
	return nil
}
