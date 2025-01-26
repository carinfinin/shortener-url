package store

import (
	"bufio"
	"encoding/json"
	"os"
)

type Producer struct {
	file   *os.File
	writer *bufio.Writer
}

func NewProducer(path string) (*Producer, error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &Producer{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}
func (p *Producer) WriteLine(line *Line) error {
	data, err := json.Marshal(&line)
	if err != nil {
		return err
	}

	if _, err := p.writer.Write(data); err != nil {
		return err
	}
	if err := p.writer.WriteByte('\n'); err != nil {
		return err
	}

	return p.writer.Flush()
}

func (p *Producer) Close() error {
	return p.file.Close()
}
