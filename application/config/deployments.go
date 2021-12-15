package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
	"icikowski.pl/myapps/common"
	"icikowski.pl/myapps/types"
)

var deployments = types.Deployments{}

func init() {
	contents, err := ioutil.ReadFile(common.PathDeploymentsFile)
	if err != nil {
		common.PrintErrorWhileMsg("reading deployments file", common.PathDeploymentsFile, err)
		os.Exit(1)
	}

	err = yaml.Unmarshal(contents, &deployments)
	if err != nil {
		common.PrintErrorWhileMsg("parsing deployments file", common.PathDeploymentsFile, err)
		os.Exit(1)
	}
}

// GetDeployments returns the list of currently deployed applications
func GetDeployments() types.Deployments {
	return deployments
}

// SetDeployments sets the list of currently deployed applications and writes it to file
func SetDeployments(d types.Deployments) {
	deployments = d
	contents, err := yaml.Marshal(deployments)
	if err != nil {
		common.PrintErrorWhileMsg("encoding data for", "deployments", err)
		os.Exit(1)
	}

	err = ioutil.WriteFile(common.PathDeploymentsFile, contents, 0644)
	if err != nil {
		common.PrintErrorWhileMsg("writing deployments file", common.PathDeploymentsFile, err)
	}
}
