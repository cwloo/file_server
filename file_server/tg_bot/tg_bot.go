package tg_bot

import (
	"strings"
	"sync"

	"github.com/cwloo/gonet/logs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	tgBot *TgBotApi
)

// <summary>
// TgBotApi
// <summary>
type TgBotApi struct {
	TgBot_Token  string
	TgBot_ChatId int64
	BotApi       *tgbotapi.BotAPI
	lock         *sync.RWMutex
}

func NewTgBot(TgBot_Token string, TgBot_ChatId int64, useTgBot bool) {
	switch tgBot {
	case nil:
		tgBot = newTgBot(TgBot_Token, TgBot_ChatId, useTgBot)
	default:
		tgBot.update(TgBot_Token, TgBot_ChatId, useTgBot)
	}
}

func newTgBot(TgBot_Token string, TgBot_ChatId int64, useTgBot bool) *TgBotApi {
	s := &TgBotApi{
		TgBot_Token:  TgBot_Token,
		TgBot_ChatId: TgBot_ChatId,
		lock:         &sync.RWMutex{},
	}
	switch useTgBot {
	case true:
		if TgBot_Token == "" || TgBot_ChatId == 0 {
			return s
		}
		s.newBotApi(TgBot_Token)
	default:
		s.resetBotApi()
	}
	return s
}

func (s *TgBotApi) update(TgBot_Token string, TgBot_ChatId int64, useTgBot bool) {
	switch useTgBot {
	case true:
		if TgBot_Token == "" || TgBot_ChatId == 0 {
			s.resetBotApi()
			return
		}
		s.newBotApi(TgBot_Token)
	default:
		s.resetBotApi()
	}
}

func (s *TgBotApi) resetBotApi() {
	s.lock.Lock()
	s.BotApi = nil
	s.lock.Unlock()
}

func (s *TgBotApi) newBotApi(TgBot_Token string) {
	s.lock.Lock()
	if s.BotApi == nil || s.TgBot_Token != TgBot_Token {
		botApi, err := tgbotapi.NewBotAPI(TgBot_Token)
		if err != nil {
			logs.LogFatal(err.Error())
		}
		s.TgBot_Token = TgBot_Token
		s.BotApi = botApi
	}
	s.lock.Unlock()
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
	s.lock.RLock()
	if s.BotApi == nil {
		return
	}
	for _, msg := range msgs {
		smsg := tgbotapi.NewMessage(s.TgBot_ChatId, strings.Join([]string{alert, msg}, ""))
		s.BotApi.Send(smsg)
	}
	s.lock.RUnlock()
}
