package bot

import (
	"fmt"
	"os"

	bt "github.com/SakoDroid/telego"
	cfg "github.com/SakoDroid/telego/configs"
)

type OptimsuVidBot struct {
	*bt.Bot
}

func Init() *bt.Bot {
	up := cfg.DefaultUpdateConfigs()

	cf := cfg.BotConfigs{BotAPI: cfg.DefaultBotAPI, APIKey: os.Getenv("TELEGRAM_API"), UpdateConfigs: up, Webhook: false, LogFileAddress: cfg.DefaultLogFile}

	bot, err := bt.NewBot(&cf)
	if err != nil {
		fmt.Printf("Bot not initialised due to this issue: %v", err)
	}

	return bot
}

// func (b *bt.Bot) SendMessage(msg string) {

// }

// func (bot *OptimsuVidBot) GetFilename() string {
// 	updateChannel := bot.GetUpdateChannel()
// 	update := <-*updateChannel

// 	// Get sticker file id
// 	fi := update.Message.Sticker.FileId

// 	// Open a file in the computer.
// 	fl, _ := os.OpenFile("sticker.webp", os.O_CREATE|os.O_WRONLY, 0666)

// 	// Gets the file info and downloads it.
// 	_, err := bot.GetFile(fi, true, fl)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fl.Close()
// }
