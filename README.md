# Mainnet node setup CLI tool

The tool upgrades a mainnet node from v0.9.0 to v1.3.0.

## Build

```shell
cd cmd && go build -o setup_util
```

## Input parameters

After running the tool asks for a few parameters.

- Directory where the node is installed. Default: /mnt/siriuschain
- REST server URL. Default: http://aldebaran.xpxsirius.io:3000
- (optional) Private key of account on behalf of which the harvest key will be registered. Required when harvest key is set and is not yet registered in the blockchain and user consents to register it.
- Port number for DBRB communications. Default: 7903

## Run

```shell
./setup_util
```