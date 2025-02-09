package bot

import (
	"fmt"
	"strings"
)

// Process the commands sent by users and returns the message as a string
// func handleCommand(identifier interface{}, command string, bot Bot) (string, error) {
func (b *BaseBot) HandleCommand(command string) string {
	var message string

	switch command {
	case "/scream":
		b.conf.Screaming = true // Enable screaming mode
		message = "Scream mode enabled!"
	case "/whisper":
		b.conf.Screaming = false // Disable screaming mode
		message = "Scream mode disabled!"
	case "/openai":
		b.conf.UseOpenAI = true
		b.conf.UseMistral = false
		b.conf.UseMETA = false
		message = "Using OpenAI GPT-4 for responses."
	case "/mistral":
		b.conf.UseOpenAI = false
		b.conf.UseMistral = true
		b.conf.UseMETA = false
		message = "Using Mistral AI Mistral-large model for responses."
	case "/meta":
		b.conf.UseOpenAI = false
		b.conf.UseMistral = false
		b.conf.UseMETA = true
		message = "Using META Llama model from Together AI for responses."
	case "/dialogflow":
		b.conf.UseDialogflow = true
		message = "Enabling Dialogflow for intent matching."
	case "/disable_dialogflow":
		b.conf.UseDialogflow = false
		message = "Dialogflow disabled."
	case "/help":
		message = "You can type the following commands:\n"
		message += "**/openai** - Use OpenAI GPT-4 for responses.\n"
		message += "**/mistral** - Use Mistral AI Mistral-large model.\n"
		message += "**/meta** - Use META Llama model from Together AI.\n"
		message += "**/dialogflow** - Enable Dialogflow for intent matching.\n"
		message += "**/disable_dialogflow** - Disable Dialogflow intent matching. Use similarity score based retrieval only.\n"
		message += "Note that only one AI model can be active at a time, while you can enable Dialogflow independently. "
		message += "OpenAI and Dialogflow enabled by default."
	default:
		message = "I don't know that command"
	}

	return message
}

// GetOpenAIResponse processes the user message and fetches a response from OpenAI API
func (b *BaseBot) GetOpenAIResponse(prompt string) (string, error) {
	client := b.aiClients.OpenAI
	response, err := client.GetResponse(prompt)
	if err != nil {
		return "", fmt.Errorf("error fetching response from OpenAI: %v", err)
	}

	// Check if response is empty or missing expected fields
	if response == "" {
		return "", fmt.Errorf("no valid response from OpenAI. Please try again later")
	}

	fmt.Printf("OpenAI response: %s \n", response)

	// Filter out "Response:" if it exists
	response = filterResponse(response)

	return response, nil
}

// GetHuggingFaceResponse processes the user message and fetches a response from Hugging Face Inference API
/*func (b *BaseBot) GetHuggingFaceResponse(prompt string) (string, error) {
	client := b.huggingfaceClient
	response, err := client.GetResponse(prompt)
	if err != nil {
		return "", fmt.Errorf("error fetching response from HuggingFace: %v", err)
	}

	// Check if response is empty or missing expected fields
	if response == "" {
		return "", fmt.Errorf("no valid response from HuggingFace. Please try again later")
	}

	fmt.Printf("Hugging Face response: %s \n", response)

	// Filter out "Response:" if it exists
	response = filterResponse(response)

	return response, nil
}*/

// GetMistralResponse processes the user message and fetches a response from Mistral AI API
func (b *BaseBot) GetMistralResponse(prompt string) (string, error) {
	client := b.aiClients.Mistral
	response, err := client.GetResponse(prompt)
	if err != nil {
		return "", fmt.Errorf("error fetching response from Mistral AI: %v", err)
	}

	// Check if response is empty or missing expected fields
	if response == "" {
		return "", fmt.Errorf("no valid response from Mistral AI. Please try again later")
	}

	fmt.Printf("Mistral response: %s \n", response)

	// Filter out "Response:" if it exists
	response = filterResponse(response)

	return response, nil
}

// GetTogetherAIResponse processes the user message and fetches a response from Together AI API
func (b *BaseBot) GetTogetherAIResponse(prompt string) (string, error) {
	if b.aiClients.TogetherAI == nil {
		return "", fmt.Errorf("error: Together AI client is not initialized")
	}
	response, err := b.aiClients.TogetherAI.GetResponse(prompt)
	if err != nil {
		return "", fmt.Errorf("error fetching response from Together AI: %v", err)
	}
	return filterResponse(response), nil
}

// filterResponse removes unwanted prefixes and phrases from the response
func filterResponse(response string) string {
	// Define a list of unwanted phrases to filter out
	unwantedPhrases := []string{
		"Response:",
		"Bot response:",
		"AI generated response:",
	}

	// Trim spaces from the input first
	response = strings.TrimSpace(response)

	// Loop through unwanted phrases and remove them if found at the beginning
	for _, phrase := range unwantedPhrases {
		if strings.HasPrefix(response, phrase) {
			response = strings.TrimPrefix(response, phrase)
			response = strings.TrimSpace(response) // Remove any leading spaces left
		}
	}

	return response
}
