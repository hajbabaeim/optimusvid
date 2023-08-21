package main

import (
	"flag"
	"log"
	"optimusvid/pkg/optimus"
	"optimusvid/pkg/system"
	"os/exec"
	"path/filepath"

	bt "github.com/SakoDroid/telego"
)

func start(bot *bt.Bot) {

	updates := bot.GetUpdateChannel()

	for update := range *updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Text == "/start" {
			_, err := bot.SendMessage(update.Message.Chat.Id, "Welcome to OptimusVid Convert Bot! Send me a video and I'll convert it for you.", "", 1, true, false)
			if err != nil {
				log.Printf("Failed to send the welcome message: %v", err)
			}
		}

		if update.Message.Video != nil {
			video := update.Message.Video
			videoDirectory := system.EnsureVideoDirectory()
			originalFilename := filepath.Join(videoDirectory, video.FileId+".mp4")
			outputVideoFilename := filepath.Join(videoDirectory, video.FileId+"_converted_video.mp4")
			outputAduioFilename := filepath.Join(videoDirectory, video.FileId+"_converted_Audio.mp3")
			originalFile := system.CreateAndOpenFile(originalFilename)
			defer originalFile.Close()

			_, err := bot.GetFile(video.FileId, true, originalFile)
			if err != nil {
				log.Println("Error while getting the file:", err)
				continue
			}
			// videoFile := CreateAndOpenFile(outputVideoFilename)
			// defer videoFile.Close()
			optimus.ConvertVideoToAudio(originalFilename, outputVideoFilename, "libxvid", "1M", "mp3", "192k")
			audioFile, _ := optimus.ExtractAudioFromVideo(originalFilename, outputAduioFilename, "mp3", "192k")
			defer audioFile.Close()

			optimus.SendAudioToUser(bot, update.Message.Chat.Id, update.Message.MessageId, audioFile, false)

			// Extract metadata using ffprobe (part of ffmpeg toolset)
			cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", originalFilename)
			// output, err := cmd.Output()
			_, err = cmd.Output()
			if err != nil {
				log.Println("Error while getting metadata:", err)
				continue
			}
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

	bot := optimus.Init()

	if bot.Run() == nil {
		go start(bot)
	}
	select {}

	// err = ConvertVideoToAudio(*input, *outputVideo, *videoCodec, *videoBitrate, *audioCodec, *audioBitrate)
	// if err != nil {
	// 	log.Fatalf("Error converting video: %v", err)
	// }

	// err = ExtractAudioFromVideo(*input, *outputAudio, *audioCodec, *audioBitrate)
	// if err != nil {
	// 	log.Fatalf("Error extracting audio: %v", err)
	// }

	// fmt.Println("Conversion successful!")
}
