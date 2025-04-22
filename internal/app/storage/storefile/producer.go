package storefile

import (
	"encoding/json"
	"github.com/carinfinin/shortener-url/internal/app/logger"
	"github.com/carinfinin/shortener-url/internal/app/models"
	"io"
	"os"
)

// Producer кодирует в json и записывает в файл.
type Producer struct {
	file    *os.File
	encoder *json.Encoder
}

// NewProducer конструктор для Producer.
func NewProducer(path string) (*Producer, error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &Producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

// WriteLine кодирует models.AuthLine в строку.
func (p *Producer) WriteLine(line *models.AuthLine) error {
	return p.encoder.Encode(line)
}

// Close завершает работу с Producer.
// 1. Перемещает указатель в начало файла
// 2. Очищает содержимое файла (truncate)
// 3. Записывает все данные из переданной map в файл
// 4. Синхронизирует изменения на диск (fsync)
// 5. Закрывает файл
func (p *Producer) Close(data map[string]models.AuthLine) error {

	_, err := p.file.Seek(0, io.SeekStart)
	if err != nil {
		logger.Log.Error("Close file Seek error: ", err)
	}
	err = p.file.Truncate(0)
	if err != nil {
		logger.Log.Error("Close file Truncate error: ", err)
	}
	for _, v := range data {
		err := p.WriteLine(&v)
		if err != nil {
			logger.Log.Error("store file close error: ", err)
		}
	}

	err = p.file.Sync()
	if err != nil {
		return err
	}
	return p.file.Close()
}
