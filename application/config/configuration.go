package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
	"icikowski.pl/myapps/common"
	"icikowski.pl/myapps/types"
)

var config = types.Configuration{
	DefaultRepository: "default",
}

func init() {
	contents, err := ioutil.ReadFile(common.PathConfigurationFile)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(contents, &config)
	if err != nil {
		panic(err)
	}
}

func GetConfiguration() types.Configuration {
	return config
}

func SetConfiguration(config types.Configuration) {
	configuration := config
	contents, err := yaml.Marshal(configuration)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(common.PathConfigurationFile, contents, 0644)
	if err != nil {
		panic(err)
	}
}
