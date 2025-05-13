package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigure(t *testing.T) {
	t.Run("info", func(t *testing.T) {
		Configure("info")
		assert.Equal(t, logrus.InfoLevel, Log.GetLevel())
	})

	t.Run("debug", func(t *testing.T) {
		Configure("debug")
		assert.Equal(t, logrus.DebugLevel, Log.GetLevel())
	})
}
