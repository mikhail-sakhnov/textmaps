package logger

import (
	"context"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"textmap/middlewares"
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
	requestID, found := ctx.Value(middlewares.RequestIDKey).(string)
	if found {
		return logrus.WithField("request_id", requestID)
	}
	return logrus.WithField("request_id", "no_request_id")
}
