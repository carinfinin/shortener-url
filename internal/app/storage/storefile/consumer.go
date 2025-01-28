package storefile

import (
	"encoding/json"
	"errors"
	"github.com/carinfinin/shortener-url/internal/app/models"
	"io"
	"os"
)

type Consumer struct {
	file    *os.File
	decoder *json.Decoder
}

func NewConsumer(path string) (*Consumer, error) {
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &Consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}
func (c *Consumer) Close() error {
	return c.file.Close()
}

func (c *Consumer) ReadAll() (map[string]string, error) {

	result := map[string]string{}
	line := models.Line{}
	for {
		err := c.decoder.Decode(&line)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		result[line.ID] = line.URL
	}

	return result, nil
}
