package eventbus

import (
	"fmt"

	"github.com/herdctl/herdctl/internal/config"
)

// LoadFromConfig registers default lifecycle event subscriptions derived from
// service hook definitions in the config. Each hook command is wired as a
// handler that prints the event for observability purposes.
func LoadFromConfig(b *Bus, cfg *config.Config, out func(string)) error {
	for _, svc := range cfg.Services {
		for _, ev := range []string{"started", "stopped", "failed", "restarted"} {
			svcName := svc.Name
			evName := ev
			if err := b.Subscribe(svcName, evName, func(s, e string, p map[string]string) {
				if out != nil {
					out(fmt.Sprintf("[eventbus] %s → %s", s, e))
				}
			}); err != nil {
				return fmt.Errorf("eventbus loader: %w", err)
			}
			_ = evName
		}
	}
	return nil
}
