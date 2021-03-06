package cli

import (
	cliv2 "github.com/urfave/cli/v2"
)

var version = "unknown"
var gitCommit = "unknown"

// MyApps is the definition of the application and its commands
var MyApps *cliv2.App = &cliv2.App{
	Name:    "myapps",
	Usage:   "Universal command line tool for managing manually installed applications",
	Version: version,
	Metadata: map[string]interface{}{
		"Git commit": gitCommit,
	},
	EnableBashCompletion: true,
	HideHelpCommand:      true,
	Authors: []*cliv2.Author{
		{
			Name:  "Piotr Icikowski",
			Email: "https://github.com/Icikowski",
		},
	},
	UseShortOptionHandling: true,
	Before:                 basicChecks,
	Commands: []*cliv2.Command{
		{
			Name:            "search",
			Usage:           "Searches for applications in configured repo(s)",
			Before:          basicChecks,
			Action:          search,
			Flags:           searchFlags,
			HideHelpCommand: false,
		},
		{
			Name:            "list",
			Aliases:         []string{"ls"},
			Usage:           "Lists installed applications",
			Before:          basicChecks,
			Action:          list,
			HideHelpCommand: true,
		},
		{
			Name:            "install",
			Usage:           "Installs the application(s)",
			Before:          allChecks,
			Action:          install,
			BashComplete:    installCompletion,
			HideHelpCommand: true,
		},
		{
			Name:            "update",
			Usage:           "Updates the application(s)",
			Before:          allChecks,
			Action:          update,
			Flags:           updateFlags,
			HideHelpCommand: true,
			BashComplete:    updateCompletion,
		},
		{
			Name:            "uninstall",
			Usage:           "Uninstalls the application",
			Before:          allChecks,
			Action:          uninstall,
			Flags:           uninstallFlags,
			BashComplete:    uninstallCompletion,
			HideHelpCommand: true,
		},
		{
			Name:    "repo",
			Aliases: []string{"repos"},
			Usage:   "Manages the application repositories",
			Subcommands: []*cliv2.Command{
				{
					Name:            "add",
					Usage:           "Adds new repository from source",
					Description:     "Adds new repository from source\n\nSource is determined by the prefix:\n  'file:'    repository loaded from file\n             e.g.: file:/tmp/some-repo.yaml\n  'github:'  repository loaded from GitHub contents\n             format: <user>/<repository>[@<branch>]/file\n             e.g.: github:Icikowski/MyApps/repositories/tools.yaml\n  'web:'     repository loaded from URL\n             e.g.: web:https://example.com/repository.yaml\n\nIf no prefix is specified, then the file loader is being used.\nSources starting with 'http://' or 'https://' will be loaded using web loader.",
					Before:          allChecks,
					Action:          addRepos,
					Flags:           addReposFlags,
					HideHelpCommand: true,
				},
				{
					Name:            "list",
					Aliases:         []string{"ls"},
					Usage:           "Shows the list of available repositories",
					Before:          basicChecks,
					Action:          listRepos,
					HideHelpCommand: true,
				},
				{
					Name:            "show",
					Usage:           "Shows the contents of the repository",
					Before:          basicChecks,
					Action:          showRepo,
					BashComplete:    repositoryNameCompletion,
					HideHelpCommand: true,
				},
				{
					Name:            "remove",
					Aliases:         []string{"rm"},
					Usage:           "Removes the repository",
					Before:          allChecks,
					Action:          removeRepos,
					BashComplete:    repositoryNameCompletion,
					HideHelpCommand: true,
				},
				{
					Name:            "default",
					Usage:           "Gets/sets the default repository",
					Before:          basicChecks,
					Action:          defaultRepo,
					BashComplete:    repositoryNameCompletion,
					HideHelpCommand: true,
				},
			},
			HideHelpCommand: true,
		},
	},
}
