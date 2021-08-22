package dbOperations

import (
	"errors"
	"fmt"

	gremcos "github.com/supplyon/gremcos"
	"github.com/supplyon/gremcos/api"
	"github.com/supplyon/gremcos/interfaces"

	"github.com/imeplusplus/dont-panic-api/app/model"
)

const (
	label = "subject"
)

func GetSubjects(cosmos gremcos.Cosmos) ([]model.Subject, error) {
	g := api.NewGraph("g")
	query := g.V().HasLabel(label)
	vertices, err := getVerticesFromQuery(cosmos, query)

	if err != nil {
		return nil, err
	}

	subjects := verticesToSubjects(vertices)
	return subjects, nil
}

func GetSubjectByName(cosmos gremcos.Cosmos, name string) (model.Subject, error) {
	var subject model.Subject

	vertex, err := getVertexByName(cosmos, name)

	if err != nil {
		return subject, err
	}

	return vertexToSubject(vertex)
}

func CreateSubject(cosmos gremcos.Cosmos, subject model.Subject) (model.Subject, error) {
	_, err := GetSubjectByName(cosmos, subject.Name)

	if err == nil {
		return model.Subject{}, errors.New("There is already a subject with name " + subject.Name)
	}

	g := api.NewGraph("g")

	query := g.AddV(label).Property("partitionKey", label)
	query = addVertexProperties(query, subject)

	vertices, err := getVerticesFromQuery(cosmos, query)

	if err != nil {
		return subject, err
	}

	if len(vertices) == 0 {
		return subject, errors.New("There is no vertex in the response")
	}

	subject, err = vertexToSubject(vertices[0])
	return subject, err
}

func UpdateSubject(cosmos gremcos.Cosmos, subject model.Subject, name string) (model.Subject, error) {
	oldSubjectVertex, err := getVertexByName(cosmos, name)

	if err != nil {
		return model.Subject{}, errors.New("There is no subject with name " + name)
	}

	g := api.NewGraph("g")
	query := addVertexProperties(g.VByStr(oldSubjectVertex.ID), subject)

	vertices, err := getVerticesFromQuery(cosmos, query)

	if err != nil {
		return model.Subject{}, err
	}

	if len(vertices) == 0 {
		return model.Subject{}, errors.New("There is no vertex in the response")
	}

	return vertexToSubject(vertices[0])
}

func DeleteSubject(cosmos gremcos.Cosmos, name string) error {
	g := api.NewGraph("g")
	query := g.V().HasLabel(label).Has("name", name).Drop()

	_, err := cosmos.ExecuteQuery(query)

	return err
}

func addVertexProperties(vertex interfaces.Vertex, subject model.Subject) interfaces.Vertex {
	vertex = vertex.
		Property("name", subject.Name).
		Property("difficulty", subject.Difficulty).
		Property("category", subject.Category)

	if len(subject.References) >= 0 {
		vertex = vertex.Property("references", subject.References[0])

		for _, r := range subject.References[1:] {
			vertex = vertex.PropertyList("references", r)
		}
	}

	return vertex
}

func getVertexByName(cosmos gremcos.Cosmos, name string) (api.Vertex, error) {
	var vertex api.Vertex

	g := api.NewGraph("g")
	query := g.V().HasLabel(label).Has("name", name)
	vertices, err := getVerticesFromQuery(cosmos, query)

	if err != nil {
		return vertex, err
	}

	if len(vertices) == 0 {
		return vertex, errors.New("There is no vertex in the response")
	}

	return vertices[0], err
}

func getVerticesFromQuery(cosmos gremcos.Cosmos, query interfaces.Vertex) ([]api.Vertex, error) {
	res, err := cosmos.ExecuteQuery(query)
	if err != nil {
		fmt.Println("Failed to execute a gremlin command " + query.String())
		//logger.Error().Err(err).Msg("Failed to execute a gremlin command")
		return nil, err
	}

	response := api.ResponseArray(res)
	vertices, err := response.ToVertices()

	if err != nil {
		fmt.Println("Failed to convert response to vertices")
	}

	return vertices, err
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

	properties := vertex.Properties

	subject.Category = properties["category"][0].Value.AsString()
	subject.Name = properties["name"][0].Value.AsString()
	subject.Difficulty = int(properties["difficulty"][0].Value.AsInt32())
	for _, p := range properties["references"] {
		subject.References = append(subject.References, p.Value.AsString())
	}

	return subject, nil
}
