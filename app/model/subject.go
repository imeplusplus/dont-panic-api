package model

type Subject struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	References []string `json:"references"`
	Difficulty int      `json:"difficulty"`
}
