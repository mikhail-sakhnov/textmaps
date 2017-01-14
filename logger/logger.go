package logger

import (
	"context"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"net/http"
)

type requestIDKeyType struct{}

var requestIDKey = requestIDKeyType{}

func TraceMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		requestID := uuid.NewV4()

		ctx := context.WithValue(context.Background(), requestIDKey, requestID.String())
		next.ServeHTTP(rw, req.WithContext(ctx))
	})
}

func LoggerFromContext(ctx context.Context) *logrus.Entry {
	requestID := ctx.Value(requestIDKey).(string)
	return logrus.WithField("request_id", requestID)
}
