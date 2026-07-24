package provider

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"

	"github.com/prayangshuuu/relay/internal/config"
)

var (
	ErrValidation       = errors.New("provider validation failed")
	ErrReservedName     = errors.New("reserved provider name")
	validIDRegex        = regexp.MustCompile(`^[a-z0-9_-]+$`)
	reservedProviderIDs = map[string]bool{
		"custom": true,
		"all":    true,
	}
)

// Validate checks the configuration for correctness.
func Validate(cfg config.ProviderConfig) error {
	if cfg.ID == "" {
		return fmt.Errorf("%w: id is required", ErrValidation)
	}

	if !validIDRegex.MatchString(cfg.ID) {
		return fmt.Errorf("%w: id can only contain lowercase letters, numbers, hyphens, and underscores", ErrValidation)
	}

	if reservedProviderIDs[cfg.ID] {
		return fmt.Errorf("%w: '%s' is a reserved id", ErrReservedName, cfg.ID)
	}

	if cfg.Name == "" {
		return fmt.Errorf("%w: name is required", ErrValidation)
	}

	if cfg.Type == "" {
		return fmt.Errorf("%w: type is required", ErrValidation)
	}

	if cfg.BaseURL != "" {
		u, err := url.ParseRequestURI(cfg.BaseURL)
		if err != nil || u.Scheme == "" || u.Host == "" {
			return fmt.Errorf("%w: base_url '%s' is invalid", ErrValidation, cfg.BaseURL)
		}
	}

	return nil
}
