package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	gremcos "github.com/supplyon/gremcos"
	"github.com/supplyon/gremcos/api"

	"github.com/imeplusplus/dont-panic-api/app/model"
)

func GetAllSubjects(cosmos gremcos.Cosmos, w http.ResponseWriter, r *http.Request) {
	g := api.NewGraph("g")
	query := g.V().HasLabel("subject")

	res, err := cosmos.ExecuteQuery(query)
	if err != nil {
		fmt.Println("Failed to execute a gremling command")
		//logger.Error().Err(err).Msg("Failed to execute a gremlin command")
		return
	}

	response := api.ResponseArray(res)
	vertices, err := response.ToVertices()

	if err == nil {
		fmt.Println(vertices)
	}

	subjects := VerticesToSubjects(vertices)

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(subjects); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func GetSubjectByName(cosmos gremcos.Cosmos, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	g := api.NewGraph("g")
	query := g.V().HasLabel("subject").Has("name", vars["name"])

	res, err := cosmos.ExecuteQuery(query)
	if err != nil {
		fmt.Println("Failed to execute a gremling command")
		//logger.Error().Err(err).Msg("Failed to execute a gremlin command")
		return
	}

	response := api.ResponseArray(res)
	vertices, err := response.ToVertices()

	if err == nil {
		fmt.Println(vertices)
	}

	var subject model.Subject
	if len(vertices) == 0 {
		return
	}

	subject, _ = VertexToSubject(vertices[0])

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(subject); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func UpdateSubject(cosmos gremcos.Cosmos, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Not implemented")
}

func DeleteSubject(cosmos gremcos.Cosmos, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Not implemented")
}

func CreateSubject(cosmos gremcos.Cosmos, w http.ResponseWriter, r *http.Request) {
	subject := model.Subject{}
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

func VerticesToSubjects(vertices []api.Vertex) []model.Subject {
	subjects := []model.Subject{}

	for _, v := range vertices {
		subject, err := VertexToSubject(v)
		if err == nil {
			subjects = append(subjects, subject)
		}
	}

	fmt.Println(JSONToString(subjects))

	return subjects
}

func VertexToSubject(vertex api.Vertex) (model.Subject, error) {
	var subject model.Subject

	if vertex.Label != "subject" {
		return subject, errors.New("Vertex is not a subject")
	}

	subject.Id = vertex.ID

	properties := vertex.Properties

	subject.Category = properties["category"][0].Value.AsString()
	subject.Name = properties["name"][0].Value.AsString()
	subject.Pk = properties["pk"][0].Value.AsString()
	//subject.Difficulty = int(properties["Difficulty"][0].Value.AsInt32())

	return subject, nil
}
