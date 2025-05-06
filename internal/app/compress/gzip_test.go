package compress

import (
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
