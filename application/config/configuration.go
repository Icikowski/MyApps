package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
	"icikowski.pl/myapps/common"
	"icikowski.pl/myapps/types"
)

var config = types.Configuration{
	DefaultRepository: "default",
	GitHubLoader: types.GitHubRepositoryLoaderConfig{
		DetectDefaultBranch: true,
		DefaultBranchName:   "master",
	},
}

func init() {
	contents, err := ioutil.ReadFile(common.PathConfigurationFile)
	if err != nil {
		common.PrintErrorWhileMsg("reading configuration file", common.PathConfigurationFile, err)
		os.Exit(1)
	}

	err = yaml.Unmarshal(contents, &config)
	if err != nil {
		common.PrintErrorWhileMsg("parsing configuration file", common.PathConfigurationFile, err)
		os.Exit(1)
	}
}

// GetConfiguration returns current configuration
func GetConfiguration() types.Configuration {
	return config
}

// SetConfiguration sets the configuration and writes it to file
func SetConfiguration(config types.Configuration) {
	configuration := config
	contents, err := yaml.Marshal(configuration)
	if err != nil {
		common.PrintErrorWhileMsg("encoding", "configuration", err)
	}

	err = ioutil.WriteFile(common.PathConfigurationFile, contents, 0644)
	if err != nil {
		common.PrintErrorWhileMsg("writing configuration file", common.PathConfigurationFile, err)
	}
}
