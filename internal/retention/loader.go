package retention

import (
	"time"

	"github.com/user/herdctl/internal/config"
)

// LoadFromConfig reads retention settings from the herd config and registers
// them with the provided Manager.
func LoadFromConfig(m *Manager, cfg *config.Config) error {
	for _, svc := range cfg.Services {
		if svc.Retention == nil {
			continue
		}
		maxAge := time.Duration(svc.Retention.MaxAgeSecs) * time.Second
		if err := m.Set(svc.Name, maxAge, svc.Retention.MaxItems); err != nil {
			return err
		}
	}
	return nil
}
