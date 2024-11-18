package openai

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Prompt based tag generation using OpenAI API
func (c *Client) AutoTagWithOpenAI(docText string) ([]string, error) {
	// Define the prompt for tag generation
	//prompt := fmt.Sprintf("Suggest relevant tags for the following content: %s", docText)

	// Define predefined tags
	predefinedTags := []string{"Account & Billing", "Order Status & Tracking", "Shipping & Returns", "Technical Troubleshooting", "Installation & Setup",
		"Product Information", "User Guide & How-To", "Software Updates & Maintenance", "Security & Privacy", "Feedback & Contact Support"}

	// Modify the prompt to include predefined tags
	prompt := fmt.Sprintf("From specifically the following tags: [ %s ], provide a comma-separated list of most relevant tags (only tags, no explanations, DO NOT use undefined tags) for the following content: %s",
		strings.Join(predefinedTags, ", "), docText)

	// Prepare the request
	request := map[string]interface{}{
		"model":      c.TagModel,
		"prompt":     prompt,
		"max_tokens": 50,
	}

	// Send the request to the OpenAI API
	response, err := c.Client.R().
		SetHeader("Authorization", "Bearer "+c.ApiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		Post("https://api.openai.com/v1/completions") // Use completion endpoint

	if err != nil {
		return nil, fmt.Errorf("error generating tags: %v", err)
	}

	// Parse the response JSON
	var result map[string]interface{}
	if err := json.Unmarshal(response.Body(), &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	fmt.Printf("Full API response: %+v\n", result)

	// Extract the text response
	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return nil, fmt.Errorf("error: missing 'choices' in API response")
	}

	text, ok := choices[0].(map[string]interface{})["text"].(string)
	if !ok {
		return nil, fmt.Errorf("error: 'text' not found in choices response")
	}

	// Parse tags from the text response (assuming they are comma-separated)
	tags := strings.Split(text, ",")
	for i := range tags {
		tags[i] = strings.TrimSpace(tags[i])
	}

	return tags, nil
}

// Generate tags using embeddings
// func (c *Client) AutoTagWithEmbeddings(docEmbedding []float64, tagEmbeddings map[string][]float64) []string {
// 	var tags []string
// 	for tag, embedding := range tagEmbeddings {
// 		if utils.CosineSimilarity(docEmbedding, embedding) > 0.7 {
// 			tags = append(tags, tag)
// 		}
// 	}
// 	return tags
// }
