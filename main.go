//go:generate statik -src=public
package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"os/signal"
	"time"
	"sync"
	"textmap/logger"
	mapsHandlers "textmap/maps/handlers"
	mapsServices "textmap/maps/services"
	"textmap/middlewares"
	_ "textmap/statik"
	"github.com/rakyll/statik/fs"
	"github.com/sirupsen/logrus"
	"context"
)
var s sync.Pool

type Application struct {
	debug   bool
	embedStatic bool
	dataDir string
	address string
	port    int
	doneCh  chan struct{}
	router  *mux.Router
	server  *http.Server

	logger *logrus.Entry
}

func (a *Application) Run() {
	go func() {
		if err := a.server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()
}

func (a *Application) Close() {
	a.doneCh <- struct{}{}
}

func (a *Application) Init() {
	a.initLogger()
	a.initRouter()
	a.initServer()
}

func (a *Application) initRouter() {
	mapService := mapsServices.NewService(
		a.dataDir,
	)
	a.router = mux.NewRouter()
	a.router.Handle("/api/list", mapsHandlers.IndexHandler{
		mapService,
	})
	a.router.Handle("/api/f/{path:.+}", mapsHandlers.MapHandler{
		mapService,
	})
	if a.embedStatic {
		a.logger.Warn("Serving static from binary")
		statikFS, err := fs.New()
		if err != nil {
			panic(err.Error())
		}
		a.router.Handle("/", http.FileServer(statikFS))
	} else {
		a.logger.Warn("Serving static from filesystem")
		a.router.Handle("/", http.FileServer(http.Dir("public")))
	}

}

func (a *Application) initLogger() {
	logger.Init(a.debug)
	a.logger = logger.FromContext(context.Background())
}

func (a *Application) initServer() {
	address := fmt.Sprintf("%s:%d", a.address, a.port)
	a.logger.WithField("address", address).Debug("Start listening")
	a.server = &http.Server{
		Handler:      middlewares.TraceMiddleware(a.router),
		Addr:         address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func NewApplication(dataDir string, address string, port int, debug, embedStatic bool) *Application {
	return &Application{
		dataDir: dataDir,
		address: address,
		port:    port,
		doneCh:  make(chan struct{}),
		debug:   debug,
		embedStatic: embedStatic,
	}
}

func main() {
	var directory string
	var address string
	var port int
	var debug bool
	var embedStatic bool

	flag.StringVar(&directory, "data", "./maps", "directory with data")
	flag.StringVar(&address, "host", "127.0.0.1", "host to listen on")
	flag.IntVar(&port, "port", 8080, "port to listen on")
	flag.BoolVar(&debug, "debug", false, "debug mode (more output)")
	flag.BoolVar(&embedStatic, "embed", true, "use static from binary or from ./public")

	flag.Parse()

	app := NewApplication(
		"/Users/msahnov/Projects/textmaps",
		address,
		port,
		debug,
		embedStatic,
	)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			app.Close()
			os.Exit(0)
		}
	}()
	app.Init()
	app.Run()
	<-app.doneCh
}
