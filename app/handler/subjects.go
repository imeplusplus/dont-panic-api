package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

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

func GetAllSubjects(cosmos gremcos.Cosmos, w http.ResponseWriter, r *http.Request) {

	g := api.NewGraph("g")
	query := g.V()

	res, err := cosmos.ExecuteQuery(query)
	if err != nil {
		fmt.Println("Failed to execute a gremling command")
		//logger.Error().Err(err).Msg("Failed to execute a gremlin command")
		return
	}

	responses := api.ResponseArray(res)
	values, err := responses.ToValues()
	if err == nil {
		fmt.Println(values)
		//logger.Info().Msgf("Received Values: %v", values)
	}

	println(responses)

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responses); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func GetSubjectByName(cosmos gremcos.Cosmos, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Not implemented")
}

func UpdateSubject(cosmos gremcos.Cosmos, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Not implemented")
}

func DeleteSubject(cosmos gremcos.Cosmos, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Not implemented")
}

func CreateSubject(cosmos gremcos.Cosmos, w http.ResponseWriter, r *http.Request) {
	subject := Subject{}
	if err := json.NewDecoder(r.Body).Decode(&subject); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//subjects = append(subjects, subject)
	response, err := json.Marshal(&subject)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
