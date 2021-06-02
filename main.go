package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"github.com/rs/zerolog"
	gremcos "github.com/supplyon/gremcos"

	"github.com/imeplusplus/dont-panic-api/handler"
)

type App struct {
	Router *mux.Router
	Cosmos gremcos.Cosmos
	Logger zerolog.Logger
}

func (app *App) Initialize() {
	host := os.Getenv("CDB_HOST")
	username := os.Getenv("CDB_USERNAME")
	password := os.Getenv("CDB_KEY")

	fmt.Println(host)
	fmt.Println(username)
	fmt.Println(password)

	app.Logger = zerolog.New(os.Stdout).Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: zerolog.TimeFieldFormat}).With().Timestamp().Logger()
	app.Cosmos, err = gremcos.New(host,
		gremcos.WithAuth(username, password),
		gremcos.WithLogger(logger),
		gremcos.NumMaxActiveConnections(10),
		gremcos.ConnectionIdleTimeout(time.Second*30),
		gremcos.MetricsPrefix("myservice"),
	)

	app.Router = mux.NewRouter()
	app.setRouters()
}

func (app *App) setRouters() {
	// Routing for subjects
	app.Get("/api/subject", app.handleRequest(handler.GetAllSubjects))
	app.Post("/api/subject", app.handleRequest(handler.CreateSubject))
	app.Get("/api/subject/{name}", app.handleRequest(handler.GetSubjectByName))
	app.Put("/papi/subject{name}", app.handleRequest(handler.UpdateSubject))
	app.Delete("/papi/subject{name}", app.handleRequest(handler.DeleteSubject))
}

// Get wraps the router for GET method
func (app *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	app.Router.HandleFunc(path, f).Methods("GET")
}

// Post wraps the router for POST method
func (app *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	app.Router.HandleFunc(path, f).Methods("POST")
}

// Put wraps the router for PUT method
func (app *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	app.Router.HandleFunc(path, f).Methods("PUT")
}

// Delete wraps the router for DELETE method
func (app *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	app.Router.HandleFunc(path, f).Methods("DELETE")
}

type RequestHandlerFunction func(cosmos gremcos.Cosmos, w http.ResponseWriter, r *http.Request)

func (app *App) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(app.Cosmos, w, r)
	}
}

// Run the app on it's router
func (app *App) Run(host string) {
	fmt.Println(http.ListenAndServe(host, app.Router))
}

var cosmos gremcos.Cosmos
var logger zerolog.Logger
var err error

func main() {

	app := App{}
	app.Initialize()

	listenAddr := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}
	app.Run(listenAddr)
}
