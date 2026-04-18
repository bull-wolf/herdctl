package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Service defines a single managed process.
type Service struct {
	Command  string   `yaml:"command"`
	Dir      string   `yaml:"dir"`
	Depends  []string `yaml:"depends_on"`
}

// Config is the top-level structure of herd.yaml.
type Config struct {
	Services map[string]Service `yaml:"services"`
}

// Load reads and validates a herd.yaml file at the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if err := validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func validate(cfg *Config) error {
	for name, svc := range cfg.Services {
		if svc.Command == "" {
			return fmt.Errorf("service %q is missing a command", name)
		}
		for _, dep := range svc.Depends {
			if _, ok := cfg.Services[dep]; !ok {
				return fmt.Errorf("service %q depends on unknown service %q", name, dep)
			}
		}
	}
	return nil
}
