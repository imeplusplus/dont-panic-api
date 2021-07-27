package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	gremcos "github.com/supplyon/gremcos"

	"github.com/imeplusplus/dont-panic-api/app/dbOperations"
	apiModel "github.com/imeplusplus/dont-panic-api/app/model/api"
	storageModel "github.com/imeplusplus/dont-panic-api/app/model/storage"
)

func GetSubjects(cosmos gremcos.Cosmos, w http.ResponseWriter, _ *http.Request) {
	res, err := dbOperations.GetSubjects(cosmos)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func GetSubject(cosmos gremcos.Cosmos, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	res, err := dbOperations.GetSubjectByName(cosmos, vars["name"])

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func UpdateSubject(cosmos gremcos.Cosmos, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	subject := apiModel.Subject{}
	err := json.NewDecoder(r.Body).Decode(&subject)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := dbOperations.UpdateSubject(cosmos, storageModel.Subject(subject), vars["name"])

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	err = json.NewEncoder(w).Encode(res)

	if err != nil {
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
	subject := apiModel.Subject{}
	var err error
	if err = json.NewDecoder(r.Body).Decode(&subject); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := dbOperations.CreateSubject(cosmos, storageModel.Subject(subject))

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(res)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
