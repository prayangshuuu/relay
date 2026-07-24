package setup

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/prayangshuuu/relay/internal/config"
	"github.com/prayangshuuu/relay/internal/profile"
	"github.com/prayangshuuu/relay/internal/provider"
	"github.com/prayangshuuu/relay/internal/tool"
)

type Wizard struct {
	cfgMgr  *config.Manager
	profMgr *profile.Manager
	provMgr *provider.Manager
	toolMgr *tool.Manager
}

func NewWizard(c *config.Manager, p *profile.Manager, pr *provider.Manager, t *tool.Manager) *Wizard {
	return &Wizard{c, p, pr, t}
}

func (w *Wizard) Run() {
	fmt.Println("==================================================")
	fmt.Println("             Welcome to Relay Setup               ")
	fmt.Println("==================================================")

	reader := bufio.NewReader(os.Stdin)

	// Step 1: Detect tools
	fmt.Println("\n[Step 1] Detecting installed AI tools...")
	commonTools := []string{"claude", "aider", "codex", "sgpt"}
	var foundTools []string
	for _, t := range commonTools {
		if _, err := exec.LookPath(t); err == nil {
			foundTools = append(foundTools, t)
			fmt.Printf("  ✓ Found: %s\n", t)
		}
	}
	if len(foundTools) == 0 {
		fmt.Println("  ⚠ No common AI tools detected in PATH.")
	}

	// Step 2: Select preferred tool
	fmt.Println("\n[Step 2] Select preferred tool")
	var selectedTool string
	for {
		fmt.Print("Enter tool executable name (e.g., claude, aider): ")
		input, _ := reader.ReadString('\n')
		selectedTool = strings.TrimSpace(input)
		if selectedTool != "" {
			break
		}
	}

	// We only add it to the manager if we actually implemented Add for tools.
	// Since tool manager fallback auto-creates tools, we don't strictly need to write a YAML yet,
	// but let's assume it gets created naturally or via a future Add command.
	// We'll skip w.toolMgr.Add for now since it doesn't exist yet and just bind the profile to it.

	// Step 3: Create provider
	fmt.Println("\n[Step 3] Select Provider Type")
	fmt.Println("Common types: anthropic, openrouter, openai, ollama")
	var provType string
	for {
		fmt.Print("Enter provider type: ")
		input, _ := reader.ReadString('\n')
		provType = strings.TrimSpace(input)
		if provType != "" {
			break
		}
	}

	// Step 4: Create provider instance
	fmt.Println("\n[Step 4] Configure Provider Instance")
	fmt.Print("Enter a name for this provider instance (e.g., work, personal) [default]: ")
	input, _ := reader.ReadString('\n')
	provName := strings.TrimSpace(input)
	if provName == "" {
		provName = "default"
	}

	fmt.Print("Enter API Key (leave blank to skip): ")
	apiKey, _ := reader.ReadString('\n')
	apiKey = strings.TrimSpace(apiKey)

	_ = apiKey // API Key logic would route to a Keyring securely here

	w.provMgr.Add(config.ProviderConfig{
		ID:   provName,
		Name: provName,
		Type: provType,
	})

	// Step 5: Create profile
	fmt.Println("\n[Step 5] Create Profile")
	fmt.Print("Enter a name for this profile (e.g., work, personal) [default]: ")
	input, _ = reader.ReadString('\n')
	profName := strings.TrimSpace(input)
	if profName == "" {
		profName = "default"
	}

	fmt.Print("Enter Default Model (e.g., claude-3-5-sonnet-20240620): ")
	model, _ := reader.ReadString('\n')
	model = strings.TrimSpace(model)

	w.profMgr.Add(config.ProfileConfig{
		Name:     profName,
		Tool:     selectedTool,
		Provider: provName,
		Model:    model,
	})

	// Step 6: Set default profile
	fmt.Println("\n[Step 6] Setting default profile...")
	globalCfg, _ := w.cfgMgr.Load()
	globalCfg.CurrentProfile = profName
	globalCfg.DefaultTool = selectedTool
	w.cfgMgr.Save(globalCfg)

	fmt.Println("\n==================================================")
	fmt.Println("             Relay is ready.                      ")
	fmt.Println("==================================================")
	fmt.Println("\nRun 'relay doctor' to validate your installation.")
}
