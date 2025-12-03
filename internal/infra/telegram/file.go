package telegram

import (
	"fmt"
	"io"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetFile(bot *tgbotapi.BotAPI, id string) (io.Reader, error) {
	file, err := bot.GetFile(tgbotapi.FileConfig{FileID: id})
	if err != nil {
		return nil, err
	}

	fileURL, err := bot.GetFileDirectURL(file.FileID)
	if err != nil {
		return nil, err
	}

	// 2. download the file
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
