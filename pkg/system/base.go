package system

import (
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/joho/godotenv"
)

const projectDirName = "OptimusVid-Convert" // change to relevant project name

func LoadEnv() {
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))

	err := godotenv.Load(string(rootPath) + `/.env`)

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func EnsureMediaDirectory(fileType string) string {
	exe, err := os.Executable()
	if err != nil {
		log.Panicf("Failed to determine the executable path: %v", err)
	}

	// Get the directory of the executable
	exeDir := filepath.Dir(exe)

	// Traverse two directories up to reach the project root
	var projectRoot string
	if projectRoot = os.Getenv("PROJECT_ROOT"); projectRoot == "" {
		projectRoot = filepath.Join(exeDir, "..", "..")
	}
	// Ensure /tmp/videos directory exists
	mediaPath := filepath.Join(projectRoot, "tmp", fileType)
	if _, err := os.Stat(mediaPath); os.IsNotExist(err) {
		err := os.MkdirAll(mediaPath, 0755)
		if err != nil {
			log.Panicf("Failed to create directory: %v", err)
		}
	}

	return mediaPath
}

func CreateAndOpenFile(filename string) *os.File {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	return file
}
