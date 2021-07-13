package dbOperations

import (
	"errors"
	"fmt"

	gremcos "github.com/supplyon/gremcos"
	"github.com/supplyon/gremcos/api"
	"github.com/supplyon/gremcos/interfaces"

	modelStorage "github.com/imeplusplus/dont-panic-api/app/modelStorage"
)

func CreateProblem(cosmos gremcos.Cosmos, problem modelStorage.Problem) (modelStorage.Problem, error) {
	_, err := GetProblemByName(cosmos, problem.Name)

	if err == nil {
		return modelStorage.Problem{}, fmt.Errorf("there is already a problem with name %v", problem.Name)
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
		return modelStorage.Problem{}, err
	}

	return problems[0], err
}

func GetProblems(cosmos gremcos.Cosmos) ([]modelStorage.Problem, error) {
	g := api.NewGraph("g")

	query := g.V().HasLabel("problem")
	res, err := cosmos.ExecuteQuery(query)
	if err != nil {
		fmt.Println("Failed to execute a gremlin command", query.String())
		return nil, err
	}

	return getProblemsFromResponse(res)
}

func GetProblemByName(cosmos gremcos.Cosmos, name string) (modelStorage.Problem, error) {
	var problem modelStorage.Problem
	g := api.NewGraph("g")
	query := g.V().HasLabel("problem").Has("name", name)

	res, err := cosmos.ExecuteQuery(query)
	if err != nil {
		fmt.Println("Failed to execute a gremlin command", query.String())
		return problem, err
	}

	problems, err := getProblemsFromResponse(res)
	if len(problems) == 0 {
		return modelStorage.Problem{}, err
	}

	return problems[0], err
}

func UpdateProblem(cosmos gremcos.Cosmos, problem modelStorage.Problem, name string) (modelStorage.Problem, error) {
	if problem.Name != name {
		return modelStorage.Problem{}, fmt.Errorf("can't change the property 'name' of the Problem")
	}
	oldProblem, err := GetProblemByName(cosmos, name)
	if err != nil {
		return oldProblem, fmt.Errorf("there is no problem with name '%v' to update in the database", name)
	}

	g := api.NewGraph("g")
	query := addProblemVertexProperties(g.VByStr(oldProblem.Id), problem)

	res, err := cosmos.ExecuteQuery(query)
	if err != nil {
		fmt.Println("Failed to execute a gremlin command", query.String())
		return problem, err
	}

	problems, err := getProblemsFromResponse(res)
	if len(problems) == 0 {
		return modelStorage.Problem{}, err
	}

	return problems[0], err
}

func DeleteProblem(cosmos gremcos.Cosmos, name string) error {
	_, err := GetProblemByName(cosmos, name)

	if err != nil {
		return fmt.Errorf("there is no problem with name '%v' to delete in the database", name)
	}

	g := api.NewGraph("g")
	query := g.V().HasLabel("problem").Has("name", name).Drop()

	_, err = cosmos.ExecuteQuery(query)

	return err
}

func addProblemVertexProperties(vertex interfaces.Vertex, problem modelStorage.Problem) interfaces.Vertex {
	vertex = vertex.
		Property("name", problem.Name).
		Property("difficulty", problem.Difficulty).
		Property("link", problem.Link)

	return vertex
}

func getProblemsFromResponse(res []interfaces.Response) ([]modelStorage.Problem, error) {
	var problems []modelStorage.Problem
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

func vertexToProblem(vertex api.Vertex) modelStorage.Problem {
	var problem modelStorage.Problem

	problem.Id = vertex.ID

	properties := vertex.Properties
	problem.Name = properties["name"][0].Value.AsString()
	problem.Difficulty = int(properties["difficulty"][0].Value.AsInt32())
	problem.Link = properties["link"][0].Value.AsString()
	problem.PartitionKey = properties["partitionKey"][0].Value.AsString()

	return problem
}
