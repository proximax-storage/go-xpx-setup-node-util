package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
)

const ()

type nodeData struct {
	harvestKey   string
	bootKey      string
	friendlyName string
	nodeHost     string
}

func main() {
	// prompt user for existing installation directory
	setupDir := promptDir()
	// retrieve data
	existingNodeData, err := getNodeData(setupDir)
	if err != nil {
		log.Fatal(err)
	}
	// check node boot key
	if !checkBootKey(existingNodeData) {
		// replace with new bootkey
	}

}

func promptDir() string {
	validate := func(path string) error {
		_, err := os.Stat(path)
		if err == nil {
			return nil
		}
		return err
	}

	prompt := promptui.Prompt{
		Label:    "Enter base installation directory",
		Validate: validate,
		Default:  "/mnt/siriuschain",
	}
	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}
	return result
}

func getNodeData(path string) (nodeData, error) {
	var data nodeData
	var err error
	data.harvestKey, err = getKeyValue(path+"/resources/config-harvestKey.properties", "harvestKey")
	if err != nil {
		return data, err
	}
	data.friendlyName, err = getKeyValue(path+"/resources/config-node.properties", "friendlyName")
	if err != nil {
		return data, err
	}
	data.nodeHost, err = getKeyValue(path+"/resources/config-node.properties", "host")
	if err != nil {
		return data, err
	}
	data.bootKey, err = getKeyValue(path+"/resources/config-user.properties", "bootKey")
	if err != nil {
		return data, err
	}
	return data, nil
}

func getKeyValue(configFile string, key string) (string, error) {
	// get value of the key from config file
	f, err := os.Open(configFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		s := strings.Split(scanner.Text(), "=")
		res1, _ := regexp.MatchString(key, s[0])
		if res1 {
			return strings.TrimSpace(s[1]), nil

		}
	}
	return nil, errors.New(key + " not found in " + configFile)
}

func checkBootKey(nodeData) bool {
	if nodeData.bootKey == nodeData.harvestKey {
		return false
	}

	if isMultiSig(nodeData.bootKey) {
		return false
	}

	return true
}
