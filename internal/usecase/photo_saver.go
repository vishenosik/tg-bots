package usecase

import (
	"context"
	"io"

	"github.com/google/uuid"
	"github.com/vishenosik/tg-bots/internal/entity"
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

func (pu *PhotoUsecase) SavePhotoWithInfo(r io.Reader, info entity.SenderInfo) (id string, err error) {
	id, err = pu.file.SaveFile(context.Background(), r, uuid.NewString()+".jpg")
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
