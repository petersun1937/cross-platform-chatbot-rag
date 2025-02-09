package togetherai

import (
	config "crossplatform_chatbot/configs"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

// Client struct for Together AI API
type Client struct {
	ApiKey string
	Model  string
	Client *resty.Client
}

// NewClient initializes a new Together AI API client
func NewClient() *Client {
	conf := config.GetConfig()

	client := resty.New()

	return &Client{
		ApiKey: conf.TogetherAIAPIKey,
		Model:  conf.TogetherAIModel,
		Client: client,
	}
}

// GetResponse sends a request to Together AI API and retrieves the response
func (c *Client) GetResponse(prompt string) (string, error) {
	request := map[string]interface{}{
		"model":       c.Model,
		"messages":    []map[string]string{{"role": "user", "content": prompt}},
		"max_tokens":  512,
		"temperature": 0.7,
	}

	response, err := c.Client.R().
		SetHeader("Authorization", "Bearer "+c.ApiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		Post("https://api.together.xyz/v1/chat/completions")

	if err != nil {
		return "", fmt.Errorf("error sending request to Together AI: %v", err)
	}

	if response.StatusCode() != 200 {
		return "", fmt.Errorf("Together AI API returned status code %d: %s", response.StatusCode(), response.String())
	}

	var result map[string]interface{}
	if err := json.Unmarshal(response.Body(), &result); err != nil {
		return "", fmt.Errorf("error parsing response from Together AI: %v", err)
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", errors.New("no response from Together AI")
	}

	text, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	if !ok {
		return "", errors.New("invalid response format from Together AI")
	}

	return cleanResponseText(text), nil
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
