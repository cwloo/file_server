package main

import (
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/uploader/file_server/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewTgBot(token string) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		logs.LogFatal(err.Error())
	}
	return bot
}

func SendTgBotMsg(msgs ...string) {
	for _, msg := range msgs {
		smsg := tgbotapi.NewMessage(config.Config.TgBot_ChatId, msg)
		TgBot.Send(smsg)
	}
}
