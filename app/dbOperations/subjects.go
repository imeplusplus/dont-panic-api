package dbOperations

import (
	"errors"
	"fmt"

	gremcos "github.com/supplyon/gremcos"
	"github.com/supplyon/gremcos/api"
	"github.com/supplyon/gremcos/interfaces"

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

	subjects := verticesToSubjects(vertices)
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

	return getSubjectFromResponse(res)
}

func InsertSubject(cosmos gremcos.Cosmos, subject model.Subject) (model.Subject, error) {
	g := api.NewGraph("g")
	query := g.AddV("subject").
		Property("name", subject.Name).
		Property("difficulty", subject.Difficulty).
		Property("category", subject.Category).
		Property("partitionKey", "subject")

	for _, r := range subject.References {
		query = query.PropertyList("references", r)
	}

	res, err := cosmos.ExecuteQuery(query)
	if err != nil {
		fmt.Println("Failed to execute a gremling command")
		//logger.Error().Err(err).Msg("Failed to execute a gremlin command")
		return subject, err
	}

	return getSubjectFromResponse(res)
}

func getSubjectFromResponse(res []interfaces.Response) (model.Subject, error) {
	var subject model.Subject
	response := api.ResponseArray(res)
	vertices, err := response.ToVertices()
	if len(vertices) == 0 {
		return subject, errors.New("Vertex is not a subject")
	}

	subject, err = vertexToSubject(vertices[0])
	if err != nil {
		return subject, err
	}

	return subject, nil
}

func verticesToSubjects(vertices []api.Vertex) []model.Subject {
	subjects := []model.Subject{}

	for _, v := range vertices {
		subject, err := vertexToSubject(v)
		if err == nil {
			subjects = append(subjects, subject)
		}
	}

	return subjects
}

func vertexToSubject(vertex api.Vertex) (model.Subject, error) {
	var subject model.Subject

	if vertex.Label != "subject" {
		return subject, errors.New("Vertex is not a subject")
	}

	subject.Id = vertex.ID

	properties := vertex.Properties

	subject.Category = properties["category"][0].Value.AsString()
	subject.Name = properties["name"][0].Value.AsString()
	subject.PartitionKey = properties["partitionKey"][0].Value.AsString()
	subject.Difficulty = int(properties["difficulty"][0].Value.AsInt32())
	for _, p := range properties["references"] {
		subject.References = append(subject.References, p.Value.AsString())
	}

	return subject, nil
}
