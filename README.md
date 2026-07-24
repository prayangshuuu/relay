<div align="center">
  <h1>Relay</h1>
  <p><b>Universal AI Provider & Environment Manager</b></p>
  
  > 🚧 **Relay is currently under active development.** Things might break, commands might change, and we are iterating quickly!
  
  [![Go Version](https://img.shields.io/github/go-mod/go-version/prayangshuuu/relay)](https://golang.org/)
  [![License](https://img.shields.io/github/license/prayangshuuu/relay)](LICENSE)
  [![Build](https://img.shields.io/github/actions/workflow/status/prayangshuuu/relay/build.yml)](https://github.com/prayangshuuu/relay/actions)
</div>

---

## Why Relay?

Relay was born out of real, daily developer frustration. 

I regularly use AI coding tools like Claude Code, Aider, and various other AI CLIs. But my workflow isn't simple. Depending on the project, the time of day, or the API limits, I constantly switch between providers like AgentRouter, OpenRouter, Anthropic, OpenAI, and local Ollama instances.

Switching between these providers meant constantly exporting new `API_KEY`s, editing `.bashrc` or Windows environment variables, restarting terminals, and juggling multiple configurations. It was incredibly repetitive. If I wanted to switch from my "work" AgentRouter account to my "personal" OpenRouter account, it broke my flow completely.

I built Relay because I wanted *one simple command* to switch providers, profiles, accounts, and tools instantly—without ever touching global environment variables manually again. 

## Features

### Implemented
- [x] **Cross-platform**: Native Go binary for Windows, macOS, and Linux.
- [x] **Lightweight**: Zero background services, daemons, or telemetry.
- [x] **Temporary Environment Injection**: Secrets are injected directly into the subprocess memory.
- [x] **Provider Management**: Manage templates for Anthropic, OpenAI, OpenRouter, AgentRouter, etc.
- [x] **Unlimited Provider Instances**: Maintain infinite accounts (e.g., `work`, `personal`, `free-tier-1`).
- [x] **Profile Management**: Link tools, instances, and models into cohesive working environments.
- [x] **Smart Switching**: Magically resolve your intent (e.g., `relay use work claude`).
- [x] **Alias Engine**: Add custom shorthand commands (e.g., `relay alias add ar agentrouter`).
- [x] **Secure Credential Storage**: Native integration with OS Keyrings (Windows Credential Manager, macOS Keychain, Linux Secret Service).
- [x] **Interactive Setup**: Terminal wizard to bootstrap your config in seconds.
- [x] **Relay Doctor**: Built-in diagnostic suite with auto-repair capabilities.
- [x] **Extension SDK (Plugins)**: Fully modular architecture for extending tools and providers.

### Planned
- [ ] **Project-local Configuration**: Automatically detect `.relay/` folders in your current working directory.
- [ ] **Importers/Exporters**: Easily share profiles with your team.

---

## Installation

Because Relay is in early development, installation is currently limited to compiling from source or using the Go toolchain.

### Go Install
If you have a Go environment set up:
```bash
go install github.com/prayangshuuu/relay@latest
```

### Build From Source
```bash
git clone https://github.com/prayangshuuu/relay.git
cd relay
go build -o relay main.go
```

---

## Quick Start

Get up and running in less than 30 seconds.

```bash
# 1. Initialize Relay (creates directories and default config)
relay init

# 2. Run the interactive setup wizard to link your tools and providers
relay setup

# 3. Validate your installation and auto-repair any issues
relay doctor

# 4. Create additional provider instances (e.g. 'work' account)
relay instance create

# 5. Create specific runtime profiles
relay profile create

# 6. Launch your AI tool! Relay will securely inject your credentials and exit
relay run claude
```

---

## Usage

Relay provides a comprehensive CLI to manage your entire AI ecosystem.

### Core Commands

- `relay --help` - View all available commands and flags.
- `relay version` - Print current Relay version.
- `relay init` - Scaffolds the `.relay` configuration directory in your OS app data folder.
- `relay setup` - Launches the interactive terminal wizard to configure your first profile.
- `relay doctor` - Runs a diagnostic check on your configuration, tools, and executable paths.
- `relay doctor --fix` - Automatically repairs missing directories or broken configurations.
- `relay current` - Displays your currently active Profile, Provider Instance, Tool, and Model.
- `relay history` - View a chronological list of environments you've recently used.
- `relay undo` - Pops the configuration stack, restoring your previously active profile.
- `relay config` - Manage raw configuration data.
- `relay completion` - Generate shell autocompletion scripts.

### Instance Management (Providers & Accounts)

- `relay instance create` - Interactively create a new provider instance and securely store the API key in your OS Keyring.
- `relay instance list` - List all configured provider instances.
- `relay instance remove [id]` - Delete a provider instance and its credentials.
- `relay instance edit [id]` - Modify an existing instance.
- `relay instance validate [id]` - Verify connectivity for an instance.

### Profile Management

- `relay profile create` - Create a new runtime profile linking a Tool, Instance, and Model.
- `relay profile list` - List all configured profiles.
- `relay profile use [id]` - Explicitly set a specific profile as active.

### Smart Switching

- `relay use [args...]` - The core engine. Pass it any combination of Profile, Instance, Tool, or Model aliases and it will intelligently switch your environment.
- `relay switch [args...]` - Exact alias for `relay use`.

### Tool Management

- `relay tool detect` - Scans your OS `PATH` to automatically discover installed AI CLI tools.
- `relay tool list` - Lists registered tools.
- `relay tool validate [id]` - Verifies the tool's executable is accessible.

### Secrets & Aliasing

- `relay secret set [id]` - Manually set an API key in the OS Keyring.
- `relay secret get [id]` - Retrieve a stored API key.
- `relay secret delete [id]` - Purge an API key from the OS Keyring.
- `relay alias add [alias] [target]` - Create a shorthand alias (e.g., `ar` -> `agentrouter`).
- `relay alias list` - View all active aliases.
- `relay alias remove [alias]` - Delete an alias.

### Plugin SDK

- `relay plugin list` - List all loaded SDK plugins (both built-in and third-party).
- `relay plugin info [id]` - View manifest details for a specific plugin.

---

## Examples

Relay's smart switcher is designed to get out of your way.

**Switching environments on the fly:**
```bash
# Switch to your OpenRouter instance
relay use openrouter

# Switch to your custom 'work' profile
relay switch work

# Switch to the AgentRouter provider and Claude tool simultaneously
relay use agentrouter claude

# Switch to a custom provider instance you created
relay use local-testing
```

**Using Aliases:**
```bash
# Set up shortcuts
relay alias add ar agentrouter
relay alias add sonnet claude-3-5-sonnet-20240620

# Use them instantly
relay use ar sonnet
```

**Launching Tools:**
```bash
# Launch Claude Code using the currently active profile
relay run claude

# Launch Aider or Codex, bypassing the active profile's default tool
relay run aider
relay run codex
```

**Diagnostics & Checks:**
```bash
# Show what your environment currently looks like
relay current

# Run the doctor to diagnose issues
relay doctor

# Test what execution would look like without actually launching the tool
relay run claude --dry-run
relay run claude --verbose
```

---

## Configuration

Relay stores its configuration safely in your user directory.

**Locations:**
- **Windows**: `%APPDATA%/Relay/`
- **macOS**: `~/Library/Application Support/Relay/`
- **Linux**: `~/.config/relay/`

**Structure:**
- `config.yaml`: The global state file containing your active profile, aliases, and history.
- `providers/`: Directory containing your Provider Template definitions.
- `instances/`: Directory containing your configured Accounts (credentials are safely in your OS Keyring).
- `profiles/`: Directory containing your full environment linkings.
- `tools/`: Directory containing tool executable rules and default arguments.
- `plugins/`: Directory where external SDK executables reside.
- `logs/` & `backups/`: Diagnostic and auto-recovery directories.

---

## Core Concepts

- **Provider**: The template for an AI API (e.g., Anthropic, OpenAI). It defines how authentication headers and endpoints are structured.
- **Instance**: A specific account utilizing a Provider. You can have unlimited instances (e.g., `work`, `personal`, `testing`) attached to a single Provider template.
- **Profile**: The ultimate runtime configuration. A profile links a specific Tool, a specific Instance, and a specific Model together alongside custom environment variables.
- **Tool**: The actual executable you want to run (e.g., `claude`, `aider`).
- **Launcher**: The execution engine that merges your existing OS environment with your secure credentials, launches the tool via `os/exec`, and immediately exits.
- **Doctor**: A comprehensive diagnostic suite that ensures your executables exist, your config isn't broken, and auto-repairs missing directories.
- **Smart Switching**: The resolution engine behind `relay use`. It handles ambiguity and recursive aliasing so you can type naturally.
- **Custom Providers**: Providers you define that aren't natively supported. Relay treats them identically to built-in templates.
- **Project Overrides (Planned)**: The ability to drop a `.relay` directory into any code repository to automatically bind that project to a specific Relay profile.

---

## Philosophy

- **No Permanent Environment Modifications**: Relay will *never* permanently inject variables into your Windows Registry or `.bashrc`. It avoids "works on my machine" state drift entirely.
- **No Background Daemons**: Relay does not linger. It figures out what variables you need, spawns the child process, and gets out of the way instantly.
- **No Telemetry**: Your configurations, your API keys, and your usage data are yours. Relay doesn't phone home.
- **Intentionally Lightweight**: Built in pure Go. It's a single binary that stays out of your way.

---

## Performance Goals

- **Small binary**: Statically compiled with standard libraries.
- **Minimal RAM**: Peak memory usage lasts only for milliseconds during execution logic.
- **No background process**: Zero idle CPU usage.
- **Fast startup**: Millisecond initialization time before handing off to the underlying AI tool.
- **Cross-platform**: Identical behavior on Windows, macOS, and Linux.

---

## Security

Security is foundational to Relay.

- **Local-first**: Everything runs entirely on your hardware.
- **OS Keyring Integration**: API Keys are heavily secured using native OS enclaves (Windows Credential Manager, macOS Keychain, Linux Secret Service).
- **Subprocess Injection**: Secrets are injected directly into the ephemeral memory space of the target tool. They are never written to disk or exposed in long-lived shell sessions.

---

## Roadmap

### Current
- [x] Core configuration structure
- [x] Provider & Profile CRUD operations
- [x] Ephemeral launcher engine
- [x] Multi-account smart switcher
- [x] Extension SDK & Plugins
- [x] OS Keyring secret management

### Next
- [ ] Project-local `.relay/config.yaml` directory overrides
- [ ] Profile import/export system for teams
- [ ] Direct Homebrew / Winget distribution

---

## Contributing

We love contributions! Relay is built to solve real workflow problems, and if you have an idea, we want to hear it.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## Development

**Clone the repo:**
```bash
git clone https://github.com/prayangshuuu/relay.git
cd relay
```

**Run locally:**
```bash
go run main.go
```

**Build:**
```bash
go build -o relay main.go
```

**Test:**
```bash
go test ./...
```

---

## License

MIT License. See [LICENSE](LICENSE) for more information.

---

## Acknowledgements

Relay takes immense inspiration from the clean CLI philosophies of tools like Docker, Bun, and Ollama, and relies heavily on the phenomenal Cobra framework from the Go ecosystem.
