// Package doctor provides diagnostics and health checks for Relay.
// Future responsibilities include validating installed tools, executable paths,
// provider configurations, API keys, and overall system health.
package doctor

// DiagnosticResult represents the outcome of a single health check.
type DiagnosticResult struct {
	Name    string
	Passed  bool
	Message string
}

// Doctor is responsible for diagnosing the Relay installation.
type Doctor interface {
	// CheckAll runs a comprehensive suite of diagnostics to ensure
	// the environment, configuration, tools, and providers are healthy.
	CheckAll() []DiagnosticResult
}
