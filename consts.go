package setup

const (
	MaxPeerConnections           = 10
	MaxHeight                    = 9007199254740991 // 2^53 − 1, max height accepted by REST because of JS constraint
	ZeroKey                      = "0000000000000000000000000000000000000000000000000000000000000000"
	DbrbConfig                   = "[dbrb]\n\ntransactionTimeout = 1h"
	DefaultInstallationDirectory = "/mnt/siriuschain"
	DefaultRestUrl               = "http://aldebaran.xpxsirius.io:3000"
	// BlockchainDockerImage TODO: replace with the name of the publid image
	BlockchainDockerImage = "249767383774.dkr.ecr.ap-southeast-1.amazonaws.com/proximax-catapult-server:develop-jenkins-build-144"
)
