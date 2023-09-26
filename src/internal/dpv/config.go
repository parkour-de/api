package dpv

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	DB struct {
		Host string `yaml:"host"`
		Root string `yaml:"root"`
		Port int    `yaml:"port"`
		User string `yaml:"user"`
		Pass string `yaml:"pass"`
	} `yaml:"db"`
	Auth struct {
		DpvSecretKey     string `yaml:"dpv_secret_key"`
		DpvTokenSeconds  int    `yaml:"dpv_token_seconds"`
		FacebookGraphUrl string `yaml:"facebook_graph_url"`
		FacebookAppId    string `yaml:"facebook_app_id"`
	} `yaml:"auth"`
}

var ConfigInstance *Config

func NewConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("could not load config file: %w", err)
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, fmt.Errorf("could not decode config file: %w", err)
	}

	return config, nil
}
