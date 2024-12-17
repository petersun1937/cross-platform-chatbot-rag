package bot

import (
	"fmt"
	"strings"
)

// Defining states
//var screaming bool
//var useOpenAI bool = true // Default to using Dialogflow

// Process the commands sent by users and returns the message as a string
// func handleCommand(identifier interface{}, command string, bot Bot) (string, error) {
func (b *BaseBot) HandleCommand(command string) string {
	var message string
	//var err error

	switch command {
	/*case "/start":
	message = "Welcome to the bot!"*/
	case "/scream":
		b.conf.Screaming = true // Enable screaming mode
		message = "Scream mode enabled!"
	case "/whisper":
		b.conf.Screaming = false // Disable screaming mode
		message = "Scream mode disabled!"
	case "/openai":
		b.conf.UseOpenAI = true
		return "Disabling Dialogflow, only using OpenAI for responses."
	case "/dialogflow":
		b.conf.UseOpenAI = false
		return "Enabling Dialogflow for intent matching."

	case "/help":
		message = "You can type /openai for basic GPT-based responses with RAG, and /dialogflow to enable Dialogflow for intent matching (Dialogflow enabled by default)!"
	/*case "/custom":
	message = "This is a custom response!"*/
	default:
		message = "I don't know that command"
	}

	return message
}

// GetOpenAIResponse processes the user message and fetches a response from OpenAI API
func (b *BaseBot) GetOpenAIResponse(prompt string) (string, error) {
	client := b.openAIclient
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
	response = filterGPTResponse(response)

	return response, nil
}

func filterGPTResponse(response string) string {
	// Check if the response starts with "Response:" and remove it
	if strings.HasPrefix(response, "Response:") {
		return strings.TrimPrefix(response, "Response:")
	}
	// Return the response without unnecessary leading/trailing spaces
	return response
}
