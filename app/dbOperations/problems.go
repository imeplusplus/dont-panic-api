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
		return model.Problem{}, fmt.Errorf("there is already a problem with name %v", problem.Name)
	}

	g := api.NewGraph("g")

	query := g.AddV("problem").Property("partitionKey", "problem")
	query = addProblemVertexProperties(query, problem)

	res, err := cosmos.ExecuteQuery(query)
	if err != nil {
		fmt.Println("Failed to execute a gremlin command", query.String())
		return problem, err
	}

	problems, err := getProblemsFromResponse(res)
	if len(problems) == 0 {
		return problem, err
	}

	return problems[0], err
}

func GetProblems(cosmos gremcos.Cosmos) ([]model.Problem, error) {
	g := api.NewGraph("g")

	query := g.V().HasLabel("problem")
	res, err := cosmos.ExecuteQuery(query)
	if err != nil {
		fmt.Println("Failed to execute a gremlin command", query.String())
		return nil, err
	}

	return getProblemsFromResponse(res)
}

func GetProblemByName(cosmos gremcos.Cosmos, name string) (model.Problem, error) {
	var problem model.Problem
	g := api.NewGraph("g")
	query := g.V().HasLabel("problem").Has("name", name)

	res, err := cosmos.ExecuteQuery(query)
	if err != nil {
		fmt.Println("Failed to execute a gremlin command", query.String())
		return problem, err
	}

	problems, err := getProblemsFromResponse(res)
	if len(problems) == 0 {
		return problem, err
	}

	return problems[0], err
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

func getProblemsFromResponse(res []interfaces.Response) ([]model.Problem, error) {
	var problems []model.Problem
	response := api.ResponseArray(res)
	vertices, _ := response.ToVertices()

	if len(vertices) == 0 {
		return problems, errors.New("there is no data with type 'api.vertex' in the response. the graph query didn't return any vertex")
	}

	for _, v := range vertices {
		problems = append(problems, vertexToProblem(v))
	}

	return problems, nil
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
