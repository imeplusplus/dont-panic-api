package dbOperations

import (
	"errors"
	"fmt"

	gremcos "github.com/supplyon/gremcos"
	"github.com/supplyon/gremcos/api"
	"github.com/supplyon/gremcos/interfaces"

	modelAPI "github.com/imeplusplus/dont-panic-api/app/modelAPI"
	modelStorage "github.com/imeplusplus/dont-panic-api/app/modelStorage"
)

func GetSubjects(cosmos gremcos.Cosmos) ([]modelStorage.Subject, error) {
	g := api.NewGraph("g")
	query := g.V().HasLabel("subject")

	res, err := cosmos.ExecuteQuery(query)

	if err != nil {
		fmt.Println("Failed to execute a gremlin command " + query.String())
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

func GetSubjectByName(cosmos gremcos.Cosmos, name string) (modelStorage.Subject, error) {
	var subject modelStorage.Subject
	g := api.NewGraph("g")
	query := g.V().HasLabel("subject").Has("name", name)

	res, err := cosmos.ExecuteQuery(query)
	if err != nil {
		fmt.Println("Failed to execute a gremlin command " + query.String())
		//logger.Error().Err(err).Msg("Failed to execute a gremlin command")
		return subject, err
	}

	return getSubjectFromResponse(res)
}

func CreateSubject(cosmos gremcos.Cosmos, subject modelAPI.Subject) (modelStorage.Subject, error) {
	_, err := GetSubjectByName(cosmos, subject.Name)

	if err == nil {
		return modelStorage.Subject{}, errors.New("There is already a subject with name " + subject.Name)
	}

	g := api.NewGraph("g")

	query := g.AddV("subject").Property("partitionKey", "subject")
	query = addVertexProperties(query, subject)

	res, err := cosmos.ExecuteQuery(query)
	if err != nil {
		fmt.Println("Failed to execute a gremlin command " + query.String())
		// logger.Error().Err(err).Msg("Failed to execute gremlin command")
		return modelStorage.Subject{}, err
	}

	return getSubjectFromResponse(res)
}

func UpdateSubject(cosmos gremcos.Cosmos, subject modelAPI.Subject, name string) (modelStorage.Subject, error) {
	oldSubject, err := GetSubjectByName(cosmos, name)

	if err != nil {
		return modelStorage.Subject{}, errors.New("There is no subject with name " + oldSubject.Name)
	}

	g := api.NewGraph("g")
	query := addVertexProperties(g.VByStr(oldSubject.Id), subject)

	res, err := cosmos.ExecuteQuery(query)
	if err != nil {
		fmt.Println("Failed to execute a gremlin command " + query.String())
		//logger.Error().Err(err).Msg("Failed to execute a gremlin command")
		return modelStorage.Subject{}, err
	}

	return getSubjectFromResponse(res)
}

func DeleteSubject(cosmos gremcos.Cosmos, name string) error {
	g := api.NewGraph("g")
	query := g.V().HasLabel("subject").Has("name", name).Drop()

	_, err := cosmos.ExecuteQuery(query)

	return err
}

func addVertexProperties(vertex interfaces.Vertex, subject modelAPI.Subject) interfaces.Vertex {
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

func getSubjectFromResponse(res []interfaces.Response) (modelStorage.Subject, error) {
	var subject modelStorage.Subject
	response := api.ResponseArray(res)
	vertices, _ := response.ToVertices()

	if len(vertices) == 0 {
		return subject, errors.New("there is no vertex in the response")
	}

	subject = vertexToSubject(vertices[0])

	return subject, nil
}

func verticesToSubjects(vertices []api.Vertex) []modelStorage.Subject {
	subjects := []modelStorage.Subject{}

	for _, v := range vertices {
		subject := vertexToSubject(v)
		subjects = append(subjects, subject)
	}

	return subjects
}

func vertexToSubject(vertex api.Vertex) modelStorage.Subject {
	var subject modelStorage.Subject

	subject.Id = vertex.ID

	properties := vertex.Properties

	subject.Category = properties["category"][0].Value.AsString()
	subject.Name = properties["name"][0].Value.AsString()
	subject.PartitionKey = properties["partitionKey"][0].Value.AsString()
	subject.Difficulty = int(properties["difficulty"][0].Value.AsInt32())
	for _, p := range properties["references"] {
		subject.References = append(subject.References, p.Value.AsString())
	}

	return subject
}
