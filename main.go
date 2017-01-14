package main

import (
	"flag"

	"fmt"
	"github.com/gorilla/mux"
	"github.com/soider/d"
	"net/http"
	"os"
	"os/signal"
	"textmap/maps"
	"time"
	"textmap/logger"
	log "github.com/sirupsen/logrus"

)

type Application struct {
	debug bool
	dataDir string
	address string
	port    int
	doneCh  chan struct{}
	router  *mux.Router
	server  *http.Server
}

func (a *Application) Run() {
	d.D("Running application")
	go func() {
		if err := a.server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()
}

func (a *Application) Close() {
	d.D("Stopping application, some cleanup")
	a.doneCh <- struct{}{}
	d.D("Stopping application, some cleanup after")
}

func (a *Application) Init() {
	a.initLogger()
	a.initRouter()
	a.initServer()
}

func (a *Application) initRouter() {
	mapService := maps.NewService(
		a.dataDir,
	)
	a.router = mux.NewRouter()
	a.router.HandleFunc("/", logger.TraceMiddleware(maps.IndexHandler{
		mapService,
	}))
	a.router.HandleFunc("/f/{path:.+}", logger.TraceMiddleware(maps.MapHandler{
		mapService,
	}))
}

func (a *Application) initLogger() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})
	// Output to stdout instead of the default stderr, could also be a file.
	log.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	if a.debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}


}

func (a *Application) initServer() {
	address := fmt.Sprintf("%s:%d", a.address, a.port)
	d.D("Listen on ", address)
	a.server = &http.Server{
		Handler:      a.router,
		Addr:         address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func NewApplication(dataDir string, address string, port int, debug bool) *Application {
	return &Application{
		dataDir: dataDir,
		address: address,
		port:    port,
		doneCh:  make(chan struct{}),
		debug: debug,
	}
}

func main() {
	var directory string
	var address string
	var port int
	var debug bool

	flag.StringVar(&directory, "data", "./maps", "directory with data")
	flag.StringVar(&address, "host", "127.0.0.1", "host to listen on")
	flag.IntVar(&port, "port", 8080, "port to listen on")
	flag.BoolVar(&debug, "debug", false, "debug mode (more output)")

	flag.Parse()

	app := NewApplication(
		directory,
		address,
		port,
		debug,
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
