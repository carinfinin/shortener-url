package store

import (
	"bufio"
	"encoding/json"
	"os"
)

type Consumer struct {
	file    *os.File
	scanner *bufio.Scanner
}

func NewConsumer(path string) (*Consumer, error) {
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &Consumer{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}
func (c *Consumer) Close() error {
	return c.file.Close()
}

func (c *Consumer) ReadAll() (map[string]string, error) {

	result := map[string]string{}
	for c.scanner.Scan() {
		data := c.scanner.Bytes()

		line := Line{}
		err := json.Unmarshal(data, &line)
		if err != nil {
			return nil, err
		}
		result[line.ID] = line.URL
	}

	return result, nil
}
