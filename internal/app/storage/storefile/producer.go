package storefile

import (
	"encoding/json"
	"github.com/carinfinin/shortener-url/internal/app/models"
	"os"
)

type Producer struct {
	file    *os.File
	encoder *json.Encoder
}

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
func (p *Producer) WriteLine(line *models.Line) error {
	return p.encoder.Encode(line)
}

func (p *Producer) Close() error {
	err := p.file.Sync()
	if err != nil {
		return err
	}
	return p.file.Close()
}
