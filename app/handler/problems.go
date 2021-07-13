package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	gremcos "github.com/supplyon/gremcos"

	"github.com/imeplusplus/dont-panic-api/app/dbOperations"
	modelStorage "github.com/imeplusplus/dont-panic-api/app/modelStorage"
)

func CreateProblem(cosmos gremcos.Cosmos, w http.ResponseWriter, r *http.Request) {
	problem := modelStorage.Problem{}
	if err := json.NewDecoder(r.Body).Decode(&problem); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	problem, err := dbOperations.CreateProblem(cosmos, problem)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(problem)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func GetProblems(cosmos gremcos.Cosmos, w http.ResponseWriter, r *http.Request) {
	problems, err := dbOperations.GetProblems(cosmos)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(problems)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func GetProblem(cosmos gremcos.Cosmos, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	problem, err := dbOperations.GetProblemByName(cosmos, vars["name"])

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(problem)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func UpdateProblem(cosmos gremcos.Cosmos, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	problem := modelStorage.Problem{}
	if err := json.NewDecoder(r.Body).Decode(&problem); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	problem, err := dbOperations.UpdateProblem(cosmos, problem, vars["name"])

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(problem)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func DeleteProblem(cosmos gremcos.Cosmos, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := dbOperations.DeleteProblem(cosmos, vars["name"])

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}