package cascade

import (
	"fmt"

	"github.com/user/herdctl/internal/config"
)

// LoadFromConfig registers cascade rules from the parsed herd.yaml config.
// Each service may define a cascade block with a policy and list of targets.
func LoadFromConfig(m *Manager, cfg *config.Config) error {
	for _, svc := range cfg.Services {
		if svc.Cascade == nil {
			continue
		}
		policy := Policy(svc.Cascade.Policy)
		if policy == "" {
			policy = PolicyIgnore
		}
		if err := m.Register(svc.Name, policy, svc.Cascade.Targets); err != nil {
			return fmt.Errorf("cascade: loading service %q: %w", svc.Name, err)
		}
	}
	return nil
}
