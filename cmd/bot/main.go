package main

import (
	"fmt"
	"log"
	"optimusvid/pkg/optimus"
	"optimusvid/pkg/system"
	"path/filepath"

	"os"
	"os/signal"
	"syscall"

	bt "github.com/SakoDroid/telego"
	objs "github.com/SakoDroid/telego/objects"
)

const (
	maxDurationSeconds = 1 * 60
)

func dismissKeyboard(bot *bt.Bot, chatID int, replyText string) {
	// Create an empty keyboard (effectively removes the previous custom keyboard)
	kb := bot.CreateKeyboard(true, false, false, "") // set the first parameter to 'true' to remove the custom keyboard

	_, err := bot.AdvancedMode().ASendMessage(chatID, replyText, "", 0, false, false, nil, false, false, kb)
	if err != nil {
		fmt.Println(err)
	}
}

func start(optimus *optimus.Optimus) {
	updates := optimus.Bot.GetUpdateChannel()

	optimus.Bot.AddHandler("/settings", func(u *objs.Update) {
		handleSettings(optimus, u)
	}, "private")

	for update := range *updates {
		if update.Message == nil {
			continue
		}

		switch update.Message.Text {
		case "/start":
			sendWelcomeMessage(optimus, update.Message.Chat.Id)
		// case "aac", "flac", "Vorbis":
		// 	handleFormatSelection(optimus, update.Message)
		case "mp3":
			optimus.Format = "mp3"
			handleFormatSelection(optimus, update.Message)
		case "flac":
			optimus.Format = "flac"
			optimus.Quality = ""
		case "32k", "64k", "96k", "128k", "192k", "256k", "320k":
			optimus.Quality = update.Message.Text
		case "/about":
			sendBotDescription(optimus, update.Message.Chat.Id, update.Message.MessageId)
		default:
			handleVideoConversion(optimus, update.Message)
		}
	}
}

func handleSettings(optimus *optimus.Optimus, u *objs.Update) {
	settingsKb := optimus.Bot.CreateKeyboard(false, false, false, "Choose an options for format of output audio file.")
	settingsKb.AddButton("mp3", 1)
	settingsKb.AddButton("flac", 1)
	_, err := optimus.Bot.AdvancedMode().ASendMessage(u.Message.Chat.Id, "Please choose a format", "", u.Message.MessageId, false, false, nil, false, false, settingsKb)
	if err != nil {
		fmt.Println(err)
	}
}

func sendWelcomeMessage(optimus *optimus.Optimus, chatID int) {
	_, err := optimus.Bot.SendMessage(chatID, "Welcome to OptimusVid Convert Bot! Send me a video and I'll convert it for you.", "", 1, true, false)
	if err != nil {
		log.Printf("Failed to send the welcome message: %v", err)
	}
}

func handleFormatSelection(optimus *optimus.Optimus, message *objs.Message) {
	qualityKb := optimus.Bot.CreateKeyboard(true, true, false, "Choose quality of audio file.")
	qualityKb.AddButton("64k", 1)
	qualityKb.AddButton("96k", 1)
	qualityKb.AddButton("128k", 2)
	qualityKb.AddButton("192k", 2)
	_, err := optimus.Bot.AdvancedMode().ASendMessage(message.Chat.Id, "Please choose a quality for audio file.", "", message.MessageId, false, false, nil, false, false, qualityKb)
	if err != nil {
		fmt.Println(err)
	}
	optimus.Format = message.Text
}

func sendBotDescription(optimus *optimus.Optimus, chatID int, messageID int) {
	description := optimus.GetBotDescription()
	_, err := optimus.Bot.SendMessage(chatID, description, "", messageID, false, false)
	if err != nil {
		log.Printf("Failed to send the video length warning: %v", err)
	}
}

func handleVideoConversion(optimus *optimus.Optimus, message *objs.Message) {
	video := message.Video

	if video.Duration > maxDurationSeconds {
		durationLimitMsg := fmt.Sprintf("The uploaded video exceeds the %d-seconds limit. Please upload a shorter video.", maxDurationSeconds)
		_, err := optimus.Bot.SendMessage(message.Chat.Id, durationLimitMsg, "", message.MessageId, false, false)
		if err != nil {
			log.Printf("Failed to send the video length warning: %v", err)
		}
		return
	}

	videoDirectory := system.EnsureVideoDirectory()
	originalFilename := filepath.Join(videoDirectory, video.FileId+".mp4")

	outputAudioFilename := filepath.Join(videoDirectory, video.FileId+fmt.Sprintf("_converted_Audio.%s", optimus.Format))
	originalFile := system.CreateAndOpenFile(originalFilename)
	fmt.Printf("1- videoDirectory: %s\n 2- originalFilename: %s\n 3- outputAudioFilename: %s\n 4- originalFile: %s\n", videoDirectory, originalFilename, outputAudioFilename, originalFile)
	defer originalFile.Close()

	_, err := optimus.Bot.GetFile(video.FileId, true, originalFile)
	if err != nil {
		log.Println("Error while getting the file:", err)
		return
	}

	audioFile, _ := optimus.ExtractAudioFromVideo(originalFilename, outputAudioFilename, optimus.Format, optimus.Quality)
	defer audioFile.Close()

	optimus.SendAudioToUser(message.Chat.Id, message.MessageId, audioFile, false)
}

func main() {
	system.LoadEnv()
	optimus := optimus.Init()
	err := optimus.Bot.Run()
	if err != nil {
		log.Fatalf("Error running bot: %v", err)
	}
	go start(optimus)

	// Create a channel to receive termination signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Block until a termination signal is received
	<-sigCh

	// Perform cleanup and shutdown operations here

	// Gracefully exit the application
	os.Exit(0)
}
