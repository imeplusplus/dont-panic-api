package model

type Subject struct {
	Name         string   `json:"name"`
	Category     string   `json:"category"`
	Difficulty   int      `json:"difficulty"`
	References   []string `json:"references"`
	Dependencies []string `json:"dependencies"`
}
