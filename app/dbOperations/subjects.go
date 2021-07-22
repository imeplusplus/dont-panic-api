package dbOperations

import (
	"fmt"

	"github.com/rs/zerolog/log"
	gremcos "github.com/supplyon/gremcos"
	"github.com/supplyon/gremcos/api"
	"github.com/supplyon/gremcos/interfaces"

	"github.com/imeplusplus/dont-panic-api/app/logger"
	storageModel "github.com/imeplusplus/dont-panic-api/app/model/storage"
)

func GetSubjects(cosmos gremcos.Cosmos) ([]storageModel.Subject, error) {
	g := api.NewGraph("g")
	query := g.V().HasLabel("subject")

	res, err := cosmos.ExecuteQuery(query)

	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("Failed to execute the gremlin command: %v", query))
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

func GetSubjectByName(cosmos gremcos.Cosmos, name string) (storageModel.Subject, error) {
	var subject storageModel.Subject
	g := api.NewGraph("g")
	query := g.V().HasLabel("subject").Has("name", name)

	res, err := cosmos.ExecuteQuery(query)
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("Failed to execute the gremlin command: %v", query))
		return subject, err
	}

	return getSubjectFromResponse(res)
}

func CreateSubject(cosmos gremcos.Cosmos, subject storageModel.Subject) (storageModel.Subject, error) {
	_, err := GetSubjectByName(cosmos, subject.Name)

	if err == nil {
		err := logger.ErrorResourceAlreadyExists{ResourceName: subject.Name}.Error()
		log.Error().Err(err).Msg(fmt.Sprintf("There is already a subject with name %s in the database", subject.Name))
		return storageModel.Subject{}, err
	}

	g := api.NewGraph("g")

	query := g.AddV("subject").Property("partitionKey", "subject")

	query = addVertexProperties(query, subject)

	res, err := cosmos.ExecuteQuery(query)
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("Failed to execute the gremlin command: %v", query))
		return storageModel.Subject{}, err
	}

	return getSubjectFromResponse(res)
}

func UpdateSubject(cosmos gremcos.Cosmos, subject storageModel.Subject, name string) (storageModel.Subject, error) {
	oldSubject, err := GetSubjectByName(cosmos, name)

	if err != nil {
		err := logger.ErrorResourceNotFound{ResourceName: name}.Error()
		log.Error().Err(err).Msg(fmt.Sprintf("Couldn't find subject with name %s in the database", name))
		return storageModel.Subject{}, err
	}

	g := api.NewGraph("g")
	query := addVertexProperties(g.VByStr(oldSubject.Id), subject)

	res, err := cosmos.ExecuteQuery(query)
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("Failed to execute the gremlin command: %v", query))
		return storageModel.Subject{}, err
	}

	return getSubjectFromResponse(res)
}

func DeleteSubject(cosmos gremcos.Cosmos, name string) error {
	g := api.NewGraph("g")
	query := g.V().HasLabel("subject").Has("name", name).Drop()

	_, err := cosmos.ExecuteQuery(query)

	return err
}

func addVertexProperties(vertex interfaces.Vertex, subject storageModel.Subject) interfaces.Vertex {
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

func getSubjectFromResponse(res []interfaces.Response) (storageModel.Subject, error) {
	var subject storageModel.Subject
	response := api.ResponseArray(res)
	vertices, _ := response.ToVertices()

	if len(vertices) == 0 {
		err := logger.ErrorResourceNotFound{ResourceName: "response"}.Error()
		log.Error().Err(err).Msg("")
		return subject, err
	}

	subject = vertexToSubject(vertices[0])

	return subject, nil
}

func verticesToSubjects(vertices []api.Vertex) []storageModel.Subject {
	subjects := []storageModel.Subject{}

	for _, v := range vertices {
		subject := vertexToSubject(v)
		subjects = append(subjects, subject)
	}

	return subjects
}

func vertexToSubject(vertex api.Vertex) storageModel.Subject {
	var subject storageModel.Subject

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
