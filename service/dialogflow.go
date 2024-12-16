package service

import (
	"context"
	"crossplatform_chatbot/bot"
	"fmt"
	"strings"

	document "crossplatform_chatbot/document_proc"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"cloud.google.com/go/dialogflow/apiv2/dialogflowpb"
)

// handleMessageDialogflow handles a message from the platform, sends it to Dialogflow for intent detection,
// retrieves the corresponding context using RAG, generates a response with OpenAI, and sends it back to the user.
func (s *Service) HandleMessageDialogflow(platform bot.Platform, chatID, message string) (string, error) {

	// Detect intent using Dialogflow
	response, err := s.fetchDialogflowResponse(chatID, message)
	if err != nil {
		return "", fmt.Errorf("error detecting intent: %v", err)
	}

	intent := response.GetQueryResult().GetIntent().GetDisplayName()
	fmt.Printf("Detected intent: %s\n", intent)

	// Fetch document context
	context, err := s.fetchDocumentContext(intent, message)
	if err != nil {
		return "", fmt.Errorf("error fetching document context: %v", err)
	}

	// Generate response using OpenAI with context
	prompt := fmt.Sprintf("Context:\n%s\nUser query: %s", context, message)
	finalResponse, err := s.client.GetResponse(prompt)
	if err != nil {
		return "", fmt.Errorf("error generating OpenAI response: %v", err)
	}

	return finalResponse, nil
}

// fetchDialogflowResponse sends the message to Dialogflow and retrieves the response with detected intent.
func (s *Service) fetchDialogflowResponse(sessionID, text string) (*dialogflowpb.DetectIntentResponse, error) {
	conf := s.botConfig
	response, err := s.detectIntentText(conf.DialogflowProjectID, sessionID, text, "en")
	if err != nil {
		return nil, fmt.Errorf("error detecting intent with Dialogflow: %v", err)
	}
	return response, nil
}

// Send a text query to Dialogflow and returns the response
func (s *Service) detectIntentText(projectID, sessionID, text, languageCode string) (*dialogflowpb.DetectIntentResponse, error) {
	ctx := context.Background()
	client, err := dialogflow.NewSessionsClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating Dialogflow client: %v", err)
	}
	defer client.Close()

	sessionPath := fmt.Sprintf("projects/%s/agent/sessions/%s", projectID, sessionID)
	req := &dialogflowpb.DetectIntentRequest{
		Session: sessionPath,
		QueryInput: &dialogflowpb.QueryInput{
			Input: &dialogflowpb.QueryInput_Text{
				Text: &dialogflowpb.TextInput{
					Text:         text,
					LanguageCode: languageCode,
				},
			},
		},
	}
	return client.DetectIntent(ctx, req)
}

// fetchDocumentContext retrieves the document chunks based on the detected intent's associated tags.
func (s *Service) fetchDocumentContext(intent, userMessage string) (string, error) {
	tags := mapTags(intent)
	if len(tags) > 0 {
		return s.retrieveChunksByTags(tags)
	}
	return s.fallbackContext(userMessage)
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

// retrieveChunksByTags fetches document chunks that match the specified tags
func (s *Service) retrieveChunksByTags(tags []string) (string, error) {
	documentEmbeddings, err := s.repository.GetDocumentChunksByTags(tags)
	if err != nil {
		return "", fmt.Errorf("error retrieving document chunks: %v", err)
	}

	var contextChunks []string
	for _, chunk := range documentEmbeddings {
		contextChunks = append(contextChunks, chunk.DocText)
	}

	return strings.Join(contextChunks, "\n"), nil
}

// fallbackContext retrieves document chunks based on similarity to the user's message, functions as basic openAI mode.
func (s *Service) fallbackContext(userMessage string) (string, error) {
	documentEmbeddings, chunkText, err := s.repository.FetchEmbeddings()
	if err != nil {
		return "", fmt.Errorf("error fetching embeddings: %v", err)
	}

	// topChunks, err := document.RetrieveTopNChunks(userMessage, documentEmbeddings, 3, chunkText, 0.75)
	// if err != nil || len(topChunks) == 0 {
	// 	return "", fmt.Errorf("no relevant chunks found: %v", err)
	// }
	topChunks, err := document.RetrieveTopNChunks(userMessage, documentEmbeddings, 3, chunkText, 0.75)
	if err != nil || len(topChunks) == 0 {
		fmt.Printf("No relevant chunks found for message: %s\n", userMessage)
		return "", nil
	}

	context := strings.Join(topChunks, "\n")
	fmt.Printf("Retrieved fallback context: %s\n", context)
	return context, nil

	//return strings.Join(topChunks, "\n"), nil
}
