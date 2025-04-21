package compress

import (
	"compress/gzip"
	"io"
	"net/http"
)

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

func (c *CompressWriter) Header() http.Header {
	return c.w.Header()
}
func (c *CompressWriter) Write(b []byte) (int, error) {
	return c.gz.Write(b)
}

func (c *CompressWriter) WriteHeader(statusCode int) {
	c.w.WriteHeader(statusCode)
}

func (c *CompressWriter) Close() error {
	return c.gz.Close()
}

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
func (c *CompressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *CompressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
