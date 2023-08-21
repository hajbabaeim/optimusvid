package main

import (
	"flag"
	"fmt"
	"log"
	"optimusvid/pkg/bot"
	"optimusvid/pkg/system"
	"os"
	"os/exec"
	"path/filepath"

	bt "github.com/SakoDroid/telego"
)

func convertVideo(inputPath string, outputPath string, videoCodec string, videoBitrate string, audioCodec string, audioBitrate string) error {
	cmd := exec.Command("ffmpeg", "-i", inputPath, "-c:v", videoCodec, "-b:v", videoBitrate, "-c:a", audioCodec, "-b:a", audioBitrate, outputPath)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to convert video: %v", err)
	}

	return nil
}

func extractAudio(inputPath string, outputPath string, audioCodec string, audioBitrate string) (*os.File, error) {
	cmd := exec.Command("ffmpeg", "-i", inputPath, "-c:a", audioCodec, "-b:a", audioBitrate, outputPath)

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	audioFile, err := os.Open(outputPath)
	if err != nil {
		return nil, err
	}

	return audioFile, nil
}

func sendAudioToUser(bot *bt.Bot, chatID int, replyTo int, audioFile *os.File, deleteAfter bool) {
	fmt.Printf("🍑🍑🍑 The chatId: %d, replyTo: %d\n\n", chatID, replyTo)
	mediaSender := bot.SendAudio(chatID, replyTo, "Test Caption", "")
	audioMsg, err := mediaSender.SendByFile(audioFile, true, false)
	if err != nil {
		log.Printf("Failed to send audio: %v\n\n", err)
		return
	}

	log.Printf("Sent audio message to chat ID %d with message ID %d\n\n", chatID, audioMsg.Result.MessageId)

	// Optionally, delete the audio file from disk after sending
	if deleteAfter {
		audioPath := audioFile.Name()
		err = os.Remove(audioPath)
		if err != nil {
			log.Printf("Failed to delete audio file: %v\n\n", err)
		}
	}
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
		projectRoot = filepath.Join(exeDir, "..", "..")
	}
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
			originalFilename := filepath.Join(videoDirectory, video.FileId+".mp4")
			outputVideoFilename := filepath.Join(videoDirectory, video.FileId+"_converted_video.mp4")
			outputAduioFilename := filepath.Join(videoDirectory, video.FileId+"_converted_Audio.mp3")
			fmt.Println("-2**>", videoDirectory, originalFilename)
			originalFile := createAndOpenFile(originalFilename)

			defer originalFile.Close()

			_, err := bot.GetFile(video.FileId, true, originalFile)
			if err != nil {
				log.Println("Error while getting the file:", err)
				continue
			}
			// videoFile := createAndOpenFile(outputVideoFilename)
			// defer videoFile.Close()
			convertVideo(originalFilename, outputVideoFilename, "libxvid", "1M", "mp3", "192k")
			audioFile, _ := extractAudio(originalFilename, outputAduioFilename, "mp3", "192k")
			defer audioFile.Close()

			// var videoInReply int
			fmt.Printf("**1**replyToMsg: %+v\n\n%+v\n\n", update.Message.MessageId, update.Message.Chat.Id)

			sendAudioToUser(bot, update.Message.Chat.Id, update.Message.MessageId, audioFile, false)

			// Extract metadata using ffprobe (part of ffmpeg toolset)
			cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", originalFilename)
			// output, err := cmd.Output()
			_, err = cmd.Output()
			if err != nil {
				log.Println("Error while getting metadata:", err)
				continue
			}

			// Respond back to the user with the filename and metadata
			// response := fmt.Sprintf("Filename: %s\nMetadata:\n%s", originalFilename, string(output))
			// msg := telego.Message{MessageID: update.Message.Chat.Id, Text: response}

			// bot.SendMessage(update.Message.Chat.Id, originalFilename, "", 1, true, false)
		}
	}

}

func main() {
	system.LoadEnv()
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
