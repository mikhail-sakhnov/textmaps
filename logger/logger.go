package logger

import (
	"github.com/sirupsen/logrus"
	"context"
	"textmap/middlewares"
	"log"
	"os"
)

func Init(debug bool) {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})
	// Output to stdout instead of the default stderr, could also be a file.
	log.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.WarnLevel)
	}
}

func FromContext(ctx context.Context) *logrus.Entry {
	requestID := ctx.Value(middlewares.RequestIDKey).(string)
	return logrus.WithField("request_id", requestID)
}