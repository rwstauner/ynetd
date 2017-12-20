package config

import (
	"encoding/json"
	"os"
)

func parseConfigFile(file string) (cfg Config, err error) {
	var fh *os.File
	fh, err = os.Open(file)
	if err != nil {
		return
	}
	err = json.NewDecoder(fh).Decode(&cfg)
	return
}
