package modelStorage

type Problem struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Difficulty   int    `json:"difficulty"`
	Link         string `json:"link"`
	PartitionKey string `json:"partitionKey"`
}
