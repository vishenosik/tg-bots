package telegram

import (
	"fmt"
	"io"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vishenosik/gocherry/pkg/bot"
)

type TgApi struct {
	bot *tgbotapi.BotAPI
}

func NewTgApiEnv() (*TgApi, error) {
	tgBot, err := bot.NewBotAPIEnv()
	if err != nil {
		return nil, err
	}

	return &TgApi{bot: tgBot.Bot()}, nil
}

func (tg *TgApi) GetFile(id string) (io.Reader, error) {
	file, err := tg.bot.GetFile(tgbotapi.FileConfig{FileID: id})
	if err != nil {
		return nil, err
	}

	fileURL, err := tg.bot.GetFileDirectURL(file.FileID)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status from telegram: %s", resp.Status)
	}
	return resp.Body, nil
}
