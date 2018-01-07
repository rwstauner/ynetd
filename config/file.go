package config

import (
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type jsonService struct {
	Proxy          map[string]string `yaml:"Proxy"`
	Command        []string          `yaml:"Command"`
	StopAfter      string            `yaml:"StopAfter"`
	StopSignal     string            `yaml:"StopSignal"`
	Timeout        string            `yaml:"Timeout"`
	WaitAfterStart string            `yaml:"WaitAfterStart"`
}

type backcompatConfig struct {
	Config       `yaml:",inline"`
	JSONServices []jsonService `yaml:"Services"`
}

func parseConfigFile(file string) (cfg Config, err error) {
	var content []byte
	content, err = ioutil.ReadFile(file)
	if err != nil {
		return
	}

	bcfg := backcompatConfig{}
	err = yaml.Unmarshal(content, &bcfg)
	if err != nil {
		return
	}

	cfg.Services = bcfg.Services

	if len(bcfg.JSONServices) > 0 {
		fmt.Fprintln(os.Stderr, "The JSON spellings are deprecated.  Please use the yaml lowercase/underscore spellings.")
		for _, js := range bcfg.JSONServices {
			cfg.Services = append(cfg.Services, Service(js))
		}
	}

	return
}
