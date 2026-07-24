package env

import (
	"strings"
	"testing"
)

func TestOSEnvironment_Build(t *testing.T) {
	e := NewOSEnvironment()

	injections := map[string]string{
		"RELAY_TEST_KEY": "relay-value",
	}

	result := e.Build(injections)

	found := false
	for _, kv := range result {
		if kv == "RELAY_TEST_KEY=relay-value" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected RELAY_TEST_KEY to be injected")
	}

	// Ensure system variables like PATH are preserved
	hasPath := false
	for _, kv := range result {
		if strings.HasPrefix(kv, "PATH=") || strings.HasPrefix(kv, "Path=") {
			hasPath = true
			break
		}
	}

	if !hasPath {
		t.Error("Expected system PATH to be preserved")
	}
}
