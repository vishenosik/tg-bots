package usecase

import (
	"context"
	"io"

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
	tgApi *telegram.TgApi
	file  FileProvider
	info  InfoProvider
}

func NewPhotoUsecase(
	file FileProvider,
	info InfoProvider,
) *PhotoUsecase {
	tgApi, err := telegram.NewTgApiEnv()
	if err != nil {
		panic(err)
	}
	return &PhotoUsecase{
		tgApi: tgApi,
		file:  file,
		info:  info,
	}
}

func (pu *PhotoUsecase) SavePhotoWithInfo(
	fileID string,
	info entity.SenderInfo,
) (string, error) {

	r, err := pu.tgApi.GetFile(fileID)
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
