// Package storage provides the abstractions for credential management.
package storage

import "errors"

var ErrCredentialNotFound = errors.New("credential not found")

// CredentialManager defines the interface for secure secret storage.
// Future responsibilities include encryption and OS-level keychain integration.
type CredentialManager interface {
	// Save stores a secret value under a specific key securely.
	Save(key string, secret string) error

	// Get retrieves a securely stored secret by key.
	Get(key string) (string, error)
}

// MockCredentialManager is a temporary implementation until encryption is added.
type MockCredentialManager struct {
	store map[string]string
}

// NewMockCredentialManager creates a simple in-memory credential manager.
func NewMockCredentialManager() *MockCredentialManager {
	return &MockCredentialManager{
		store: make(map[string]string),
	}
}

// Save stores the secret in memory.
func (m *MockCredentialManager) Save(key string, secret string) error {
	m.store[key] = secret
	return nil
}

// Get retrieves the secret from memory.
func (m *MockCredentialManager) Get(key string) (string, error) {
	val, ok := m.store[key]
	if !ok {
		return "", ErrCredentialNotFound
	}
	return val, nil
}
