package handlers

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"textmap/maps/entities"
)

type MapHandler struct {
	Srv interface {
		GetMapTextContent(ctx context.Context, path string) (entities.SingleMap, error)
	}
}

func (lh MapHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	switch req.Method {
	case "GET":
		reqData := mux.Vars(req)
		lh.handleGet(ctx, rw, reqData["path"])
	default:
		http.Error(rw, "Not implemented", 405)
	}
}

func (lh MapHandler) handleGet(ctx context.Context, rw http.ResponseWriter, reqPath string) {
	singleMap, err := lh.Srv.GetMapTextContent(ctx, reqPath)
	if err != nil {
		http.Error(rw, "Internal server error: "+err.Error(), 500)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(rw)
	err = encoder.Encode(singleMap)
	if err != nil {
		http.Error(rw, err.Error(), 500)
	}
}
