package bot

import (
	"errors"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vishenosik/gocherry/pkg/bot"
	"github.com/vishenosik/gocherry/pkg/logs"
	"github.com/vishenosik/tg-bots/internal/entity"
)

type PhotoProvider interface {
	SavePhotoWithInfo(
		bot *tgbotapi.BotAPI,
		fileID string,
		info entity.SenderInfo,
	) (id string, err error)
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
	tb.Route("save", p.saveHandler())
}

// handleSavePhoto downloads the photo file and saves sender info
func (p *photoApi) saveHandler() bot.HandlerFunc {

	return func(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {

		if msg.IsCommand() {
			text := "Send a *photo* with caption `/save` and I will save it and your info."
			resp := tgbotapi.NewMessage(msg.Chat.ID, text)
			resp.ParseMode = "Markdown"
			_, _ = bot.Send(resp)
			return
		}

		err := p.save(bot, msg)
		if err != nil {
			p.log.Error("failed to save photo", logs.Error(err))
			resp := tgbotapi.NewMessage(msg.Chat.ID, "❌ Error while saving your photo, try again later.")
			resp.ReplyToMessageID = msg.MessageID
			_, _ = bot.Send(resp)
			return
		}

		success := tgbotapi.NewMessage(msg.Chat.ID, "✅ Photo and sender info saved.")
		success.ReplyToMessageID = msg.MessageID
		_, _ = bot.Send(success)
	}
}

func (p *photoApi) save(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) error {

	// choose the highest resolution photo (last element)
	photos := msg.Photo
	photo := photos[len(photos)-1]

	if msg.From == nil {
		return errors.New("no sender info in message")
	}
	from := msg.From

	id, err := p.provider.SavePhotoWithInfo(bot, photo.FileID, entity.SenderInfo{
		UserID:    from.ID,
		UserName:  from.UserName,
		FirstName: from.FirstName,
		LastName:  from.LastName,
		ChatID:    msg.Chat.ID,
	})
	if err != nil {
		return err
	}

	p.log.Info("photo saved",
		slog.String("id", id),
		slog.Int64("chat_id", msg.Chat.ID),
	)

	return nil
}
