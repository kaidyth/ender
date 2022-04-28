package client

import (
	"github.com/99designs/keyring"
)

var (
	keychainName = "com.kaidyth.ender"
)

func GetKeyring(chest string) (keyring.Keyring, error) {
	if chest == "" {
		chest = "default"
	}
	ring, err := keyring.Open(keyring.Config{
		ServiceName:                    keychainName + "." + chest,
		KeychainName:                   keychainName + "." + chest,
		KeyCtlScope:                    "user",
		KeyCtlPerm:                     0x3f3f0000,
		KeychainTrustApplication:       true,
		KeychainSynchronizable:         true,
		KeychainAccessibleWhenUnlocked: false,
		AllowedBackends: []keyring.BackendType{
			keyring.KeyCtlBackend,
			keyring.KeychainBackend,
		},
	})

	return ring, err
}
