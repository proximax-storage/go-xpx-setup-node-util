package setup

import (
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type ConfigUpdater struct {
	SetupDir       string
	HarvestKey     string
	BootKey        string
	SubscriberPort string
	Port           string
	ApiPort        string
	DbrbPort       string
}

func NewConfigUpdater(setupDir string) *ConfigUpdater {
	c := &ConfigUpdater{}
	c.SetupDir = setupDir

	file := openFile(setupDir + "/resources/config-harvesting.properties")
	c.HarvestKey = getKeyValue(file, "harvestKey")
	CloseFile(file)

	file = openFile(setupDir + "/resources/config-node.properties")
	c.Port = getKeyValue(file, "port")
	c.ApiPort = getKeyValue(file, "apiPort")
	CloseFile(file)

	file = openFile(setupDir + "/resources/config-messaging.properties")
	c.SubscriberPort = getKeyValue(file, "subscriberPort")
	CloseFile(file)

	file = openFile(setupDir + "/resources/config-user.properties")
	c.BootKey = getKeyValue(file, "bootKey")
	CloseFile(file)

	return c
}

func (c *ConfigUpdater) SaveNewConfig() {
	skipFiles := make(map[string]bool)
	skipFiles["config-database.properties"] = true
	skipFiles["config-dbrb.properties"] = true
	skipFiles["config-extensions-server.properties"] = true
	skipFiles["config-node.properties"] = true
	skipFiles["config-user.properties"] = true

	copyDir(c.SetupDir+"/resources", c.SetupDir+"/new_resources", skipFiles)

	CreateFile(c.SetupDir+"/new_resources/config-dbrb.properties", []string{DbrbConfig})

	copyFileAndReplace(c.SetupDir+"/resources/config-database.properties", c.SetupDir+"/new_resources/config-database.properties", map[string]string{
		"catapult.mongo.plugins.service":       "false",
		"catapult.mongo.plugins.operation":     "false",
		"catapult.mongo.plugins.supercontract": "false",
		"catapult.mongo.plugins.committee":     "true\ncatapult.mongo.plugins.dbrb = true",
	})
	copyFileAndReplace(c.SetupDir+"/resources/config-extensions-server.properties", c.SetupDir+"/new_resources/config-extensions-server.properties", map[string]string{
		"extension.harvesting":      "false",
		"extension.unbondedpruning": "true\nextension.fastfinality = true",
	})
	copyFileAndReplace(c.SetupDir+"/resources/config-node.properties", c.SetupDir+"/new_resources/config-node.properties", map[string]string{
		"apiPort": c.ApiPort + "\ndbrbPort = " + c.DbrbPort,
	})
	copyFileAndReplace(c.SetupDir+"/resources/config-user.properties", c.SetupDir+"/new_resources/config-user.properties", map[string]string{
		"bootKey": c.BootKey,
	})
}

func openFile(path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Error opening file ", path, ": ", err)
	}

	return file
}

func getKeyValue(file *os.File, key string) string {
	_, err := file.Seek(0, 0)
	if err != nil {
		log.Fatal("Error moving to the beginning of file: ", err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), "=")
		matched, _ := regexp.MatchString(key, s[0])
		if matched {
			return strings.TrimSpace(s[1])
		}
	}

	log.Fatal(key, " not found")

	return ""
}

func copyDir(src, dst string, skipFiles map[string]bool) {
	entries, err := os.ReadDir(src)
	if err != nil {
		log.Fatal("Failed to read directory ", src, ": ", err)
	}

	err = os.RemoveAll(dst)
	if err != nil {
		log.Fatal("Failed to remove directory", dst, ": ", err)
	}

	err = os.Mkdir(dst, 0755)
	if err != nil {
		log.Fatal("Failed to create directory", dst, ": ", err)
	}

	for _, entry := range entries {
		if _, skip := skipFiles[entry.Name()]; !skip {
			CopyFile(filepath.Join(src, entry.Name()), filepath.Join(dst, entry.Name()))
		}
	}
}

func copyFileAndReplace(src, dst string, replacements map[string]string) {
	source, err := os.Open(src)
	if err != nil {
		log.Fatal("Failed to open file ", src, ": ", err)
	}
	defer CloseFile(source)

	destination, err := os.Create(dst)
	if err != nil {
		log.Fatal("Failed to create file ", dst, ": ", err)
	}
	defer CloseFile(destination)

	scanner := bufio.NewScanner(source)
	lines := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		s := strings.Split(line, "=")
		key := strings.TrimSpace(s[0])
		if value, ok := replacements[key]; ok {
			line = key + " = " + value
		}
		newLines := strings.Split(line, "\n")
		for _, newLine := range newLines {
			lines = append(lines, newLine)
		}
	}

	for i, line := range lines {
		if i > 0 {
			previousLine := lines[i-1]
			if line == previousLine {
				continue
			}
		}

		_, err = io.WriteString(destination, line+"\n")
		if err != nil {
			log.Fatal("Failed to write to file ", dst, ": ", err)
		}
	}
}
