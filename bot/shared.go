package bot

import (
	config "crossplatform_chatbot/configs"
	openai "crossplatform_chatbot/openai"
	"crossplatform_chatbot/utils"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/line/line-bot-sdk-go/linebot"
)

// Defining states
var screaming bool
var useOpenAI bool = true // Default to using Dialogflow

// Process the commands sent by users and returns the message as a string
// func handleCommand(identifier interface{}, command string, bot Bot) (string, error) {
func handleCommand(command string) string {
	var message string
	//var err error

	switch command {
	case "/start":
		message = "Welcome to the bot!"
	case "/scream":
		screaming = true // Enable screaming mode
		message = "Scream mode enabled!"
	case "/whisper":
		screaming = false // Disable screaming mode
		message = "Scream mode disabled!"
	case "/openai":
		useOpenAI = true
		return "Switched to OpenAI for responses."
	case "/dialogflow":
		useOpenAI = false
		return "Switched to Dialogflow for responses."
	// case "/menu":
	// 	// Handle menu sending based on platform
	// 	/*switch platform {
	// 	case LINE:
	// 		if event, ok := identifier.(*linebot.Event); ok {
	// 			err = sendLineMenu(event.ReplyToken) // Send a menu to LINE
	// 		} else {
	// 			err = fmt.Errorf("invalid identifier type for LINE platform")
	// 		}
	// 	case TELEGRAM:
	// 		if chatID, ok := identifier.(int64); ok {
	// 			err = sendMenu(chatID) // Send a menu to Telegram
	// 		} else {
	// 			err = fmt.Errorf("invalid identifier type for Telegram platform")
	// 		}
	// 	}*/
	// 	err = bot.sendMenu(identifier)
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	return "", nil
	case "/help":
		message = "Here are some commands you can use: /start, /help, /scream, /whisper, /menu. You can also type /openai for GPT-based responses, and /dialogflow to switch to rule-based Dialogflow responses!"
	case "/custom":
		message = "This is a custom response!"
	default:
		message = "I don't know that command"
	}

	return message
}

// handleMessageDialogflow handles messages from different platforms
func handleMessageDialogflow(platform Platform, identifier interface{}, text string, bot Bot) {
	// Determine sessionID based on platform
	sessionID, err := getSessionID(platform, identifier)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Send the message to Dialogflow and receive a response
	conf := config.GetConfig()
	response, err := utils.DetectIntentText(conf.DialogflowProjectID, sessionID, text, "en")
	if err != nil {
		fmt.Printf("Error detecting intent: %v\n", err)
		return
	}

	// Process and send the Dialogflow response to the appropriate platform
	if err := bot.handleDialogflowResponse(response, identifier); err != nil {
		fmt.Println(err)
	}
}

// getSessionID extracts the session ID based on the platform and identifier
func getSessionID(platform Platform, identifier interface{}) (string, error) {
	switch platform {
	case LINE:
		if event, ok := identifier.(*linebot.Event); ok {
			return event.Source.UserID, nil
		}
		return "", fmt.Errorf("invalid LINE event identifier")
	case TELEGRAM:
		if message, ok := identifier.(*tgbotapi.Message); ok {
			return strconv.FormatInt(message.Chat.ID, 10), nil
		}
		return "", fmt.Errorf("invalid Telegram message identifier")
	case FACEBOOK:
		if recipientID, ok := identifier.(string); ok {
			return recipientID, nil
		}
		return "", fmt.Errorf("invalid Messenger recipient identifier")
	case GENERAL:
		if sessionID, ok := identifier.(string); ok {
			return sessionID, nil
		}
		return "", fmt.Errorf("invalid session identifier")
	default:
		return "", fmt.Errorf("unsupported platform")
	}
}

// GetOpenAIResponse processes the user message and fetches a response from OpenAI API
func GetOpenAIResponse(prompt string) (string, error) {
	client := openai.NewClient()
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
