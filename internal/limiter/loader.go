package limiter

import (
	"fmt"

	"github.com/user/herdctl/internal/config"
)

// LoadFromConfig reads concurrency limits from the herd config and registers
// them with the provided Limiter. Services without a concurrency value are
// skipped.
func LoadFromConfig(l *Limiter, cfg *config.Config) error {
	for name, svc := range cfg.Services {
		if svc.Concurrency <= 0 {
			continue
		}
		if err := l.Set(name, svc.Concurrency); err != nil {
			return fmt.Errorf("limiter: service %q: %w", name, err)
		}
	}
	return nil
}
