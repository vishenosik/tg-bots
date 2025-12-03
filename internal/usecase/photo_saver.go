package usecase

import (
	"context"
	"io"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/vishenosik/tg-bots/internal/entity"
	"github.com/vishenosik/tg-bots/internal/infra/telegram"
)

type FileProvider interface {
	SaveFile(ctx context.Context, file io.Reader, fileName string) (string, error)
}

type InfoProvider interface {
	SaveSenderInfo(entity.SenderInfo) error
}

type PhotoUsecase struct {
	file FileProvider
	info InfoProvider
}

func NewPhotoUsecase(
	file FileProvider,
	info InfoProvider,
) *PhotoUsecase {
	return &PhotoUsecase{
		file: file,
		info: info,
	}
}

func (pu *PhotoUsecase) SavePhotoWithInfo(
	bot *tgbotapi.BotAPI,
	fileID string,
	info entity.SenderInfo,
) (string, error) {

	r, err := telegram.GetFile(bot, fileID)
	if err != nil {
		return "", err
	}

	id, err := pu.file.SaveFile(context.Background(), r, uuid.NewString()+".jpg")
	if err != nil {
		return "", err
	}

	info.PhotoID = id
	err = pu.info.SaveSenderInfo(info)
	if err != nil {
		return "", err
	}
	return id, nil
}
