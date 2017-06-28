package cfg

import (
	"os"

	"github.com/naoina/toml"
)

type Config struct {
	Project []ProjectConf
}

type ProjectConf struct {
	Name     string
	Root     string
	Dirs     []string
	Exclude  []string
	CtagArgs string
	Tags     string
}

func ReadConfig(files ...string) (*Config, error) {
	for _, file := range files {
		f, err := os.Open(file)
		if os.IsNotExist(err) {
			continue
		} else if err != nil {
			return nil, err
		}
		var config Config
		err = toml.NewDecoder(f).Decode(&config)
		if err != nil {
			return nil, err
		}
		return &config, nil
	}
	return nil, os.ErrNotExist
}
