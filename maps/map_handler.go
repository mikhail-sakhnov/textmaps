package maps

import (
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
	"github.com/soider/d"
)

type MapHandler struct {
	Srv interface {
		GetMapTextContent(path string) (SingleMap, error)
	}
}

func (lh MapHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		reqData := mux.Vars(req)
		lh.handleGet(rw, reqData["path"])
	default:
		http.Error(rw, "Not implemented", 405)
	}
}

func (lh MapHandler) handleGet(rw http.ResponseWriter, reqPath string) {
	singleMap, err := lh.Srv.GetMapTextContent(reqPath)
	if err != nil {
		http.Error(rw,
			err.Error(),
			500)
		return
	}
	encoder := json.NewEncoder(rw)
	d.D(singleMap)
	encoder.Encode(singleMap)
}
