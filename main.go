package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Subject struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	References []string `json:"references"`
	Difficulty int      `json:"difficulty"`
}

var subjects = []Subject{}

func main() {
	r := mux.NewRouter()
	usersR := r.PathPrefix("/subject").Subrouter()
	usersR.Path("").Methods(http.MethodGet).HandlerFunc(getAllSubjects)
	usersR.Path("").Methods(http.MethodPost).HandlerFunc(createSubject)
	usersR.Path("/{id}").Methods(http.MethodGet).HandlerFunc(getSubjectByID)
	usersR.Path("/{id}").Methods(http.MethodPut).HandlerFunc(updateSubject)
	usersR.Path("/{id}").Methods(http.MethodDelete).HandlerFunc(deleteSubject)
	fmt.Println("Start listening")
	fmt.Println(http.ListenAndServe("", r))
}

func getAllSubjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(subjects); err != nil {
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
