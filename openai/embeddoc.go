package openai

import (
	"encoding/json"
	"fmt"
)

// EmbedText converts text to an embedding vector using OpenAI's embedding model
func (c *Client) EmbedText(text string) ([]float64, error) {
	//fmt.Println("Starting document embedding process...")

	//client := NewClient() // Create a new OpenAI client
	request := map[string]interface{}{
		"model": c.EmbModel,
		"input": text,
	}

	response, err := c.Client.R().
		SetHeader("Authorization", "Bearer "+c.ApiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		Post("https://api.openai.com/v1/embeddings")

	if err != nil {
		return nil, fmt.Errorf("error embedding document: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(response.Body(), &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	//data := result["data"].([]interface{})
	// Check if "data" exists in the response and is not nil
	data, ok := result["data"].([]interface{})
	if !ok || data == nil || len(data) == 0 {
		fmt.Printf("Full API response: %+v\n", result) // Log the entire response
		return nil, fmt.Errorf("error: missing or invalid 'data' in API response")
	}
	embedding := data[0].(map[string]interface{})["embedding"].([]interface{})

	//fmt.Println("Converting embedding to []float64...")
	embeddingFloat := make([]float64, len(embedding))
	for i, v := range embedding {
		embeddingFloat[i] = v.(float64)
	}

	//fmt.Println("Document embedding complete.")
	return embeddingFloat, nil
}

// EmbedSentencesBatch sends multiple sentences in a single request and gets embeddings for all of them
// func EmbedSentencesBatch(sentences []string) ([][]float64, error) {
// 	client := NewClient()
// 	request := map[string]interface{}{
// 		"model": "text-embedding-ada-002",
// 		"input": sentences, // Send all sentences in one request
// 	}

// 	response, err := client.Client.R().
// 		SetHeader("Authorization", "Bearer "+client.ApiKey).
// 		SetHeader("Content-Type", "application/json").
// 		SetBody(request).
// 		Post("https://api.openai.com/v1/embeddings")

// 	if err != nil {
// 		return nil, fmt.Errorf("error embedding sentences: %v", err)
// 	}

// 	var result map[string]interface{}
// 	if err := json.Unmarshal(response.Body(), &result); err != nil {
// 		return nil, fmt.Errorf("error parsing response: %v", err)
// 	}

// 	// Check for the "data" field in the response
// 	data, ok := result["data"].([]interface{})
// 	if !ok || data == nil {
// 		return nil, fmt.Errorf("error: missing or invalid 'data' in API response")
// 	}

// 	// Process embeddings for each sentence
// 	embeddings := make([][]float64, len(data))
// 	for i, item := range data {
// 		embedding, ok := item.(map[string]interface{})["embedding"].([]interface{})
// 		if !ok {
// 			return nil, fmt.Errorf("error: missing 'embedding' field for sentence %d", i)
// 		}
// 		embeddingFloat := make([]float64, len(embedding))
// 		for j, v := range embedding {
// 			embeddingFloat[j] = v.(float64)
// 		}
// 		embeddings[i] = embeddingFloat
// 	}

// 	return embeddings, nil
// }
