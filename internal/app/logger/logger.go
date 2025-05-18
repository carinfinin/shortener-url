package logger

import (
	"github.com/sirupsen/logrus"
)

// Log global constant containing the *logrus.Logger object.
var Log *logrus.Logger = logrus.New()

// Configure sets the logging level.
func Configure(lvl string) error {
	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		return err
	}
	Log.SetLevel(level)

	return nil
}
