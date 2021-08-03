package model

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
)

type Subject struct {
	Id           string   `json:"id"`
	Name         string   `json:"name"`
	References   []string `json:"references"`
	Difficulty   int      `json:"difficulty"`
	Category     string   `json:"category"`
	PartitionKey string   `json:"partitionKey"`
}

func PrettyPrint(subjects ...interface{}) string {
	subjectJSON, err := json.MarshalIndent(subjects, "", "  ")
	if err != nil {
		log.Error().Stack().Err(err).Msg("Couldn't make resource pretty to print")
		return fmt.Sprintf("%v", subjects)
	}
	return string(subjectJSON)
}
