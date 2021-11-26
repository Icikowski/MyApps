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
	"icikowski.pl/myapps/types"
)

var addReposFlags = []cliv2.Flag{
	&cliv2.BoolFlag{
		Name:    "force",
		Aliases: []string{"f"},
		Usage:   "force add repository if it already exists",
	},
}

func addRepos(ctx *cliv2.Context) error {
	args := ctx.Args()
	if !args.Present() {
		return exitErrMsg("at least one input file must be specified")
	}

	fmt.Println("Processing", color.New(color.FgHiWhite, color.Bold).Sprint(args.Len()), "file(s)...")

	errOcurred := false
	for _, filename := range args.Slice() {
		contents, err := ioutil.ReadFile(filename)
		if err != nil {
			printErrorWhileMsg("reading file", filename, err)
			errOcurred = true
			continue
		}

		var repo types.Repository
		if err := yaml.Unmarshal(contents, &repo); err != nil {
			printErrorWhileMsg("parsing file", filename, err)
			errOcurred = true
			continue
		}

		target := fmt.Sprintf("%s/%s.yaml", common.PathRepositories, repo.Name)

		if _, err := os.Stat(target); err == nil && !ctx.Bool("force") {
			printErrorWhileMsg("adding repository", repo.Name, errors.New("repository already exists"))
			errOcurred = true
			continue
		}

		if err := ioutil.WriteFile(target, contents, 0644); err != nil {
			printErrorWhileMsg("storing repository", repo.Name, err)
			errOcurred = true
			continue
		}

		printSuccessfully("added repository", repo.Name)
	}

	if errOcurred {
		return exitWarnMsg("some repositories were not stored due to errors")
	}
	return nil
}

func listRepos(ctx *cliv2.Context) error {
	config.GetRepositories().Print(config.GetConfiguration().DefaultRepository)
	return nil
}

func showRepo(ctx *cliv2.Context) error {
	if ctx.Args().Len() != 1 {
		return exitErrMsg("exactly one repository name must be specified")
	}

	repoName := ctx.Args().First()
	repo, ok := config.GetRepositories().FindByName(repoName)
	if !ok {
		return exitErrMsg(fmt.Sprint("repository", color.BlueString(repoName), "not found"))
	}

	repo.Print()
	return nil
}

func removeRepos(ctx *cliv2.Context) error {
	args := ctx.Args()
	if !args.Present() {
		return exitErrMsg("at least one repository must be specified")
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
			printErrorWhileMsg("removing repository", repoName, errors.New("repository not found"))
			errOcurred = true
			continue
		}

		if repo.Name == config.GetConfiguration().DefaultRepository {
			printErrorWhileMsg("removing repository", repoName, errors.New("cannot remove default repository"))
			errOcurred = true
			continue
		}

		if _, ok := usedRepos[repo.Name]; ok {
			printErrorWhileMsg("removing repository", repoName, errors.New("cannot remove repository because some applications from it are installed"))
			errOcurred = true
			continue
		}

		target := fmt.Sprintf("%s/%s.yaml", common.PathRepositories, repoName)

		if err := os.Remove(target); err != nil {
			printErrorWhileMsg("removing repository", repoName, err)
			errOcurred = true
			continue
		}

		printSuccessfully("removed repository", repoName)
	}

	if errOcurred {
		return exitWarnMsg("some repositories were not removed due to errors")
	}
	return nil
}
