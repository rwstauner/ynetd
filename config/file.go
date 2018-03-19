package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

func parseConfigFile(file string) (cfg Config, err error) {
	var content []byte
	content, err = ioutil.ReadFile(file)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		return
	}

	return
}
