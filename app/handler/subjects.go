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
	storageSubjects, err := dbOperations.GetSubjects(cosmos)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	apiSubjects := []apiModel.Subject{}

	for _, subject := range storageSubjects {
		apiSubjects = append(apiSubjects, NewApiSubject(subject, nil))
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(apiSubjects)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func GetSubject(cosmos gremcos.Cosmos, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	storageSubject, err := dbOperations.GetSubjectByName(cosmos, vars["name"])

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	apiSubject := NewApiSubject(storageSubject, nil)

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(apiSubject)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func UpdateSubject(cosmos gremcos.Cosmos, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	apiSubject := apiModel.Subject{}
	err := json.NewDecoder(r.Body).Decode(&apiSubject)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	storageSubject, err := dbOperations.UpdateSubject(cosmos, NewStorageSubject(apiSubject), vars["name"])
	apiSubject = NewApiSubject(storageSubject, nil)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	err = json.NewEncoder(w).Encode(apiSubject)

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
	apiSubject := apiModel.Subject{}
	var err error
	if err = json.NewDecoder(r.Body).Decode(&apiSubject); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	storageSubject, err := dbOperations.CreateSubject(cosmos, NewStorageSubject(apiSubject))
	apiSubject = NewApiSubject(storageSubject, nil)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(apiSubject)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func NewApiSubject(storageSubject storageModel.Subject, dependencies []string) apiModel.Subject {
	apiSubject := new(apiModel.Subject)

	apiSubject.Name = storageSubject.Name
	apiSubject.Category = storageSubject.Category
	apiSubject.Difficulty = storageSubject.Difficulty
	apiSubject.References = storageSubject.References
	apiSubject.Dependencies = dependencies

	return *apiSubject
}

func NewStorageSubject(apiSubject apiModel.Subject) storageModel.Subject {
	storageSubject := new(storageModel.Subject)

	storageSubject.Name = apiSubject.Name
	storageSubject.Category = apiSubject.Category
	storageSubject.Difficulty = apiSubject.Difficulty
	storageSubject.References = apiSubject.References

	return *storageSubject
}
