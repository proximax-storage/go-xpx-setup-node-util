package setup

import (
	"context"
	"log"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	crypto "github.com/proximax-storage/go-xpx-crypto"
)

type BootKeyChecker struct {
	Client *sdk.Client
}

func NewBootKeyChecker(client *sdk.Client) *BootKeyChecker {
	return &BootKeyChecker{
		Client: client,
	}
}

func (b *BootKeyChecker) Check(ctx context.Context, bootKey string, harvestKey string) bool {
	return IsValidPrivateKey(b.Client, bootKey) && bootKey != harvestKey && !b.IsMultiSig(ctx, bootKey)
}

func (b *BootKeyChecker) IsMultiSig(ctx context.Context, bootKey string) bool {
	account, err := b.Client.NewAccountFromPrivateKey(bootKey)
	if err != nil {
		log.Fatal("Error creating account: ", err)
	}

	multiSig, _ := b.Client.Account.GetMultisigAccountInfo(ctx, account.Address)

	return multiSig != nil
}

func (b *BootKeyChecker) Generate() (bootKey string) {
	keyPair, err := crypto.NewRandomKeyPair()
	if err != nil {
		log.Fatal("Error generating new boot key: ", err)
	}

	log.Println("Generated new boot key")

	return keyPair.PrivateKey.String()
}
