package optimus

import (
	"bytes"
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
	// Create default update configs
	up := cfg.DefaultUpdateConfigs()

	// Create bot configs with default values and environment variable for API key
	cf := cfg.BotConfigs{
		BotAPI:         cfg.DefaultBotAPI,
		APIKey:         os.Getenv("TELEGRAM_API"),
		UpdateConfigs:  up,
		Webhook:        false,
		LogFileAddress: cfg.DefaultLogFile,
	}

	// Create a new bot instance
	bot, err := bt.NewBot(&cf)
	if err != nil {
		fmt.Printf("Bot not initialized due to this issue: %v", err)
	}

	// Initialize Optimus struct with bot, format, and quality
	return &Optimus{
		Bot:     bot,
		Format:  "mp3",
		Quality: "128k",
	}
}

func (optimus *Optimus) ExtractAudioFromVideo(inputPath string, outputPath string, audioCodec string, audioBitrate string) (*os.File, error) {
	fmt.Printf("🚀 audioCodec: %s\n 🚀audioBitrate: %s\n 🚀inputPath: %s\n🚀 outputPath: %s\n", audioCodec, audioBitrate, inputPath, outputPath)
	var cmd *exec.Cmd

	switch optimus.Format {
	case "flac":
		fmt.Println("Changed output format to FLAC.")
		cmd = exec.Command("ffmpeg", "-y", "-i", inputPath, "-vn", "-c:a", "flac", outputPath)
	case "mp3":
		fmt.Println("Changed output format to MP3.")
		cmd = exec.Command("ffmpeg", "-y", "-i", inputPath, "-vn", "-c:a", "libmp3lame", "-b:a", audioBitrate, outputPath)
	case "wav":
		fmt.Println("Changed output format to WAV.")
		cmd = exec.Command("ffmpeg", "-y", "-i", inputPath, "-vn", "-c:a", "pcm_s16le", outputPath)
	default:
		return nil, fmt.Errorf("Unsupported format: %s", optimus.Format)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		fmt.Println(fmt.Sprintf("Failed to execute command: %v\nOutput:\n%s\nError:\n%s", cmd.Args, stdout.String(), stderr.String()))
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

func (optimus *Optimus) GetBotDescription() string {
	return `OptimusVidBot is your premier video-to-audio converter bot. Whether you're looking to transform music videos into MP3s or capture the audio from informative content, OptimusVidBot offers a broad range of output formats including AAC, MP3, Vorbis, FLAC, and WAV. But that's not all! You're in full control of the audio quality, allowing you to choose bitrates from as low as 32k (great for spoken content) to as high as 320k (perfect for high-fidelity music). Why compromise when you can get the best? Dive in now and experience the richness of sound with OptimusVidBot.`
}
