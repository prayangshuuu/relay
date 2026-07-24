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
)

type Match struct {
	Name string
	Type EntityType
}

type Resolver struct {
	paths config.PathManager
}

func NewResolver(paths config.PathManager) *Resolver {
	return &Resolver{paths: paths}
}

func (r *Resolver) exists(dir, name string) bool {
	_, err := os.Stat(filepath.Join(dir, fmt.Sprintf("%s.yaml", name)))
	return err == nil
}

func (r *Resolver) Resolve(arg string) (Match, error) {
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

	// Special case: if absolutely nothing matches, we can assume they are switching to an
	// on-the-fly tool, but only if they explicitly specify that behavior.
	// For 'relay use', it's safer to require the tool configuration to exist or error out.
	if len(matches) == 0 {
		return Match{}, fmt.Errorf("could not resolve '%s' to any profile, provider, or tool", arg)
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
