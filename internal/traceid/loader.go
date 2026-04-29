package traceid

import (
	"fmt"

	"github.com/example/herdctl/internal/config"
)

// LoadFromConfig generates an initial trace ID for every service defined
// in the config and returns the populated store.
func LoadFromConfig(cfg *config.Config) (*Store, error) {
	s := New()
	for name := range cfg.Services {
		if _, err := s.Generate(name); err != nil {
			return nil, fmt.Errorf("traceid: failed to generate for %q: %w", name, err)
		}
	}
	return s, nil
}
