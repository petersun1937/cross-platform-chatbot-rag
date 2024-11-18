package bot

import (
	document "crossplatform_chatbot/document_proc"
	"crossplatform_chatbot/utils"
	"fmt"
	"strconv"
	"strings"

	"cloud.google.com/go/dialogflow/apiv2/dialogflowpb"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/line/line-bot-sdk-go/linebot"
)

// handleMessageDialogflow handles a message from the platform, sends it to Dialogflow for intent detection,
// retrieves the corresponding context using RAG, generates a response with OpenAI, and sends it back to the user.
func (b *BaseBot) handleMessageDialogflow(platform Platform, identifier interface{}, msg string, bot Bot) {
	// Determine session ID
	sessionID, err := getSessionID(platform, identifier)
	if err != nil {
		fmt.Printf("Error getting session ID: %v\n", err)
		return
	}

	// Send message to Dialogflow and get intent
	dialogflowResponse, err := b.fetchDialogflowResponse(sessionID, msg)
	if err != nil {
		fmt.Printf("Error fetching Dialogflow response: %v\n", err)
		return
	}

	// Capture intent from Dialogflow response
	intent := dialogflowResponse.GetQueryResult().GetIntent().GetDisplayName()
	fmt.Printf("Captured Intent: %s\n", intent)

	// Retrieve document context using intent-based tags
	context, err := b.fetchDocumentContext(intent, msg)
	if err != nil {
		fmt.Printf("Error fetching document context: %v\n", err)
		return
	}
	fmt.Printf("Retrieved Context:\n%s\n", context)

	// Generate response using OpenAI with context
	prompt := fmt.Sprintf("Context:\n%s\nUser query: %s", context, msg)

	finalResponse, err := b.GetOpenAIResponse(prompt)
	//finalResponse, err := b.openAIclient.GetResponse(prompt)
	if err != nil {
		fmt.Printf("Error generating OpenAI response: %v\n", err)
		return
	}

	// Send final response back to the platform
	if err := bot.sendResponse(sessionID, finalResponse); err != nil {
		fmt.Printf("Error sending response: %v\n", err)
	}
}

// HandleDialogflowIntent processes the user message, maps the intent, and generates a response.
// func (b *BaseBot) HandleDialogflowIntent(message string) (string, error) {
// 	// Map user message to an intent
// 	intent := mapIntent(message)

// 	// Retrieve tags based on intent
// 	tags := mapTags(intent)
// 	if len(tags) == 0 {
// 		return "Intent not recognized for RAG context.", nil
// 	}

// 	// Retrieve document chunks based on tags
// 	contextChunks, err := b.retrieveChunksByTags(tags)
// 	if err != nil {
// 		return "", fmt.Errorf("error retrieving document chunks: %w", err)
// 	}

// 	// Create a prompt with context for OpenAI
// 	context := strings.Join(contextChunks, "\n")
// 	prompt := fmt.Sprintf("Context:\n%s\nUser query: %s", context, message)

// 	// Get response from OpenAI
// 	response, err := b.openAIclient.GetResponse(prompt)
// 	if err != nil {
// 		return "", fmt.Errorf("error generating RAG response: %w", err)
// 	}

// 	return response, nil
// }

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

// fetchDialogflowResponse sends the message to Dialogflow and retrieves the response with detected intent.
func (b *BaseBot) fetchDialogflowResponse(sessionID, text string) (*dialogflowpb.DetectIntentResponse, error) {
	conf := b.conf
	response, err := utils.DetectIntentText(conf.DialogflowProjectID, sessionID, text, "en")
	if err != nil {
		return nil, fmt.Errorf("error detecting intent with Dialogflow: %v", err)
	}
	return response, nil
}

// fetchDocumentContext retrieves the document chunks based on the detected intent's associated tags.
func (b *BaseBot) fetchDocumentContext(intent string, userMessage string) (string, error) {
	// Define tags based on intent
	tags := mapTags(intent)
	var context string
	var err error

	if len(tags) > 0 {
		// Use tag-based retrieval if tags are available
		context, err = b.retrieveChunksByTags(tags)
		if err != nil {
			return "", err
		}
	} else {
		// Fallback to similarity-based chunk retrieval if no tags are found, same as basic OpenAI mode
		context, err = b.fallbackContext(userMessage)
		if err != nil {
			return "", err
		}
	}

	return context, nil
	/*
		if len(tags) == 0 {
				return "", fmt.Errorf("no tags defined for intent: %s", intent)
			}

			// Retrieve document chunks from the database
			contextChunks, err := b.retrieveChunksByTags(tags)
			if err != nil {
				return "", fmt.Errorf("error retrieving document chunks: %w", err)
			}

			// Combine chunks into a single context
			return strings.Join(contextChunks, "\n"), nil
	*/
}

// mapIntent maps user input to predefined intents.
func mapIntent(message string) string {
	switch {
	case strings.Contains(strings.ToLower(message), "faq"):
		return "FAQ Intent"
	case strings.Contains(strings.ToLower(message), "product"):
		return "Product Inquiry Intent"
	case strings.Contains(strings.ToLower(message), "troubleshoot"):
		return "Troubleshooting Intent"
	case strings.Contains(strings.ToLower(message), "install"):
		return "Installation Intent"
	default:
		return "Default Intent"
	}
}

// Defines tags associated with an intent.
func mapTags(intent string) []string {
	switch intent {
	case "FAQ Intent":
		return []string{"FAQs", "Product Information", "User Guide & How-To"}
	case "Product Inquiry Intent":
		return []string{"Product Information", "Account & Billing", "Order Status & Tracking"}
	case "Troubleshooting Intent":
		return []string{"Technical Troubleshooting", "Installation & Setup", "Security & Privacy"}
	case "Installation Intent":
		return []string{"Installation & Setup"}
	default:
		return nil
	}
}

// retrieveChunksByTags fetches document chunks that match the specified tags and ranks them by similarity.
// func (b *BaseBot) retrieveChunksByTags(tags []string, userMessage string) (string, error) {
// 	// Step 1: Fetch document embeddings for the matching tags
// 	documentEmbeddings, chunkText, err := b.dao.GetDocumentChunksByTagsWithEmbeddings(tags)
// 	if err != nil {
// 		return "", fmt.Errorf("error retrieving document chunks by tags: %w", err)
// 	}

// 	if len(documentEmbeddings) == 0 {
// 		return "", fmt.Errorf("no document chunks found for the specified tags")
// 	}

// 	// Step 2: Rank the retrieved chunks based on similarity to the user's query
// 	topChunks, err := document.RetrieveTopNChunks(userMessage, documentEmbeddings, b.embConfig.NumTopChunks, chunkText, b.embConfig.ScoreThreshold)
// 	if err != nil {
// 		return "", fmt.Errorf("error ranking document chunks by similarity: %w", err)
// 	}

// 	if len(topChunks) == 0 {
// 		return "", fmt.Errorf("no relevant document chunks found based on similarity")
// 	}

// 	// Step 3: Combine the top-ranked chunks into a single context string
// 	return strings.Join(topChunks, "\n"), nil
// }

// retrieveChunksByTags fetches document chunks that match the specified tags
func (b *BaseBot) retrieveChunksByTags(tags []string) (string, error) {
	// Fetch document embeddings from the database where tags match
	documentEmbeddings, err := b.dao.GetDocumentChunksByTags(tags)
	if err != nil {
		return "", fmt.Errorf("error retrieving document chunks: %w", err)
	}

	// Extract text content from matched embeddings
	var contextChunks []string
	for _, chunk := range documentEmbeddings {
		contextChunks = append(contextChunks, chunk.DocText)
	}

	return strings.Join(contextChunks, "\n"), nil
}

// fallbackContext retrieves document chunks based on similarity to the user's message, functions as basic openAI mode.
func (b *BaseBot) fallbackContext(userMessage string) (string, error) {
	// Fetch all document embeddings and chunk text
	documentEmbeddings, chunkText, err := b.dao.FetchEmbeddings()
	if err != nil {
		return "", fmt.Errorf("error fetching embeddings: %w", err)
	}

	// Use similarity-based chunk retrieval
	topChunks, err := document.RetrieveTopNChunks(userMessage, documentEmbeddings, b.embConfig.NumTopChunks, chunkText, b.embConfig.ScoreThreshold)
	if err != nil {
		return "", fmt.Errorf("error retrieving similar document chunks: %w", err)
	}

	if len(topChunks) == 0 {
		return "", fmt.Errorf("no relevant document chunks found for the user message")
	}

	return strings.Join(topChunks, "\n"), nil
}
