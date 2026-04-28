package proxy

import (
	"fmt"

	"github.com/example/herdctl/internal/config"
)

// LoadFromConfig registers proxy rules from the herd.yaml service definitions.
// Each service may define a proxy block with port and upstream fields.
func LoadFromConfig(p *Proxy, cfg *config.Config) error {
	for _, svc := range cfg.Services {
		if svc.Proxy == nil {
			continue
		}
		if err := p.Register(svc.Name, svc.Proxy.Port, svc.Proxy.Upstream); err != nil {
			return fmt.Errorf("proxy loader: %w", err)
		}
	}
	return nil
}
