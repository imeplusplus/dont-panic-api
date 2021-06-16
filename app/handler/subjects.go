package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	gremcos "github.com/supplyon/gremcos"

	"github.com/imeplusplus/dont-panic-api/app/dbOperations"
	"github.com/imeplusplus/dont-panic-api/app/model"
)

func GetSubjects(cosmos gremcos.Cosmos, w http.ResponseWriter, r *http.Request) {
	subjects, err := dbOperations.GetSubjects(cosmos)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(subjects); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func GetSubject(cosmos gremcos.Cosmos, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	subject, err := dbOperations.GetSubjectByName(cosmos, vars["name"])

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(subject); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func UpdateSubject(cosmos gremcos.Cosmos, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	subject := model.Subject{}
	var err error
	if err = json.NewDecoder(r.Body).Decode(&subject); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	subject, err = dbOperations.UpdateSubject(cosmos, subject, vars["name"])

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(subject); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func DeleteSubject(cosmos gremcos.Cosmos, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := dbOperations.DeleteSubject(cosmos, vars["name"])

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}

func CreateSubject(cosmos gremcos.Cosmos, w http.ResponseWriter, r *http.Request) {
	subject := model.Subject{}
	var err error
	if err = json.NewDecoder(r.Body).Decode(&subject); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	subject, err = dbOperations.CreateSubject(cosmos, subject)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(subject); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
