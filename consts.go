package setup

const (
	MaxPeerConnections           = 10
	MaxHeight                    = 9007199254740991 // 2^53 − 1, max height accepted by REST because of JS constraint
	ZeroKey                      = "0000000000000000000000000000000000000000000000000000000000000000"
	DbrbConfig                   = "[dbrb]\n\ntransactionTimeout = 1h"
	DefaultInstallationDirectory = "/mnt/siriuschain"
	DefaultRestUrl               = "https://arcturus.xpxsirius.io"
	BlockchainDockerImage        = "proximax/proximax-sirius-chain:v1.4.2-bullseye"
	NewBlockchainVersion         = "1.4.2.0"
)
