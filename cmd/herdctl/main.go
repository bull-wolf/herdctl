package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/herdctl/internal/config"
	"github.com/user/herdctl/internal/deps"
	"github.com/user/herdctl/internal/env"
	"github.com/user/herdctl/internal/health"
	"github.com/user/herdctl/internal/logs"
	"github.com/user/herdctl/internal/runner"
	"github.com/user/herdctl/internal/status"
	"github.com/user/herdctl/internal/status/printer"
)

const defaultConfigFile = "herd.yaml"

func main() {
	configFile := flag.String("config", defaultConfigFile, "path to herd.yaml config file")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}

	cfg, err := config.Load(*configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
	}

	logStore := logs.New()
	statStore := status.New()
	healthChecker := health.New()
	envResolver := env.New(cfg)
	depGraph := deps.New(cfg)
	r := runner.New(cfg, logStore, statStore, envResolver, depGraph)

	cmd := args[0]
	switch cmd {
	case "start":
		services := args[1:]
		if len(services) == 0 {
			// start all services in dependency order
			ordered, err := depGraph.Order()
			if err != nil {
				fmt.Fprintf(os.Stderr, "dependency error: %v\n", err)
				os.Exit(1)
			}
			services = ordered
		}
		for _, svc := range services {
			if err := r.Start(svc); err != nil {
				fmt.Fprintf(os.Stderr, "failed to start %s: %v\n", svc, err)
				os.Exit(1)
			}
			fmt.Printf("started %s\n", svc)
		}

	case "stop":
		services := args[1:]
		if len(services) == 0 {
			for _, svc := range cfg.Services {
				services = append(services, svc.Name)
			}
		}
		for _, svc := range services {
			if err := r.Stop(svc); err != nil {
				fmt.Fprintf(os.Stderr, "failed to stop %s: %v\n", svc, err)
			}
			fmt.Printf("stopped %s\n", svc)
		}

	case "status":
		printer.Print(statStore, os.Stdout)

	case "logs":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "usage: herdctl logs <service> [--tail N]")
			os.Exit(1)
		}
		svc := args[1]
		tail := 20
		for i, a := range args[2:] {
			if a == "--tail" && i+1 < len(args[2:]) {
				fmt.Sscanf(args[2:][i+1], "%d", &tail)
			}
		}
		entries := logStore.TailService(svc, tail)
		for _, e := range entries {
			fmt.Printf("[%s] %s: %s\n", e.Timestamp.Format("15:04:05"), e.Service, strings.TrimRight(e.Line, "\n"))
		}

	case "health":
		if len(args) < 2 {
			// check all services that have health endpoints
			for _, svc := range cfg.Services {
				if svc.HealthURL == "" {
					continue
				}
				result := healthChecker.Check(svc.Name, svc.HealthURL)
				fmt.Printf("%-20s %s\n", svc.Name, result.Status)
			}
		} else {
			svc := args[1]
			result, ok := healthChecker.Get(svc)
			if !ok {
				fmt.Fprintf(os.Stderr, "no health result for %s\n", svc)
				os.Exit(1)
			}
			fmt.Printf("%s: %s\n", svc, result.Status)
		}

	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`herdctl — manage multi-service local dev environments

Usage:
  herdctl [--config herd.yaml] <command> [args]

Commands:
  start [service...]   Start one or all services
  stop  [service...]   Stop one or all services
  status               Show status of all services
  logs  <service>      Tail logs for a service (--tail N, default 20)
  health [service]     Check health endpoints`)
}
