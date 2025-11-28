package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/vishenosik/gocherry/pkg/errors"
	"github.com/vishenosik/tg-bots/internal/entity"
)

func init() {

}

type CSVService struct {
	path string
	once sync.Once
}

func NewCSVStore(path string) *CSVService {
	cs := &CSVService{
		path: path,
	}
	return cs
}

func (cs *CSVService) SaveSenderInfo(info entity.SenderInfo) error {
	return cs.write(
		time.Now().Format(time.RFC3339),
		fmt.Sprint(info.UserID),
		info.UserName,
		info.FirstName,
		info.LastName,
		fmt.Sprint(info.ChatID),
		info.PhotoID,
	)
}

// createSendersCSVIfNotExists creates CSV with header if not exists
func (cs *CSVService) createCSVIfNotExists() error {

	if _, err := os.Stat(cs.path); err == nil {
		return err
	}

	f, err := os.Create(cs.path)
	if err != nil {
		return errors.Wrap(err, "cannot create csv file")
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	header := []string{"timestamp", "user_id", "username", "first_name", "last_name", "chat_id", "photo_path"}
	if err := w.Write(header); err != nil {
		return errors.Wrap(err, "cannot write header to csv file")
	}
	return nil
}

func (cs *CSVService) write(records ...string) error {
	cs.once.Do(func() {
		if err := cs.createCSVIfNotExists(); err != nil {
			panic(err)
		}
	})
	f, err := os.OpenFile(cs.path, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write(records); err != nil {
		return err
	}

	return nil
}
