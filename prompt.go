package setup

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/manifoldco/promptui"
	crypto "github.com/proximax-storage/go-xpx-crypto"
)

func PromptDir() string {
	return promptBase("Enter base installation directory", DefaultInstallationDirectory, "Invalid base installation directory", func(path string) error {
		_, err := os.Stat(path)
		return err
	})
}

func PromptRestUrl() string {
	return promptBase("Enter REST server URL", DefaultRestUrl, "Invalid REST server URL", func(url string) error {
		return nil
	})
}

func PromptKey(label string) string {
	return promptBase(label, "", "Invalid key", func(input string) error {
		_, err := crypto.NewPrivateKeyfromHexString(input)
		return err
	})
}

func PromptDbrbPort(label string, configUpdater *ConfigUpdater) string {
	return promptBase(label, "7903", "Invalid port", func(input string) error {
		if input == configUpdater.SubscriberPort {
			return errors.New("port value is already set as subscriber port in config-messaging.properties")
		}

		if input == configUpdater.Port {
			return errors.New("port value is already set as P2P port in config-node.properties")
		}

		if input == configUpdater.ApiPort {
			return errors.New("port value is already set as API port in config-node.properties")
		}

		port, err := strconv.Atoi(input)
		if err != nil {
			return err
		}

		if port < 1024 || port > 65535 {
			return errors.New("port is not in recommended range 1024-65535")
		}

		return nil
	})
}

func PromptConfirmation(label string) bool {
	result := promptBase(label+" (y/n)", "", "", func(input string) error {
		if input != "y" && input != "n" {
			return errors.New("y/n expected")
		}

		return nil
	})

	return result == "y"
}

func PromptInfo(label string) {
	promptBase(label, "", "", func(input string) error {
		return nil
	})
}

func promptBase(label string, defaultValue string, errorMessage string, validate promptui.ValidateFunc) string {
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
		Default:  defaultValue,
	}

	result, err := prompt.Run()
	if err != nil {
		log.Fatal(errorMessage, ": ", err)
	}

	return result
}
