package compress

import (
	"bytes"
	"compress/gzip"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http/httptest"
	"testing"
)

func TestCompress(t *testing.T) {

	data := []byte("Привет компресс тест")
	w := httptest.NewRecorder()
	wcom := NewCompressWriter(w)
	_, err := wcom.Write(data)
	wcom.Close()

	assert.NoError(t, err, err)

	gr, err := gzip.NewReader(w.Body)
	assert.NoError(t, err, err)
	value, err := io.ReadAll(gr)
	gr.Close()
	assert.NoError(t, err, err)
	assert.Equal(t, data, value)
}

func TestDecompress(t *testing.T) {
	data := []byte("Привет компресс тест")

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write(data)
	gz.Close()

	gr, err := NewCompressReader(io.NopCloser(&buf))
	assert.NoError(t, err, err)

	decompressedData, err := io.ReadAll(gr)
	assert.NoError(t, err)
	gr.Close()
	assert.Equal(t, data, decompressedData)

}
