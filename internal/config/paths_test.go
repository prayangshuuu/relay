package config

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestOSPathManager(t *testing.T) {
	pm, err := NewOSPathManager()
	if err != nil {
		t.Fatalf("Failed to create OSPathManager: %v", err)
	}

	baseDir := pm.BaseDir()
	if baseDir == "" {
		t.Error("BaseDir should not be empty")
	}

	if runtime.GOOS == "linux" {
		if !strings.HasSuffix(baseDir, "relay") {
			t.Errorf("Expected base directory to end with 'relay' on linux, got %s", baseDir)
		}
	} else {
		if !strings.HasSuffix(baseDir, "Relay") {
			t.Errorf("Expected base directory to end with 'Relay' on %s, got %s", runtime.GOOS, baseDir)
		}
	}

	if pm.ConfigFile() != filepath.Join(baseDir, "config.yaml") {
		t.Errorf("Unexpected ConfigFile path: %s", pm.ConfigFile())
	}
	if pm.ProvidersDir() != filepath.Join(baseDir, "providers") {
		t.Errorf("Unexpected ProvidersDir path: %s", pm.ProvidersDir())
	}
	if pm.ProfilesDir() != filepath.Join(baseDir, "profiles") {
		t.Errorf("Unexpected ProfilesDir path: %s", pm.ProfilesDir())
	}
	if pm.ToolsDir() != filepath.Join(baseDir, "tools") {
		t.Errorf("Unexpected ToolsDir path: %s", pm.ToolsDir())
	}
	if pm.BackupsDir() != filepath.Join(baseDir, "backups") {
		t.Errorf("Unexpected BackupsDir path: %s", pm.BackupsDir())
	}
	if pm.LogsDir() != filepath.Join(baseDir, "logs") {
		t.Errorf("Unexpected LogsDir path: %s", pm.LogsDir())
	}
}
