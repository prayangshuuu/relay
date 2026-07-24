package provider

import (
	"testing"

	"github.com/prayangshuuu/relay/internal/config"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.ProviderConfig
		wantErr bool
	}{
		{
			name: "Valid Config",
			cfg: config.ProviderConfig{
				ID:      "test-provider",
				Name:    "Test Provider",
				Type:    "openai",
				BaseURL: "https://api.test.com/v1",
			},
			wantErr: false,
		},
		{
			name: "Missing ID",
			cfg: config.ProviderConfig{
				Name: "Test Provider",
				Type: "openai",
			},
			wantErr: true,
		},
		{
			name: "Invalid ID",
			cfg: config.ProviderConfig{
				ID:   "Test_Provider", // uppercase not allowed
				Name: "Test",
				Type: "openai",
			},
			wantErr: true,
		},
		{
			name: "Reserved ID",
			cfg: config.ProviderConfig{
				ID:   "custom",
				Name: "Test",
				Type: "openai",
			},
			wantErr: true,
		},
		{
			name: "Missing Name",
			cfg: config.ProviderConfig{
				ID:   "test",
				Type: "openai",
			},
			wantErr: true,
		},
		{
			name: "Missing Type",
			cfg: config.ProviderConfig{
				ID:   "test",
				Name: "Test",
			},
			wantErr: true,
		},
		{
			name: "Invalid URL",
			cfg: config.ProviderConfig{
				ID:      "test",
				Name:    "Test",
				Type:    "openai",
				BaseURL: "not-a-url",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Validate(tt.cfg); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
