package huggingface

import (
	config "crossplatform_chatbot/configs"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

// Client struct to hold Hugging Face API configuration
type Client struct {
	ApiKey string
	Model  string
	Client *resty.Client
}

// NewClient initializes and returns a new Hugging Face API client
func NewClient() *Client {
	conf := config.GetConfig()

	client := resty.New()
	return &Client{
		ApiKey: conf.HuggingFaceAPIKey,
		Model:  conf.HuggingFaceModel,
		Client: client,
	}
}

// GetResponse sends a request to the Hugging Face Inference API and retrieves the response
func (c *Client) GetResponse(prompt string) (string, error) {
	request := map[string]interface{}{
		"inputs": prompt,
	}

	// Send the request to Hugging Face API
	response, err := c.Client.R().
		SetHeader("Authorization", "Bearer "+c.ApiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		Post(fmt.Sprintf("https://api-inference.huggingface.co/models/%s", c.Model))

	if err != nil {
		return "", fmt.Errorf("error sending request to Hugging Face: %v", err)
	}

	if response.StatusCode() != 200 {
		return "", fmt.Errorf("error: HuggingFace API returned status code %d: %s", response.StatusCode(), response.String())
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(response.Body(), &result); err != nil {
		return "", fmt.Errorf("error parsing response from Hugging Face: %v", err)
	}

	if len(result) == 0 {
		return "", errors.New("no response from Hugging Face")
	}

	text, ok := result[0]["generated_text"].(string)
	if !ok {
		return "", errors.New("invalid response format from Hugging Face")
	}

	return cleanResponseText(text), nil
}

// cleanResponseText trims unnecessary prefixes from the text response
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
