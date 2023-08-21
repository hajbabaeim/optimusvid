package optimus

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	bt "github.com/SakoDroid/telego"
	cfg "github.com/SakoDroid/telego/configs"
)

type Optimus struct {
	Bot     *bt.Bot
	Format  string
	Quality string
}

func Init() *Optimus {
	up := cfg.DefaultUpdateConfigs()

	cf := cfg.BotConfigs{BotAPI: cfg.DefaultBotAPI, APIKey: os.Getenv("TELEGRAM_API"), UpdateConfigs: up, Webhook: false, LogFileAddress: cfg.DefaultLogFile}

	bot, err := bt.NewBot(&cf)
	if err != nil {
		fmt.Printf("Bot not initialised due to this issue: %v", err)
	}

	return &Optimus{Bot: bot, Format: "mp3", Quality: "192k"}
}

func (optimus *Optimus) ConvertVideoToAudio(inputPath string, outputPath string, videoCodec string, videoBitrate string, audioCodec string, audioBitrate string) error {
	cmd := exec.Command("ffmpeg", "-i", inputPath, "-c:v", videoCodec, "-b:v", videoBitrate, "-c:a", audioCodec, "-b:a", audioBitrate, outputPath)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to convert video: %v", err)
	}

	return nil
}

func (optimus *Optimus) ExtractAudioFromVideo(inputPath string, outputPath string, audioCodec string, audioBitrate string) (*os.File, error) {
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

func (optimus *Optimus) SendAudioToUser(chatID int, replyTo int, audioFile *os.File, deleteAfter bool) {
	mediaSender := optimus.Bot.SendAudio(chatID, replyTo, "Test Caption", "")
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
