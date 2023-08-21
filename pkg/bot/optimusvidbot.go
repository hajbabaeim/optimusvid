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
