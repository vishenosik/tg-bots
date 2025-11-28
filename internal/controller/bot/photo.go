package bot

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vishenosik/gocherry/pkg/bot"
	"github.com/vishenosik/gocherry/pkg/logs"
	"github.com/vishenosik/tg-bots/internal/entity"
)

type PhotoProvider interface {
	SavePhotoWithInfo(r io.Reader, info entity.SenderInfo) (id string, err error)
}

type photoApi struct {
	provider PhotoProvider
	log      *slog.Logger
}

func NewPhotoApi(provider PhotoProvider) *photoApi {
	return &photoApi{
		provider: provider,
		log:      logs.SetupLogger().With(logs.AppComponent("photo_bot")),
	}
}

func (p *photoApi) Route(tb *bot.TelegramBot) {
	tb.Route("save", p.savePhoto())
}

// handleSavePhoto downloads the photo file and saves sender info
func (p *photoApi) savePhoto() bot.HandlerFunc {

	return func(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {

		errorMsg := func(err error) {

			if err != nil {
				p.log.Error("failed to save photo", logs.Error(err))
			}

			resp := tgbotapi.NewMessage(msg.Chat.ID, "❌ Error while saving your photo, try again later.")
			resp.ReplyToMessageID = msg.MessageID
			_, _ = bot.Send(resp)
		}

		if msg.IsCommand() {
			text := "Send a *photo* with caption `/save` and I will save it and your info."
			resp := tgbotapi.NewMessage(msg.Chat.ID, text)
			resp.ParseMode = "Markdown"
			_, _ = bot.Send(resp)
			return
		}

		// 1. choose the highest resolution photo (last element)
		photos := msg.Photo
		photo := photos[len(photos)-1]

		file, err := bot.GetFile(tgbotapi.FileConfig{FileID: photo.FileID})
		if err != nil {
			errorMsg(err)
			return
		}

		fileURL, err := bot.GetFileDirectURL(file.FileID)
		if err != nil {
			errorMsg(err)
			return
		}

		// 2. download the file
		resp, err := http.Get(fileURL)
		if err != nil {
			errorMsg(err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			errorMsg(fmt.Errorf("bad status from telegram: %s", resp.Status))
			return
		}

		if msg.From == nil {
			errorMsg(fmt.Errorf("no sender info in message"))
			return
		}
		from := msg.From

		id, err := p.provider.SavePhotoWithInfo(resp.Body, entity.SenderInfo{
			UserID:    from.ID,
			UserName:  from.UserName,
			FirstName: from.FirstName,
			LastName:  from.LastName,
			ChatID:    msg.Chat.ID,
		})
		if err != nil {
			errorMsg(err)
			return
		}

		p.log.Info("photo saved",
			slog.String("id", id),
			slog.Int64("chat_id", msg.Chat.ID),
		)

		success := tgbotapi.NewMessage(msg.Chat.ID, "✅ Photo and sender info saved.")
		success.ReplyToMessageID = msg.MessageID
		_, _ = bot.Send(success)
	}
}
