package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "herd.yaml")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoad_Valid(t *testing.T) {
	path := writeTemp(t, `
version: "1"
services:
  web:
    command: go run main.go
    env:
      PORT: "3000"
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(cfg.Services))
	}
	if cfg.Services["web"].Env["PORT"] != "3000" {
		t.Errorf("expected PORT=3000")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/herd.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_MissingCommand(t *testing.T) {
	path := writeTemp(t, `
version: "1"
services:
  broken:
    dir: ./somewhere
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing command")
	}
}

func TestLoad_UnknownDependency(t *testing.T) {
	path := writeTemp(t, `
version: "1"
services:
  api:
    command: go run main.go
    depends_on:
      - db
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for unknown dependency")
	}
}
