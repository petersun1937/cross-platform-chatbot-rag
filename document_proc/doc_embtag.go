package document_proc

import (
	config "crossplatform_chatbot/configs"
	"crossplatform_chatbot/openai"
	"crossplatform_chatbot/repository"
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

// StoreDocumentChunks splits the given text into chunks and embed them one by one via embedding model
func StoreDocumentChunks(Filename, sessionID, text string, embConfig config.EmbeddingConfig, client *openai.Client, dao repository.DAO) error {
	// Chunk text and embed each chunk
	chunks := OverlapChunk(text, embConfig.ChunkSize, embConfig.MinChunkSize)
	for i, chunk := range chunks {
		embedding, err := client.EmbedText(chunk)
		if err != nil {
			return fmt.Errorf("error embedding chunk %d: %v", i, err)
		}

		chunkID := fmt.Sprintf("%s_chunk_%d_%s", Filename, i, sessionID)
		if err := dao.CreateDocumentEmbedding(Filename, sessionID, chunkID, chunk, embedding); err != nil {
			return fmt.Errorf("error storing chunks: %v", err)
		}
	}
	return nil
}

// AutoTagDocumentEmbeddings
// func AutoTagDocumentEmbeddings(sessionID string, client *openai.Client, dao repository.DAO, tagEmbeddings map[string][]float64) ([]string, error) {
// 	chunkEmbeddings, err := dao.GetChunkEmbeddings(sessionID)
// 	if err != nil {
// 		return nil, fmt.Errorf("error retrieving chunk embeddings: %w", err)
// 	}

// 	combinedEmbedding, err := utils.AverageEmbeddings(chunkEmbeddings)
// 	if err != nil {
// 		return nil, fmt.Errorf("error combining chunk embeddings: %w", err)
// 	}

// 	tags := client.AutoTagWithEmbeddings(combinedEmbedding, tagEmbeddings)
// 	return tags, nil
// }
