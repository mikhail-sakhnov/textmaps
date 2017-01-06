package main

import (
	"github.com/soider/d"
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"textmap/maps"
)

type Application struct {
	dataDir string
	port    int
	doneCh chan struct{}
	router *mux.Router
	server *http.Server
}

func (a *Application) Run() {
	d.D("Running application")
	a.doneCh = make(chan struct{})
	d.D(a.server.ListenAndServe())
}

func (a *Application) Close()  {
	d.D("Stopping application, some cleanup")
	a.doneCh  <- struct{}{}
}

func (a *Application) Init() {
	a.initRouter()
	a.initLogger()
	a.initServer()
}

func (a *Application) initRouter() {
	mapService := maps.NewService(
		a.dataDir,
	)
	a.router = mux.NewRouter()
	a.router.Handle("/", maps.IndexHandler{
		mapService,
	})
	a.router.Handle("/f/{path:.+}", maps.MapHandler{
		mapService,
	})
}

func (a *Application) initLogger() {

}

func (a *Application) initServer() {
	d.D("test me")
	a.server = &http.Server{
		Handler: a.router,
		Addr:         "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func NewApplication(dataDir string, port int) *Application {
	return &Application{
		dataDir: dataDir,
		port:    port,
	}
}

func main() {
	app := NewApplication(
		"/Users/msahnov/Projects/textmaps",
		8080,
	)
	//c := make(chan os.Signal, 1)
	//signal.Notify(c, os.Interrupt)
	//go func(){
	//	for _ = range c {
	//		app.Close()
	//		os.Exit(0)
	//	}
	//}()
	app.Init()
	app.Run()

	<- app.doneCh
}
