package fileservice

import (
	"context"
	"io"

	"github.com/vishenosik/file-svc-sdk/client"
)

type FileService struct {
	client *client.FileServiceClient
}

func NewFileService() (*FileService, error) {

	fsClient, err := client.NewFileServiceClientEnv()
	if err != nil {
		return nil, err
	}

	return &FileService{
		client: fsClient,
	}, nil
}

func (fs *FileService) SaveFile(ctx context.Context, file io.Reader, fileName string) (string, error) {
	resp, err := fs.client.V1().Upload(ctx, file, fileName)
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}
