package switcher

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/prayangshuuu/relay/internal/config"
)

type EntityType string

const (
	TypeProfile  EntityType = "profile"
	TypeProvider EntityType = "provider"
	TypeTool     EntityType = "tool"
	TypeModel    EntityType = "model"
)

type Match struct {
	Name string
	Type EntityType
}

type Resolver struct {
	paths   config.PathManager
	aliases map[string]string
}

func NewResolver(paths config.PathManager, aliases map[string]string) *Resolver {
	return &Resolver{paths: paths, aliases: aliases}
}

func (r *Resolver) exists(dir, name string) bool {
	_, err := os.Stat(filepath.Join(dir, fmt.Sprintf("%s.yaml", name)))
	return err == nil
}

func (r *Resolver) Resolve(arg string) (Match, error) {
	// 1. Resolve Aliases
	visited := make(map[string]bool)
	for {
		if target, ok := r.aliases[arg]; ok {
			if visited[target] {
				break // prevent infinite alias loops
			}
			visited[target] = true
			arg = target
		} else {
			break
		}
	}

	var matches []Match

	// Check Profile
	if r.exists(r.paths.ProfilesDir(), arg) {
		matches = append(matches, Match{Name: arg, Type: TypeProfile})
	}

	// Check Provider
	if r.exists(r.paths.ProvidersDir(), arg) {
		matches = append(matches, Match{Name: arg, Type: TypeProvider})
	}

	// Check Tool
	if r.exists(r.paths.ToolsDir(), arg) {
		matches = append(matches, Match{Name: arg, Type: TypeTool})
	}

	// Fallback: If nothing on disk matches, assume it's a Model override
	if len(matches) == 0 {
		return Match{Name: arg, Type: TypeModel}, nil
	}

	if len(matches) == 1 {
		return matches[0], nil
	}

	// Ambiguity Resolution
	fmt.Printf("Target '%s' matched multiple configurations:\n", arg)
	for i, m := range matches {
		fmt.Printf("  %d: %s\n", i+1, m.Type)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Which did you mean? (enter number): ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		var selection int
		_, err := fmt.Sscanf(text, "%d", &selection)
		if err == nil && selection >= 1 && selection <= len(matches) {
			return matches[selection-1], nil
		}
		fmt.Println("Invalid selection. Please try again.")
	}
}
