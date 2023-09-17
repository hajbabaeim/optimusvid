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
	fmt.Printf("---->>>> %s, %s", audioCodec, audioBitrate)
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

func (optimus *Optimus) GetBotDescription() string {
	return `OptimusVidBot is your premier video-to-audio converter bot. Whether you're looking to transform music videos into MP3s or capture the audio from informative content, OptimusVidBot offers a broad range of output formats including AAC, MP3, Vorbis, FLAC, and WAV. But that's not all! You're in full control of the audio quality, allowing you to choose bitrates from as low as 32k (great for spoken content) to as high as 320k (perfect for high-fidelity music). Why compromise when you can get the best? Dive in now and experience the richness of sound with OptimusVidBot.`
}
