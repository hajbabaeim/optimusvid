package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"

	"github.com/joho/godotenv"
)

const projectDirName = "bot" // change to relevant project name

func loadEnv() {
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))

	err := godotenv.Load(string(rootPath) + `/.env`)

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func convertVideo(inputPath string, outputPath string, videoCodec string, videoBitrate string, audioCodec string, audioBitrate string) error {
	cmd := exec.Command("ffmpeg", "-i", inputPath, "-c:v", videoCodec, "-b:v", videoBitrate, "-c:a", audioCodec, "-b:a", audioBitrate, outputPath)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to convert video: %v", err)
	}

	return nil
}

func extractAudio(inputPath string, outputPath string, audioCodec string, audioBitrate string) error {
	cmd := exec.Command("ffmpeg", "-i", inputPath, "-c:a", audioCodec, "-b:a", audioBitrate, outputPath)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to extract audio: %v", err)
	}

	return nil
}

func main() {
	loadEnv()
	// Define and parse flags
	input := flag.String("input", "sample.mp4", "Path to the input video file")
	outputVideo := flag.String("outputVideo", "output.avi", "Path for the converted video file")
	outputAudio := flag.String("outputAudio", "output.mp3", "Path for the extracted audio file")
	videoCodec := flag.String("videoCodec", "libxvid", "Video codec to be used")
	videoBitrate := flag.String("videoBitrate", "1M", "Video bitrate")
	audioCodec := flag.String("audioCodec", "mp3", "Audio codec to be used")
	audioBitrate := flag.String("audioBitrate", "192k", "Audio bitrate")

	flag.Parse()

	err := convertVideo(*input, *outputVideo, *videoCodec, *videoBitrate, *audioCodec, *audioBitrate)
	if err != nil {
		log.Fatalf("Error converting video: %v", err)
	}

	err = extractAudio(*input, *outputAudio, *audioCodec, *audioBitrate)
	if err != nil {
		log.Fatalf("Error extracting audio: %v", err)
	}

	fmt.Println("Conversion successful!")
}
