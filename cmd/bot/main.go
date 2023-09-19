package main

import (
	"fmt"
	// "log"
	"optimusvid/pkg/optimus"
	"optimusvid/pkg/system"
	"path/filepath"

	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	bt "github.com/SakoDroid/telego"
	objs "github.com/SakoDroid/telego/objects"
)

const (
	maxDurationSeconds = 1 * 60
)

// A map to store user states
var userStates = make(map[int]string)

// Check if the user is currently expected to make a quality selection
func isAwaitingQualitySelection(optimus *optimus.Optimus, chatID int) bool {
	state, exists := userStates[chatID]
	return exists && state == "awaiting_quality_selection"
}

// Set the user's state
func setUserState(optimus *optimus.Optimus, chatID int, state string) {
	userStates[chatID] = state
}

// Clear the user's state
func clearUserState(optimus *optimus.Optimus, chatID int) {
	delete(userStates, chatID)
}

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

		fmt.Println("-----------------------the output value: ", update.Message.Text)
		switch update.Message.Text {
		case "/start":
			sendWelcomeMessage(optimus, update.Message.Chat.Id)
		// case "aac", "flac", "Vorbis":
		// 	handleFormatSelection(optimus, update.Message)
		case "mp3":
			optimus.Format = "mp3"
			createQualityKeyboard(optimus, update.Message)
			// handleFormatSelection(optimus, update.Message)
		case "flac":
			optimus.Format = "flac"
			optimus.Quality = ""
		case "32k", "64k", "96k", "128k", "192k", "256k", "320k":
			// If you're awaiting a quality selection, handle it here
			if isAwaitingQualitySelection(optimus, update.Message.Chat.Id) {
				go handleQualitySelection(optimus, update.Message)
			}
		case "/about":
			sendBotDescription(optimus, update.Message.Chat.Id, update.Message.MessageId)
		default:
			go handleVideoConversion(optimus, update.Message)
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

func handleQualitySelection(optimus *optimus.Optimus, message *objs.Message) {
	// Process the user's response
	optimus.Quality = message.Text
	fmt.Println("----> [][][] the message.Text here is : ", message.Text)

	// Hide the keyboard and send a confirmation message
	_, err := optimus.Bot.SendMessage(message.Chat.Id, "Quality selected: "+message.Text, "", message.MessageId, false, false)
	if err != nil {
		fmt.Println(err)
	}

	// Clear the user's state
	clearUserState(optimus, message.Chat.Id)
}

func createQualityKeyboard(optimus *optimus.Optimus, message *objs.Message) {
	qualityKb := optimus.Bot.CreateKeyboard(true, true, false, "Choose quality of audio file.")
	qualityKb.AddButton("64k", 1)
	qualityKb.AddButton("96k", 1)
	qualityKb.AddButton("128k", 2)
	qualityKb.AddButton("192k", 2)
	fmt.Println("----> [2][2][2] the message.Text here is : ", message.Text)
	_, err := optimus.Bot.AdvancedMode().ASendMessage(message.Chat.Id, "Please choose a quality for audio file.", "", message.MessageId, false, false, nil, false, false, qualityKb)
	if err != nil {
		fmt.Println(err)
	}
	setUserState(optimus, message.Chat.Id, "awaiting_quality_selection")
}

func handleFormatSelection(optimus *optimus.Optimus, message *objs.Message) {
	// Check if the user has already made a selection
	optimus.Quality = message.Text
	fmt.Println("----> [1][1][1] the message.Text here is : ", message.Text)

	// After receiving the value, hide the keyboard
	_, err := optimus.Bot.SendMessage(message.Chat.Id, "Quality selected: "+message.Text, "", message.MessageId, false, false)
	if err != nil {
		fmt.Println(err)
	}

	// Clear the user's state or context
	clearUserState(optimus, message.Chat.Id)
}

func sendBotDescription(optimus *optimus.Optimus, chatID int, messageID int) {
	description := optimus.GetBotDescription()
	_, err := optimus.Bot.SendMessage(chatID, description, "", messageID, false, false)
	if err != nil {
		log.Printf("Failed to send the video length warning: %v", err)
	}
}

func handleVideoConversion(optimus *optimus.Optimus, message *objs.Message) error {

	video := message.Video

	if video.Duration > maxDurationSeconds {
		durationLimitMsg := fmt.Sprintf("The uploaded video exceeds the %d-seconds limit. Please upload a shorter video.", maxDurationSeconds)
		return fmt.Errorf("failed to send the video length warning: %s", durationLimitMsg)
	}

	videoDirectory := system.EnsureVideoDirectory()
	originalFilename := filepath.Join(videoDirectory, video.FileId+".mp4")

	log.WithFields(log.Fields{
		"videoDirectory":   videoDirectory,
		"originalFilename": originalFilename,
		"optimus.Format":   optimus.Format,
	}).Info("Video conversion started")

	outputAudioFilename := filepath.Join(videoDirectory, video.FileId+fmt.Sprintf("_converted_Audio.%s", optimus.Format))
	originalFile := system.CreateAndOpenFile(originalFilename)
	defer originalFile.Close()

	_, err := optimus.Bot.GetFile(video.FileId, true, originalFile)
	if err != nil {
		return fmt.Errorf("error while getting the file: %v", err)
	}
	fmt.Printf("🍑 originalFilename: %s\n 🍑outputAudioFilename: %s\n 🍑optimus.Format: %s\n🍑 optimus.Quality: %s\n", originalFilename, outputAudioFilename, optimus.Format, optimus.Quality)

	audioFile, err := optimus.ExtractAudioFromVideo(originalFilename, outputAudioFilename, optimus.Format, optimus.Quality)
	if err != nil {
		log.Printf("Failed to extract audio from video: %v", err)
		return err
	}
	defer audioFile.Close()

	optimus.SendAudioToUser(message.Chat.Id, message.MessageId, audioFile, false)

	return nil
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
