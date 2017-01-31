package handlers

import (
	"context"
	"encoding/json"
	"github.com/soider/d"
	"net/http"
	"textmap/logger"
	"textmap/maps/entities"
)

type IndexHandler struct {
	Srv interface {
		GetAllMaps(ctx context.Context) (entities.MapsCollection, error)
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
	maps, err := ih.Srv.GetAllMaps(ctx)
	logger.FromContext(ctx).
		WithField("maps_count", len(maps)).
		WithField("error", err).
		Debug("Loaded maps")
	d.D(maps)
	if err != nil {
		http.Error(rw, "Internal server error: "+err.Error(), 500)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(rw)
	err = encoder.Encode(maps)
	if err != nil {
		http.Error(rw, "Internal server error: "+err.Error(), 500)
	}
}
