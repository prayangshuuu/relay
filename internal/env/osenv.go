package env

import (
	"fmt"
	"os"
	"strings"
)

// OSEnvironment implements Environment by reading the current host OS environment
// and overlaying custom variables without mutating the host state.
type OSEnvironment struct{}

// NewOSEnvironment creates a new OSEnvironment.
func NewOSEnvironment() *OSEnvironment {
	return &OSEnvironment{}
}

// Build merges the injected variables with the current system environment.
func (e *OSEnvironment) Build(injections map[string]string) []string {
	// Start with the current environment
	current := os.Environ()

	if len(injections) == 0 {
		return current
	}

	// Create a map to efficiently track and override existing variables
	envMap := make(map[string]string)
	var keys []string

	for _, kv := range current {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			envMap[key] = parts[1]
			keys = append(keys, key)
		}
	}

	// Override or add new variables
	for k, v := range injections {
		if _, exists := envMap[k]; !exists {
			keys = append(keys, k)
		}
		envMap[k] = v
	}

	// Rebuild the final array
	result := make([]string, 0, len(keys))
	for _, k := range keys {
		result = append(result, fmt.Sprintf("%s=%s", k, envMap[k]))
	}

	return result
}

// Clean performs any necessary cleanup. For OSEnvironment, this is a no-op
// since we never modified the host system.
func (e *OSEnvironment) Clean() error {
	return nil
}
