package dbOperations

import (
	"errors"
	"fmt"

	gremcos "github.com/supplyon/gremcos"
	"github.com/supplyon/gremcos/api"
	"github.com/supplyon/gremcos/interfaces"

	"github.com/imeplusplus/dont-panic-api/app/model"
)

func CreateProblem(cosmos gremcos.Cosmos, problem model.Problem) (model.Problem, error) {
	_, err := GetProblemByName(cosmos, problem.Name)

	if err == nil {
		return model.Problem{}, errors.New("There is already a problem with name " + problem.Name)
	}

	g := api.NewGraph("g")

	query := g.AddV("problem").Property("partitionKey", "problem")
	query = addProblemVertexProperties(query, problem)

	res, err := cosmos.ExecuteQuery(query)
	if err != nil {
		fmt.Println("Failed to execute a gremlin command", query.String())
		return problem, err
	}

	return getProblemFromResponse(res)
}

func GetProblemByName(cosmos gremcos.Cosmos, name string) (model.Problem, error) {
	var problem model.Problem
	g := api.NewGraph("g")
	query := g.V().HasLabel("Problem").Has("Name", name)

	res, err := cosmos.ExecuteQuery(query)
	if err != nil {
		fmt.Println("Failed to execute a gremlin command", query.String())
		return problem, err
	}

	return getProblemFromResponse(res)
}

func addProblemVertexProperties(vertex interfaces.Vertex, problem model.Problem) interfaces.Vertex {
	vertex = vertex.
		Property("name", problem.Name).
		Property("difficulty", problem.Difficulty)

	if len(problem.Subjects) > 0 {
		vertex = vertex.Property("subjects", problem.Subjects[0])

		for _, s := range problem.Subjects[1:] {
			vertex = vertex.PropertyList("subjects", s)
		}
	}

	return vertex
}

func getProblemFromResponse(res []interfaces.Response) (model.Problem, error) {
	var problem model.Problem
	response := api.ResponseArray(res)
	vertices, _ := response.ToVertices()

	if len(vertices) == 0 {
		return problem, errors.New("there is no vertex in the response")
	}

	problem = vertexToProblem(vertices[0])

	return problem, nil
}

func vertexToProblem(vertex api.Vertex) model.Problem {
	var problem model.Problem

	problem.Id = vertex.ID

	properties := vertex.Properties
	problem.Name = properties["name"][0].Value.AsString()
	problem.Difficulty = int(properties["difficulty"][0].Value.AsInt32())
	problem.PartitionKey = properties["partitionKey"][0].Value.AsString()
	for _, p := range properties["subjects"] {
		problem.Subjects = append(problem.Subjects, p.Value.AsString())
	}

	return problem
}
