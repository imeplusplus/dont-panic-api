package model

type Subject struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	References []string `json:"references"`
	Difficulty int      `json:"difficulty"`
	Category   string   `json:"category"`
	Pk         string   `json:"pk"`
}
