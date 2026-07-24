package keyring

import (
	"github.com/zalando/go-keyring"
)

const serviceName = "relay.providers"

// Manager defines the interface for interacting with the OS credential store.
type Manager interface {
	Set(instanceID, secret string) error
	Get(instanceID string) (string, error)
	Delete(instanceID string) error
}

type osKeyringManager struct{}

// NewManager returns a keyring manager backed by the native OS credential store.
func NewManager() Manager {
	return &osKeyringManager{}
}

func (m *osKeyringManager) Set(instanceID, secret string) error {
	return keyring.Set(serviceName, instanceID, secret)
}

func (m *osKeyringManager) Get(instanceID string) (string, error) {
	return keyring.Get(serviceName, instanceID)
}

func (m *osKeyringManager) Delete(instanceID string) error {
	return keyring.Delete(serviceName, instanceID)
}
