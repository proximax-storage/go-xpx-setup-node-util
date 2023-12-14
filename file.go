package setup

import (
	"io"
	"log"
	"os"
)

func CreateFile(path string, content []string) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatal("Failed to create file ", path, ": ", err)
	}
	defer CloseFile(file)

	for _, line := range content {
		_, err = io.WriteString(file, line+"\n")
		if err != nil {
			log.Fatal("Failed to write to file ", path, ": ", err)
		}
	}
}

func CopyFile(src, dst string) {
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

	_, err = io.Copy(destination, source)
	if err != nil {
		log.Fatal("Failed to copy file from ", src, " to ", dst, ": ", err)
	}
}

func CloseFile(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Fatal("Failed to close file: ", err)
	}
}
