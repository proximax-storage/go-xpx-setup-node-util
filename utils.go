package setup

import (
	"log"
	"strconv"
	"strings"

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

func ParseBlockchainVersion(version string) sdk.BlockChainVersion {
	s := strings.Split(version, ".")

	if len(s) != 4 {
		log.Fatal("Wrong version of blockchain, expected format 2.2.2.2")
	}

	numbers := make([]uint16, len(s))

	for i, value := range s {
		number, err := strconv.Atoi(value)
		if err != nil {
			log.Fatal("Error parsing version number: ", err)
		}
		numbers[i] = uint16(number)
	}

	return sdk.NewBlockChainVersion(numbers[0], numbers[1], numbers[2], numbers[3])
}
