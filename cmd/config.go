package cmd

// cspell:ignore godaddy nsdyndns

import (
	"os"

	"github.com/simon3z/nsdyndns"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Domain  string `yaml:"domain"`
	Host    string `yaml:"host"`
	GoDaddy struct {
		Key    string `yaml:"key"`
		Secret string `yaml:"secret"`
	}
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

type ConfigNameService struct {
	Name    string
	Service nsdyndns.NameService
}

func (c *Config) GetNameServices() []*ConfigNameService {
	return []*ConfigNameService{
		{"godaddy", nsdyndns.NewGoDaddyService(c.GoDaddy.Key, c.GoDaddy.Secret)},
	}
}
