package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	ErrPathEmpty   = errors.New("Path is empty")
	ErrCantRead    = errors.New("Can't read")
	ErrUnmarsharll = errors.New("Unmarshal error!")
)

type Config struct {
	LogLevel string `yaml:"LogLevel"`
	ServerIp string `yaml:"ServerIp"`
	DataBase string `yaml:"DataBase"`
}

func LoadConfig(path string) (*Config, error) {

	if len(path) == 0 {
		return nil, fmt.Errorf("%w %v", ErrPathEmpty, path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w %v", ErrCantRead, path)
	}

	var config Config

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrUnmarsharll)
	}

	return &config, nil
}
