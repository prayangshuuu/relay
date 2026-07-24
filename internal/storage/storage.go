// Package storage handles safe, cross-platform file operations.
// Future responsibilities include atomic writes, safe directory creation,
// and reading/writing configuration files without corruption or race conditions.
package storage

// Storage provides an abstraction over file system interactions.
type Storage interface {
	// ReadYAML reads a file from the specified path and unmarshals it into out.
	ReadYAML(path string, out interface{}) error

	// WriteYAML marshals the 'in' interface to YAML and writes it atomically
	// to the specified path to prevent file corruption.
	WriteYAML(path string, in interface{}) error

	// EnsureDir safely creates a directory and its parents if they do not exist.
	EnsureDir(path string) error
}
