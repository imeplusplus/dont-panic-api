package model

type Problem struct {
	Id           string   `json:"id"`
	Name         string   `json:"name"`
	Subjects     []string `json:"subjects"`
	Difficulty   int      `json:"difficulty"`
	PartitionKey string   `json:"partitionKey"`
}
