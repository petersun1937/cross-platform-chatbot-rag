package mistral

import (
	config "crossplatform_chatbot/configs"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

// Client struct for Mistral AI API
type Client struct {
	ApiKey string
	Model  string
	Client *resty.Client
}

// NewClient initializes a new Mistral AI API client
func NewClient() *Client {
	conf := config.GetConfig()

	client := resty.New()

	return &Client{
		ApiKey: conf.MistralAPIKey,
		Model:  conf.MistralModel,
		Client: client,
	}
}

// GetResponse sends a request to Mistral AI API and retrieves the response
func (c *Client) GetResponse(prompt string) (string, error) {
	request := map[string]interface{}{
		"model":       c.Model, // Mistral model name
		"messages":    []map[string]string{{"role": "user", "content": prompt}},
		"max_tokens":  512,
		"temperature": 0.7,
	}

	// Send request to Mistral API
	response, err := c.Client.R().
		SetHeader("Authorization", "Bearer "+c.ApiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		Post("https://api.mistral.ai/v1/chat/completions") // Mistral API endpoint

	if err != nil {
		return "", fmt.Errorf("error sending request to Mistral: %v", err)
	}

	if response.StatusCode() != 200 {
		return "", fmt.Errorf("error: Mistral API returned status code %d: %s", response.StatusCode(), response.String())
	}

	var result map[string]interface{}
	if err := json.Unmarshal(response.Body(), &result); err != nil {
		return "", fmt.Errorf("error parsing response from Mistral: %v", err)
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", errors.New("no response from Mistral")
	}

	text, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	if !ok {
		return "", errors.New("invalid response format from Mistral")
	}

	// Clean up the response to remove unnecessary prefixes
	cleanedText := cleanResponseText(text)

	return cleanedText, nil
}

// cleanResponseText trims unwanted prefixes from the text
func cleanResponseText(text string) string {
	lowerText := strings.ToLower(text)

	if strings.HasPrefix(lowerText, "response:") {
		text = strings.TrimSpace(text[len("response:"):])
	}

	if strings.HasPrefix(lowerText, "assistant:") {
		text = strings.TrimSpace(text[len("assistant:"):])
	}

	if strings.HasPrefix(lowerText, "bot:") {
		text = strings.TrimSpace(text[len("bot:"):])
	}
	return text
}
