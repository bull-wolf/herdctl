package priority

import "github.com/seanmorris/herdctl/internal/config"

// LoadFromConfig reads priority values from the herd config and registers
// them with the provided Manager. Services without an explicit priority
// are left at the default (Normal).
func LoadFromConfig(m *Manager, cfg *config.Config) error {
	for name, svc := range cfg.Services {
		if svc.Priority == 0 {
			continue
		}
		if err := m.Set(name, Level(svc.Priority)); err != nil {
			return err
		}
	}
	return nil
}
