// Package storage handles safe, cross-platform file operations.
package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LocalStorage implements Storage using the local filesystem.
type LocalStorage struct{}

// NewLocalStorage creates a new LocalStorage.
func NewLocalStorage() *LocalStorage {
	return &LocalStorage{}
}

// ReadYAML reads and unmarshals a YAML file.
func (s *LocalStorage) ReadYAML(path string, out interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}

	if err := yaml.Unmarshal(data, out); err != nil {
		return fmt.Errorf("failed to parse YAML in %s: %w", path, err)
	}

	return nil
}

// WriteYAML atomically writes the interface as YAML to the specified path.
func (s *LocalStorage) WriteYAML(path string, in interface{}) error {
	data, err := yaml.Marshal(in)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	dir := filepath.Dir(path)
	if err := s.EnsureDir(dir); err != nil {
		return err
	}

	// Create a temporary file in the same directory for atomic rename
	tempFile, err := os.CreateTemp(dir, ".tmp-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	tempName := tempFile.Name()

	// Ensure temp file is closed and cleaned up in case of failure
	defer func() {
		tempFile.Close()
		os.Remove(tempName)
	}()

	if _, err := tempFile.Write(data); err != nil {
		return fmt.Errorf("failed to write to temporary file: %w", err)
	}

	// Sync to ensure data is on disk before renaming
	if err := tempFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync temporary file: %w", err)
	}

	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("failed to close temporary file: %w", err)
	}

	// Rename is atomic on POSIX and generally safe enough on Windows
	if err := os.Rename(tempName, path); err != nil {
		return fmt.Errorf("failed to atomically rename file to %s: %w", path, err)
	}

	return nil
}

// EnsureDir creates the directory if it does not exist.
func (s *LocalStorage) EnsureDir(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}
	return nil
}
