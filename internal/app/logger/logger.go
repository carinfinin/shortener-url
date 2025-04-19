package logger

import (
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger = logrus.New()

func Configure(lvl string) error {
	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		return err
	}
	Log.SetLevel(level)

	return nil
}
