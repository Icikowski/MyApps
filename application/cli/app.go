package cli

import (
	cliv2 "github.com/urfave/cli/v2"
)

var MyApps *cliv2.App = &cliv2.App{
	Name:                 "myapps",
	Usage:                "Universal command line tool for managing manually installed applications",
	Version:              "1.0.0",
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
			Usage:           "Searchs for applications in configured repo(s)",
			Before:          basicChecks,
			Action:          search,
			Flags:           searchFlags,
			HideHelpCommand: true,
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
			HideHelpCommand: true,
		},
		{
			Name:            "update",
			Usage:           "Updates the application(s)",
			Before:          allChecks,
			Action:          update,
			Flags:           updateFlags,
			HideHelpCommand: true,
		},
		{
			Name:            "uninstall",
			Usage:           "Uninstalls the application",
			Before:          allChecks,
			Action:          uninstall,
			Flags:           uninstallFlags,
			HideHelpCommand: true,
		},
		{
			Name:  "repos",
			Usage: "Manages the application repositories",
			Subcommands: []*cliv2.Command{
				{
					Name:            "add",
					Usage:           "Adds new repository from file",
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
					HideHelpCommand: true,
				},
				{
					Name:            "remove",
					Aliases:         []string{"rm"},
					Usage:           "Removes the repository",
					Before:          allChecks,
					Action:          removeRepos,
					HideHelpCommand: true,
				},
			},
			HideHelpCommand: true,
		},
	},
}
