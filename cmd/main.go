package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/proximax-storage/go-xpx-chain-sdk/tools/health"
	"github.com/proximax-storage/go-xpx-chain-sdk/tools/health/packets"
	crypto "github.com/proximax-storage/go-xpx-crypto"
	setup "github.com/proximax-storage/go-xpx-setup-node-util"
)

func main() {
	//////////////////////////////////////////////////////////////////////////
	// prompt user for existing installation directory
	setupDir := setup.PromptDir()
	err := os.Chdir(setupDir)
	if err != nil {
		log.Fatal("Error changing working directory: ", err)
	}

	//////////////////////////////////////////////////////////////////////////
	// connect to the REST server
	restUrl := setup.PromptRestUrl()
	log.Println("Connecting to REST server")
	ctx := context.Background()
	cfg, err := sdk.NewConfig(ctx, []string{restUrl})
	if err != nil {
		log.Fatal("Error connecting to REST server: ", err)
	}
	client := sdk.NewClient(http.DefaultClient, cfg)

	// read configuration files
	log.Println("Reading configuration files")
	configUpdater := setup.NewConfigUpdater(setupDir)

	//////////////////////////////////////////////////////////////////////////
	// validate boot key and generate new one if it's invalid
	log.Println("Validating boot key")
	bootKeyChecker := setup.NewBootKeyChecker(client)
	if !bootKeyChecker.Check(ctx, configUpdater.BootKey, configUpdater.HarvestKey) {
		log.Print("Boot key is invalid. Generating new one...")
		configUpdater.BootKey = bootKeyChecker.Generate()
	}

	//////////////////////////////////////////////////////////////////////////
	// check if harvest key is set, if yes then register it if it's not registered yet
	harvestKeyChecker := setup.NewHarvestKeyChecker(client)
	log.Println("Checking harvest public key registration")
	if setup.IsValidPrivateKey(client, configUpdater.HarvestKey) && !harvestKeyChecker.IsRegistered(ctx, configUpdater.HarvestKey) {
		log.Println("Harvest key is not registered")
		if setup.PromptConfirmation("Register harvest key?") {
			signerPrivateKey := setup.PromptKey("Enter private key to sign the transaction to register the harvest key")
			harvestKeyChecker.Register(ctx, cfg, configUpdater.HarvestKey, signerPrivateKey)
		}
	}

	//////////////////////////////////////////////////////////////////////////
	// prompt user for DBRB connection port
	configUpdater.DbrbPort = setup.PromptDbrbPort("Enter port number for DBRB communications", configUpdater)

	//////////////////////////////////////////////////////////////////////////
	// save updated config files in a separate directory
	log.Println("Saving new configuration")
	configUpdater.SaveNewConfig()

	//////////////////////////////////////////////////////////////////////////
	// get the upgrade height
	log.Println("Checking network version")
	networkVersion, err := client.Network.GetNetworkVersionAtHeight(ctx, setup.MaxHeight)
	if err != nil {
		log.Fatal("Error getting network version: ", err)
	}

	requiredNetworkVersion := setup.ParseBlockchainVersion(setup.RequiredBlockchainVersion)
	if networkVersion.BlockChainVersion < requiredNetworkVersion {
		log.Fatal("Expected network version ", requiredNetworkVersion.String(), " is not yet set")
	}

	//////////////////////////////////////////////////////////////////////////
	// Wait for the upgrade height
	if setup.PromptConfirmation("Check if network reached the upgrade height?") {
		bootAccount, err := client.NewAccountFromPrivateKey(configUpdater.BootKey)
		if err != nil {
			log.Fatal("Error creating boot account: ", err)
		}

		nodeInfos := []*health.NodeInfo{
			{
				IdentityKey: bootAccount.KeyPair.PublicKey,
				Endpoint:    "localhost:" + configUpdater.Port,
			},
		}

		keyPair, err := crypto.NewRandomKeyPair()
		if err != nil {
			log.Fatal("Error creating random key pair: ", err)
		}

		for {
			pool, err := health.NewNodeHealthCheckerPool(keyPair, nodeInfos, packets.NoneConnectionSecurity, true, setup.MaxPeerConnections)
			if err != nil {
				log.Printf("Error creating health checker pool: %s\n", err)
				if setup.PromptConfirmation("Retry check if network reached the upgrade height?") {
					continue
				} else {
					break
				}
			}

			err = pool.WaitHeightAll(uint64(networkVersion.StartedHeight) - 1)
			if err != nil {
				log.Printf("Error waiting for the upgrade height: %s\n", err)
				if setup.PromptConfirmation("Retry check if network reached the upgrade height?") {
					continue
				} else {
					break
				}
			}

			err = pool.WaitAllHashesEqual(uint64(networkVersion.StartedHeight) - 1)
			if err != nil {
				log.Printf("Error waiting for equal last block hashes: %s\n", err)
				if setup.PromptConfirmation("Retry check if network reached the upgrade height?") {
					continue
				} else {
					break
				}
			}

			setup.PromptInfo("Network has reached the upgrade height, press any key to continue")
			break
		}
	}

	//////////////////////////////////////////////////////////////////////////
	// Stop node
	log.Println("Shutting down node")
	cmd := exec.Command("docker-compose", "down")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(out))
		log.Fatal("Error shutting down node: ", err)
	}

	if _, err = os.Stat(setupDir + "/data/server.lock"); err == nil {
		err = os.Remove(setupDir + "/data/server.lock")
		if err != nil {
			log.Fatal("Error removing server.lock: ", err)
		}
	}

	//////////////////////////////////////////////////////////////////////////
	// Update config
	log.Println("Updating configuration files")
	err = os.Mkdir(setupDir+"/backup", 0755)
	err = os.Rename(setupDir+"/resources", setupDir+"/backup/resources")
	if err != nil {
		log.Fatal("Error backing up config files: ", err)
	}

	err = os.Rename(setupDir+"/new_resources", setupDir+"/resources")
	if err != nil {
		log.Fatal("Error updating config files: ", err)
	}

	//////////////////////////////////////////////////////////////////////////
	// update docker-compose.yml
	log.Println("Updating docker compose script")
	dockerComposeFile := setupDir + "/docker-compose.yml"
	setup.CopyFile(dockerComposeFile, setupDir+"/backup/docker-compose.yml")
	dockerComposeBytes, err := os.ReadFile(dockerComposeFile)
	if err != nil {
		log.Fatal("Error reading docker-compose.yml: ", err)
	}
	dockerComposeLines := strings.Split(string(dockerComposeBytes), "\n")
	for len(dockerComposeLines) > 2 && dockerComposeLines[len(dockerComposeLines)-2] == "" {
		length := len(dockerComposeLines)
		index := length - 2
		dockerComposeLines = append(dockerComposeLines[:index], dockerComposeLines[index+1:]...)
	}
	catapultSectionFound := false
	isCatapultPortsSection := false
	for index, line := range dockerComposeLines {
		if catapultSectionFound {
			if strings.Contains(line, "image") {
				dockerComposeLines[index] = "    image: " + setup.BlockchainDockerImage
			} else if isCatapultPortsSection {
				if !strings.Contains(line, "-") {
					dockerComposeLines = append(dockerComposeLines[:index+1], dockerComposeLines[index:]...)
					dockerComposeLines[index] = "      - " + configUpdater.DbrbPort + ":" + configUpdater.DbrbPort
					break
				} else {
					if strings.Contains(line, configUpdater.DbrbPort) {
						break
					}
				}
			} else if strings.Contains(line, "ports:") {
				isCatapultPortsSection = true
			}
		} else if strings.Contains(line, "catapult:") || strings.Contains(line, "catapult-api-node:") || strings.Contains(line, "mainnet-peer:") {
			catapultSectionFound = true
		}
	}
	if !catapultSectionFound {
		log.Fatal("docker-compose.yml has unexpected format")
	}
	setup.CreateFile(dockerComposeFile, dockerComposeLines)

	log.Printf("Success!")
	log.Printf("Node upgraded to version %s\n", networkVersion.BlockChainVersion.String())
	log.Printf("Node can now be booted up with command: docker-compose up -d")
}
