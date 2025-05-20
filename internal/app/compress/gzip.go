package compress

import (
	"compress/gzip"
	"io"
	"net/http"
)

// CompressWriter реализует интерфейс http.ResponseWriter и сжимает данные
// с использованием gzip.
type CompressWriter struct {
	w  http.ResponseWriter
	gz *gzip.Writer
}

// NewCompressWriter constructor for type CompressWriter.
func NewCompressWriter(w http.ResponseWriter) *CompressWriter {
	return &CompressWriter{
		w:  w,
		gz: gzip.NewWriter(w),
	}
}

// Header возврашает заголовки http.ResponseWriter.
func (c *CompressWriter) Header() http.Header {
	return c.w.Header()
}

// Write сжимает и записывает данные.
func (c *CompressWriter) Write(b []byte) (int, error) {
	return c.gz.Write(b)
}

// WriteHeader записывает заголовки в CompressWriter
func (c *CompressWriter) WriteHeader(statusCode int) {
	c.w.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer
func (c *CompressWriter) Close() error {
	return c.gz.Close()
}

// CompressReader читает сжатые данные.
type CompressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

// NewCompressReader constructor for type CompressReader.
func NewCompressReader(r io.ReadCloser) (*CompressReader, error) {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return &CompressReader{
		r:  r,
		zr: gz,
	}, nil
}

// Read reads compressed data
func (c *CompressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close closed gzip.Reader
func (c *CompressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
