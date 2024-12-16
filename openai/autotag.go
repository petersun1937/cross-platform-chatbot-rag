package openai

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkoukk/tiktoken-go"
)

// Prompt based tag generation using OpenAI API
func (c *Client) AutoTagWithOpenAI(docText string) ([]string, error) {
	// Define the prompt for tag generation
	//prompt := fmt.Sprintf("Suggest relevant tags for the following content: %s", docText)

	// Define predefined tags
	predefinedTags := []string{"Account & Billing", "Order Status & Tracking", "Shipping & Returns", "Technical Troubleshooting", "Installation & Setup",
		"Product Information", "User Guide & How-To", "Software Updates & Maintenance", "Security & Privacy", "Feedback & Contact Support"}

	// Prepare the full predefined tags as part of the prompt
	tagList := strings.Join(predefinedTags, ", ")
	reminder := "Reminder: DO NOT include any additional undefined tags. Provide ONLY tags from the list."

	// Construct the base prompt
	/*prompt := fmt.Sprintf("From specifically the following tags: [ %s ], provide a comma-separated list of most relevant tags (only tags, no explanations, DO NOT use undefined tags) for the following content: %s. Reminder: Provide only tags, no explanations, DO NOT use undefined tags",
	tagList, docText)*/
	basePrompt := fmt.Sprintf(`From specifically the following tags: [ %s ], 
	provide a comma-separated list of most relevant tags (only tags, no explanations, DO NOT use undefined tags) for the following content: %s.
	
	%s`, tagList, docText, reminder)

	// Tokenize the prompt into words
	/*tokens := strings.Fields(prompt)
	tokensize := len(tokens)

	// Check if the prompt exceeds maxTagTokens and trim from the end if necessary
	if tokensize > c.TagTokenSize {
		prompt = strings.Join(tokens[:c.TagTokenSize], " ") // Retain only the allowed number of tokens
	}*/

	// Initialize OpenAI tokenizer for token counting
	tkm, err := tiktoken.GetEncoding("cl100k_base") // Replace with model-specific encoding
	if err != nil {
		return nil, fmt.Errorf("error initializing tokenizer: %v", err)
	}

	// Tokenize the prompt to count tokens accurately
	tokens := tkm.Encode(basePrompt, nil, nil)
	// Maximum prompt size (reserving 50 tokens for completion)
	maxPromptSize := c.TagTokenSize - 50

	// Check if the prompt exceeds MaxPromptSize
	if len(tokens) > maxPromptSize {
		// Trim the content section (docText) to fit within the limit
		tokens = tokens[:maxPromptSize]
		// Decode the trimmed tokens back to text
		basePrompt = tkm.Decode(tokens)

		// Ensure the reminder is still included at the end
		if !strings.HasSuffix(basePrompt, reminder) {
			// If trimming removed the reminder, re-add it
			basePrompt += "\n\n" + reminder
		}
	}

	// Prepare the request
	request := map[string]interface{}{
		"model":      c.TagModel,
		"prompt":     basePrompt,
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
