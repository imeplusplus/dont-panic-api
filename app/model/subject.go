package model

type Subject struct {
	Difficulty   int      `json:"difficulty"`
	Name         string   `json:"name"`
	Category     string   `json:"category"`
	References   []string `json:"references"`
	Dependencies []string `json:"dependencies"`
}
