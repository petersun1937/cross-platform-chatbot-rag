package document_proc

import (
	"fmt"
)

// GetRelevantTags scores each tag's embedding against the query embedding
// and returns a list of tags that pass the threshold.
func GetRelevantTags(queryEmbedding []float64, tagEmbeddings map[string][]float64, threshold float64) []string {
	var relevantTags []string

	for tag, embedding := range tagEmbeddings {
		cosineScore := cosineSimilarity(queryEmbedding, embedding)
		keywordScore := keywordMatchScore(tag, tag) // Simple keyword match for tags
		combinedScore := weightedScore(cosineScore, keywordScore)

		if combinedScore >= threshold {
			relevantTags = append(relevantTags, tag)
			fmt.Printf("Tag: %s, Score: %.4f\n", tag, combinedScore)
		}
	}
	return relevantTags
}
