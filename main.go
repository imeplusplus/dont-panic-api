package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"github.com/rs/zerolog"
	gremcos "github.com/supplyon/gremcos"
	"github.com/supplyon/gremcos/api"
)

type Subject struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	References []string `json:"references"`
	Difficulty int      `json:"difficulty"`
}

var subjects = []Subject{}
var cosmos gremcos.Cosmos
var logger zerolog.Logger
var err error

func main() {

	host := os.Getenv("CDB_HOST")
	username := os.Getenv("CDB_USERNAME")
	password := os.Getenv("CDB_KEY")

	fmt.Println(host)
	fmt.Println(username)
	fmt.Println(password)

	logger = zerolog.New(os.Stdout).Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: zerolog.TimeFieldFormat}).With().Timestamp().Logger()
	cosmos, err = gremcos.New(host,
		gremcos.WithAuth(username, password),
		gremcos.WithLogger(logger),
		gremcos.NumMaxActiveConnections(10),
		gremcos.ConnectionIdleTimeout(time.Second*30),
		gremcos.MetricsPrefix("myservice"),
	)

	if err != nil {
		return
	}

	r := mux.NewRouter()
	usersR := r.PathPrefix("/api/subject").Subrouter()
	usersR.Path("").Methods(http.MethodGet).HandlerFunc(getAllSubjects)
	usersR.Path("").Methods(http.MethodPost).HandlerFunc(createSubject)
	usersR.Path("/{id}").Methods(http.MethodGet).HandlerFunc(getSubjectByID)
	usersR.Path("/{id}").Methods(http.MethodPut).HandlerFunc(updateSubject)
	usersR.Path("/{id}").Methods(http.MethodDelete).HandlerFunc(deleteSubject)
	fmt.Println("Start listening")

	listenAddr := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}

	fmt.Println(http.ListenAndServe(listenAddr, r))
}

func getAllSubjects(w http.ResponseWriter, r *http.Request) {

	g := api.NewGraph("g")
	query := g.V()

	res, err := cosmos.ExecuteQuery(query)

	if err != nil {
		logger.Error().Err(err).Msg("Failed to execute a gremlin command")
		return
	}

	responses := api.ResponseArray(res)
	values, err := responses.ToValues()
	if err == nil {
		logger.Info().Msgf("Received Values: %v", values)
	}

	println(responses)

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responses); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getSubjectByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Not implemented")
}

func updateSubject(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Not implemented")
}

func deleteSubject(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Not implemented")
}

func createSubject(w http.ResponseWriter, r *http.Request) {
	subject := Subject{}
	if err := json.NewDecoder(r.Body).Decode(&subject); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	subjects = append(subjects, subject)
	response, err := json.Marshal(&subject)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
