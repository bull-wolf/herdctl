package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const DefaultConfigFile = "herd.yaml"

type Service struct {
	Command   string            `yaml:"command"`
	Dir       string            `yaml:"dir"`
	Env       map[string]string `yaml:"env"`
	DependsOn []string          `yaml:"depends_on"`
}

type Config struct {
	Version  string             `yaml:"version"`
	Services map[string]Service `yaml:"services"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file %q not found; run 'herdctl init' to create one", path)
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	for name, svc := range c.Services {
		if svc.Command == "" {
			return fmt.Errorf("service %q is missing required field 'command'", name)
		}
		for _, dep := range svc.DependsOn {
			if _, ok := c.Services[dep]; !ok {
				return fmt.Errorf("service %q depends on unknown service %q", name, dep)
			}
		}
	}
	return nil
}
