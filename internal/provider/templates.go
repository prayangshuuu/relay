package provider

import "github.com/prayangshuuu/relay/internal/config"

// Templates returns the built-in provider configurations.
func Templates() map[string]config.ProviderConfig {
	return map[string]config.ProviderConfig{
		"anthropic": {
			ID:                   "anthropic",
			Name:                 "Anthropic",
			Type:                 "anthropic",
			BaseURL:              "https://api.anthropic.com/v1",
			AuthenticationType:   "header",
			EnvironmentVariables: []string{"ANTHROPIC_API_KEY"},
			Enabled:              true,
		},
		"openrouter": {
			ID:                   "openrouter",
			Name:                 "OpenRouter",
			Type:                 "openai-compatible",
			BaseURL:              "https://openrouter.ai/api/v1",
			AuthenticationType:   "bearer",
			EnvironmentVariables: []string{"OPENROUTER_API_KEY"},
			Enabled:              true,
		},
		"agentrouter": {
			ID:                   "agentrouter",
			Name:                 "AgentRouter",
			Type:                 "openai-compatible",
			BaseURL:              "https://agentrouter.ai/api/v1",
			AuthenticationType:   "bearer",
			EnvironmentVariables: []string{"AGENTROUTER_API_KEY"},
			Enabled:              true,
		},
		"openai": {
			ID:                   "openai",
			Name:                 "OpenAI",
			Type:                 "openai",
			BaseURL:              "https://api.openai.com/v1",
			AuthenticationType:   "bearer",
			EnvironmentVariables: []string{"OPENAI_API_KEY"},
			Enabled:              true,
		},
		"google-ai": {
			ID:                   "google-ai",
			Name:                 "Google AI",
			Type:                 "google",
			BaseURL:              "https://generativelanguage.googleapis.com/v1beta",
			AuthenticationType:   "header",
			EnvironmentVariables: []string{"GEMINI_API_KEY"},
			Enabled:              true,
		},
		"groq": {
			ID:                   "groq",
			Name:                 "Groq",
			Type:                 "openai-compatible",
			BaseURL:              "https://api.groq.com/openai/v1",
			AuthenticationType:   "bearer",
			EnvironmentVariables: []string{"GROQ_API_KEY"},
			Enabled:              true,
		},
		"ollama": {
			ID:                   "ollama",
			Name:                 "Ollama",
			Type:                 "openai-compatible",
			BaseURL:              "http://localhost:11434/v1",
			AuthenticationType:   "none",
			EnvironmentVariables: []string{},
			Enabled:              true,
		},
		"lm-studio": {
			ID:                   "lm-studio",
			Name:                 "LM Studio",
			Type:                 "openai-compatible",
			BaseURL:              "http://localhost:1234/v1",
			AuthenticationType:   "none",
			EnvironmentVariables: []string{},
			Enabled:              true,
		},
	}
}
