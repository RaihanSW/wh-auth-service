package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const fname = "application_%s.json"

func ReadConfig(env string) (*ApplicationConfig, error) {
	var cfg ApplicationConfig

	err := readJSON(fmt.Sprintf(fname, env), &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func readJSON(fname string, cfg any) error {
	data, err := os.ReadFile(fname)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(data), &cfg)
	if err != nil {
		return err
	}

	return nil
}
