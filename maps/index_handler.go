package maps

import (
	"net/http"
	"github.com/soider/d"
)

type IndexHandler struct {
	Srv interface{
		GetAllMaps() (MapsCollection, error)
	}
}

func (ih IndexHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		ih.handleGet(rw)
	default:
		http.Error(rw, "Not implemented", 405)
	}
}

func (ih IndexHandler) handleGet(rw http.ResponseWriter) {
	maps, err := ih.Srv.GetAllMaps()
	if err!=nil {
		http.Error(rw, "Internal server error", 500)
	}
	d.D(maps)
}