package main

import (
	"flag"
	"fmt"
	"log"
	"optimusvid/pkg/bot"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	bt "github.com/SakoDroid/telego"
	"github.com/joho/godotenv"
)

const projectDirName = "OptimusVid-Convert" // change to relevant project name

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

func createAndOpenFile(filename string) *os.File {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	return file
}

func ensureVideoDirectory() string {
	exe, err := os.Executable()
	if err != nil {
		log.Panicf("Failed to determine the executable path: %v", err)
	}

	// Get the directory of the executable
	exeDir := filepath.Dir(exe)

	// Traverse two directories up to reach the project root
	var projectRoot string
	if projectRoot = os.Getenv("PROJECT_ROOT"); projectRoot == "" {
		fmt.Printf("projectRoot -**-: %s\n\n", projectRoot)
		projectRoot = filepath.Join(exeDir, "..", "..")
	}
	fmt.Printf("projectRoot: %s\n\n", projectRoot)
	// Ensure /tmp/videos directory exists
	videoPath := filepath.Join(projectRoot, "tmp", "videos")
	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		err := os.MkdirAll(videoPath, 0755)
		if err != nil {
			log.Panicf("Failed to create directory: %v", err)
		}
	}

	return videoPath
}

func start(bot *bt.Bot) {

	updates := bot.GetUpdateChannel()

	for update := range *updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Video != nil {

			video := update.Message.Video
			fmt.Println("-1-->", video.Duration, video.FileId)

			videoDirectory := ensureVideoDirectory()
			filename := filepath.Join(videoDirectory, video.FileId+".mp4")
			outputFilename := filepath.Join(videoDirectory, video.FileId+"_converted.mp4")
			fmt.Println("-2**>", videoDirectory, filename)
			file := createAndOpenFile(filename)
			defer file.Close()

			videoFile, err := bot.GetFile(video.FileId, true, file)
			if err != nil {
				log.Println("Error while getting the file:", err)
				continue
			}
			fmt.Println("-2-->", videoFile)
			convertVideo(filename, outputFilename, "libxvid", "1M", "mp3", "192k")

			// Extract metadata using ffprobe (part of ffmpeg toolset)
			cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", filename)
			output, err := cmd.Output()
			if err != nil {
				log.Println("Error while getting metadata:", err)
				continue
			}

			// Respond back to the user with the filename and metadata
			response := fmt.Sprintf("Filename: %s\nMetadata:\n%s", filename, string(output))
			fmt.Println("-4-->", response[0])
			// msg := telego.Message{MessageID: update.Message.Chat.Id, Text: response}
			bot.SendMessage(update.Message.Chat.Id, filename, "", 1, true, false)
		}
	}

}

func main() {
	loadEnv()
	// Define and parse flags
	// input := flag.String("input", "sample.mp4", "Path to the input video file")
	// outputVideo := flag.String("outputVideo", "output.avi", "Path for the converted video file")
	// outputAudio := flag.String("outputAudio", "output.mp3", "Path for the extracted audio file")
	// videoCodec := flag.String("videoCodec", "libxvid", "Video codec to be used")
	// videoBitrate := flag.String("videoBitrate", "1M", "Video bitrate")
	// audioCodec := flag.String("audioCodec", "mp3", "Audio codec to be used")
	// audioBitrate := flag.String("audioBitrate", "192k", "Audio bitrate")

	flag.Parse()

	bot := bot.Init()

	if bot.Run() == nil {
		go start(bot)
	}
	select {}

	// err = convertVideo(*input, *outputVideo, *videoCodec, *videoBitrate, *audioCodec, *audioBitrate)
	// if err != nil {
	// 	log.Fatalf("Error converting video: %v", err)
	// }

	// err = extractAudio(*input, *outputAudio, *audioCodec, *audioBitrate)
	// if err != nil {
	// 	log.Fatalf("Error extracting audio: %v", err)
	// }

	// fmt.Println("Conversion successful!")
}
