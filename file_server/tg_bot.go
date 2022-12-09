package main

import (
	"strings"

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

func TgWarnMsg(msgs ...string) {
	tgBotMsg("⚠️", msgs...)
}

func TgSuccMsg(msgs ...string) {
	tgBotMsg("✅", msgs...)
}

func TgErrMsg(msgs ...string) {
	tgBotMsg("❌", msgs...)
}

func tgBotMsg(alert string, msgs ...string) {
	for _, msg := range msgs {
		smsg := tgbotapi.NewMessage(config.Config.TgBot_ChatId, strings.Join([]string{alert, msg}, ""))
		TgBot.Send(smsg)
	}
}
