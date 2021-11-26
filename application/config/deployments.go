package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
	"icikowski.pl/myapps/common"
	"icikowski.pl/myapps/types"
)

var deployments = types.Deployments{}

func init() {
	contents, err := ioutil.ReadFile(common.PathDeploymentsFile)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(contents, &deployments)
	if err != nil {
		panic(err)
	}
}

func GetDeployments() types.Deployments {
	return deployments
}

func SetDeployments(d types.Deployments) {
	deployments = d
	contents, err := yaml.Marshal(deployments)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(common.PathDeploymentsFile, contents, 0644)
	if err != nil {
		panic(err)
	}
}
