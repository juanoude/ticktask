package gkeyring

import (
	"github.com/zalando/go-keyring"
)

type KeyringWallet struct{}

const (
	service = "ticktask-cli"
)

func GetWallet() *KeyringWallet {
	return &KeyringWallet{}
}

func (wallet *KeyringWallet) StoreKey(key string, value string) error {
	return keyring.Set(service, key, value)
}

func (wallet *KeyringWallet) GetKey(key string) (string, error) {
	return keyring.Get(service, key)
}
