package hooks

import (
	"fmt"
	"os/exec"

	"github.com/user/herdctl/internal/config"
)

// LoadFromConfig registers shell-command hooks derived from the herd.yaml
// service definitions into the provided Registry.
//
// Each service may define hook commands under:
//   hooks:
//     before_start: "echo starting"
//     after_start:  "echo started"
//     before_stop:  "echo stopping"
//     after_stop:   "echo stopped"
func LoadFromConfig(cfg *config.Config, r *Registry) error {
	for _, svc := range cfg.Services {
		svc := svc // capture
		for rawEvent, cmdStr := range svc.Hooks {
			event := Event(rawEvent)
			switch event {
			case EventBeforeStart, EventAfterStart, EventBeforeStop, EventAfterStop:
				// valid
			default:
				return fmt.Errorf("service %q: unknown hook event %q", svc.Name, rawEvent)
			}
			cmd := cmdStr
			r.Register(svc.Name, event, func(service string, ev Event) error {
				return exec.Command("sh", "-c", cmd).Run()
			})
		}
	}
	return nil
}
