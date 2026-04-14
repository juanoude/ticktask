// Package gkeyring implements the WalletLayer interface using the system keyring.
// It provides secure storage for sensitive credentials (like AWS access keys)
// using the operating system's native secret storage:
//   - Linux: Secret Service API (GNOME Keyring, KWallet)
//   - macOS: Keychain
//   - Windows: Credential Manager
package gkeyring

import (
	"github.com/zalando/go-keyring"
)

// KeyringWallet implements WalletLayer using the system keyring.
type KeyringWallet struct{}

// service is the application identifier used in the system keyring.
// All TickTask credentials are stored under this service name.
const (
	service = "ticktask-cli"
)

// GetWallet returns a new KeyringWallet instance.
func GetWallet() *KeyringWallet {
	return &KeyringWallet{}
}

// StoreKey saves a secret value in the system keyring.
// The key is used as the "username" field in the keyring entry.
func (wallet *KeyringWallet) StoreKey(key string, value string) error {
	return keyring.Set(service, key, value)
}

// GetKey retrieves a secret value from the system keyring.
// Returns an error if the key doesn't exist.
func (wallet *KeyringWallet) GetKey(key string) (string, error) {
	return keyring.Get(service, key)
}
