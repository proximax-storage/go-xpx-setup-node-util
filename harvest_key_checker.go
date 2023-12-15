package setup

import (
	"context"
	"log"
	"time"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket"
	sync "github.com/proximax-storage/go-xpx-chain-sync"
)

type HarvestKeyChecker struct {
	Client *sdk.Client
}

func NewHarvestKeyChecker(client *sdk.Client) *HarvestKeyChecker {
	return &HarvestKeyChecker{
		Client: client,
	}
}

func (h *HarvestKeyChecker) IsRegistered(ctx context.Context, harvestKey string) bool {
	harvesterAccount, err := h.Client.NewAccountFromPrivateKey(harvestKey)
	if err != nil {
		log.Fatal("Error creating harvester account: ", err)
	}

	harvester, err := h.Client.Account.GetAccountHarvesting(ctx, harvesterAccount.Address)
	if err != nil {
		log.Fatal("Error requesting harvesting info: ", err)
	}

	return harvester != nil
}

func (h *HarvestKeyChecker) Register(ctx context.Context, cfg *sdk.Config, harvestKey string, signerPrivateKey string) {
	signerAccount, err := h.Client.NewAccountFromPrivateKey(signerPrivateKey)
	if err != nil {
		log.Fatal("Error creating signer account: ", err)
	}

	harvesterAccount, err := h.Client.NewAccountFromPublicKey(harvestKey)
	if err != nil {
		log.Fatal("Error creating harvester account: ", err)
	}

	ws, err := websocket.NewClient(ctx, cfg)
	if err != nil {
		log.Fatal("Error creating websocket client: ", err)
	}

	htx, err := h.Client.NewHarvesterTransaction(
		sdk.NewDeadline(time.Hour),
		sdk.AddHarvester,
		harvesterAccount,
	)

	res, err := sync.Announce(ctx, cfg, ws, signerAccount, htx)
	if err != nil {
		log.Fatal("Error announcing AddHarvesterTransaction: ", err)
	}

	if res.Err() != nil {
		log.Fatal("Error confirming AddHarvesterTransaction: ", res.Err())
	}

	log.Println("Harvest key registered successfully")
}
