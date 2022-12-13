package tg_bot

import (
	"strings"

	"github.com/cwloo/gonet/logs"
	// "github.com/cwloo/uploader/file_server/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	tgBot *TgBotApi
)

// <summary>
// TgBotApi
// <summary>
type TgBotApi struct {
	TgBot_ChatId int64
	BotApi       *tgbotapi.BotAPI
}

func NewTgBot(TgBot_Token string, TgBot_ChatId int64) {
	tgBot = newTgBot(TgBot_Token, TgBot_ChatId)
}

func newTgBot(TgBot_Token string, TgBot_ChatId int64) *TgBotApi {
	s := &TgBotApi{
		TgBot_ChatId: TgBot_ChatId,
	}
	botApi, err := tgbotapi.NewBotAPI(TgBot_Token)
	if err != nil {
		logs.LogFatal(err.Error())
	}
	s.BotApi = botApi
	return s
}

func TgWarnMsg(msgs ...string) {
	if tgBot == nil {
		return
	}
	tgBot.tgBotMsg("⚠️", msgs...)
}

func TgSuccMsg(msgs ...string) {
	if tgBot == nil {
		return
	}
	tgBot.tgBotMsg("✅", msgs...)
}

func TgErrMsg(msgs ...string) {
	if tgBot == nil {
		return
	}
	tgBot.tgBotMsg("❌", msgs...)
}

func (s *TgBotApi) tgBotMsg(alert string, msgs ...string) {
	if s.BotApi == nil {
		return
	}
	for _, msg := range msgs {
		smsg := tgbotapi.NewMessage(s.TgBot_ChatId, strings.Join([]string{alert, msg}, ""))
		s.BotApi.Send(smsg)
	}
}
