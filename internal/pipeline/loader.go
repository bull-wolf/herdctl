package pipeline

import (
	"fmt"

	"github.com/example/herdctl/internal/config"
)

// LoadFromConfig registers pipelines for every service that declares pipeline
// stages in the herd.yaml configuration.
func LoadFromConfig(p *Pipeline, cfg *config.Config) error {
	for _, svc := range cfg.Services {
		if len(svc.Pipeline) == 0 {
			continue
		}
		if err := p.Register(svc.Name, svc.Pipeline); err != nil {
			return fmt.Errorf("pipeline: registering %q: %w", svc.Name, err)
		}
	}
	return nil
}
