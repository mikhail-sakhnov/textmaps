package middlewares

import (
	"context"
	"github.com/satori/go.uuid"
	"net/http"
)

type requestIDKeyType struct{}

var RequestIDKey = requestIDKeyType{}

func TraceMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		requestID := uuid.NewV4()
		ctx := context.WithValue(context.Background(), RequestIDKey, requestID.String())
		next.ServeHTTP(rw, req.WithContext(ctx))
	})
}


