package launcher

import (
	"testing"
)

// Mock objects
type MockEnv struct {
	lastEnv []string
}

func (m *MockEnv) Build(injections map[string]string) []string {
	m.lastEnv = []string{"MOCK=true"}
	for k, v := range injections {
		m.lastEnv = append(m.lastEnv, k+"="+v)
	}
	return m.lastEnv
}
func (m *MockEnv) Clean() error { return nil }

type MockShell struct {
	lastExe  string
	lastArgs []string
	lastEnv  []string
}

func (m *MockShell) Exec(executable string, args []string, env []string) error {
	m.lastExe = executable
	m.lastArgs = args
	m.lastEnv = env
	return nil
}

type MockTool struct {
	name string
}

func (m *MockTool) ExecutableName() string                  { return m.name }
func (m *MockTool) SupportedEnvironmentVariables() []string { return nil }
func (m *MockTool) LaunchMethod() string                    { return "exec" }

func TestEngine_DryRun(t *testing.T) {
	e := NewEngine(&MockEnv{}, &MockShell{})

	cfg := Config{
		Tool:   &MockTool{name: "go"}, // Use 'go' as it's guaranteed to be in PATH
		DryRun: true,
	}

	if err := e.Launch(cfg); err != nil {
		t.Fatalf("Launch failed: %v", err)
	}
}

func TestEngine_Execution(t *testing.T) {
	shell := &MockShell{}
	e := NewEngine(&MockEnv{}, shell)

	cfg := Config{
		Tool: &MockTool{name: "go"},
		Args: []string{"version"},
	}

	if err := e.Launch(cfg); err != nil {
		t.Fatalf("Launch failed: %v", err)
	}

	if shell.lastArgs[0] != "version" {
		t.Errorf("Expected args to contain 'version', got %v", shell.lastArgs)
	}
}
