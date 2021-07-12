package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	gremcos "github.com/supplyon/gremcos"

	"github.com/imeplusplus/dont-panic-api/app/dbOperations"
	"github.com/imeplusplus/dont-panic-api/app/model"
)

func CreateProblem(cosmos gremcos.Cosmos, w http.ResponseWriter, r *http.Request) {
	problem := model.Problem{}
	log.Println("New post with body: ", r.Body)
	if err := json.NewDecoder(r.Body).Decode(&problem); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println("New problem with fields: ", problem)

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
