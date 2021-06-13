package dbOperations

import (
	"errors"
	"fmt"

	gremcos "github.com/supplyon/gremcos"
	"github.com/supplyon/gremcos/api"

	"github.com/imeplusplus/dont-panic-api/app/model"
)

func GetSubjects(cosmos gremcos.Cosmos) ([]model.Subject, error) {
	g := api.NewGraph("g")
	query := g.V().HasLabel("subject")

	res, err := cosmos.ExecuteQuery(query)

	if err != nil {
		fmt.Println("Failed to execute a gremling command")
		//logger.Error().Err(err).Msg("Failed to execute a gremlin command")
		return nil, err
	}

	response := api.ResponseArray(res)
	vertices, err := response.ToVertices()

	if err != nil {
		return nil, err
	}

	subjects := VerticesToSubjects(vertices)
	return subjects, nil
}

func GetSubjectByName(cosmos gremcos.Cosmos, name string) (model.Subject, error) {
	var subject model.Subject
	g := api.NewGraph("g")
	query := g.V().HasLabel("subject").Has("name", name)

	res, err := cosmos.ExecuteQuery(query)
	if err != nil {
		fmt.Println("Failed to execute a gremling command")
		//logger.Error().Err(err).Msg("Failed to execute a gremlin command")
		return subject, err
	}

	response := api.ResponseArray(res)
	vertices, err := response.ToVertices()
	if len(vertices) == 0 {
		return subject, errors.New("Vertex is not a subject")
	}

	subject, err = VertexToSubject(vertices[0])
	if err != nil {
		return subject, err
	}

	return subject, nil
}

func VerticesToSubjects(vertices []api.Vertex) []model.Subject {
	subjects := []model.Subject{}

	for _, v := range vertices {
		subject, err := VertexToSubject(v)
		if err == nil {
			subjects = append(subjects, subject)
		}
	}

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
