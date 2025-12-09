package bot

import (
	"fmt"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vishenosik/gocherry/pkg/logs"

	bot "github.com/vishenosik/tg-bot-engine"
)

type affirmationsApi struct {
	log *slog.Logger
}

func NewAffirmationsApi() *affirmationsApi {
	return &affirmationsApi{
		log: logs.SetupLogger().With(logs.AppComponent("affirmations_bot")),
	}
}

func (a *affirmationsApi) Route(tb *bot.TelegramBot) {

	tb.Route("send_affirmations_daily", a.sendAffirmationsDailyHandler())
}

// handleSavePhoto downloads the photo file and saves sender info
func (a *affirmationsApi) sendAffirmationsDailyHandler() bot.HandlerFunc {

	return func(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {

		greet := "Привет, сладкие булочки, вот вам аффирмация дня"
		affirmation := "У вас все получится ❤️"

		success := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("%s: %s", greet, affirmation))
		success.ReplyToMessageID = msg.MessageID
		_, _ = bot.Send(success)

		sticker := tgbotapi.NewSticker(msg.Chat.ID, tgbotapi.FileID("CAACAgIAAxkBAAMyaSa9ePcR-V0xNf6oMCkaesSUmaAAAnmWAAIwYTBJohBtEH8TlgI2BA"))
		_, _ = bot.Send(sticker)
	}
}
