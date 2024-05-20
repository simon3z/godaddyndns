package cmd

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Key    string `yaml:"key"`
	Secret string `yaml:"secret"`
	Domain string `yaml:"domain"`
	Host   string `yaml:"host"`
}

func LoadConfiguration(path string) (*Config, error) {
	f, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	config := new(Config)

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(config)

	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) FullDomain() string {
	return c.Host + "." + c.Domain
}
