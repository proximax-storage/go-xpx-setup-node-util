package setup

import (
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

func IsValidPrivateKey(client *sdk.Client, key string) bool {
	if key == ZeroKey {
		return false
	}

	_, err := client.NewAccountFromPrivateKey(key)
	if err != nil {
		return false
	}

	return true
}
