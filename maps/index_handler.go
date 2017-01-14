package maps

import (
	"github.com/soider/d"
	"net/http"
	"textmap/logger"
	"context"
)

type IndexHandler struct {
	Srv interface {
		GetAllMaps() (MapsCollection, error)
	}
}

func (ih IndexHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	switch req.Method {
	case "GET":
		ih.handleGet(ctx, rw)
	default:
		http.Error(rw, "Not implemented", 405)
	}
}

func (ih IndexHandler) handleGet(ctx context.Context, rw http.ResponseWriter) {
	maps, err := ih.Srv.GetAllMaps()
	if err != nil {
		http.Error(rw, "Internal server error", 500)
	}
	d.D(maps)
}
