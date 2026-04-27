package quota

import (
	"fmt"

	"github.com/user/herdctl/internal/config"
)

// LoadFromConfig reads quota definitions from the herd config and registers
// them with the provided Manager. Services without quotas are silently skipped.
func LoadFromConfig(m *Manager, cfg *config.Config) error {
	for _, svc := range cfg.Services {
		q := svc.Quota
		if q.CPU > 0 {
			if err := m.Set(svc.Name, KindCPU, q.CPU); err != nil {
				return fmt.Errorf("quota loader: %w", err)
			}
		}
		if q.Memory > 0 {
			if err := m.Set(svc.Name, KindMemory, q.Memory); err != nil {
				return fmt.Errorf("quota loader: %w", err)
			}
		}
		if q.Procs > 0 {
			if err := m.Set(svc.Name, KindProcs, float64(q.Procs)); err != nil {
				return fmt.Errorf("quota loader: %w", err)
			}
		}
	}
	return nil
}
