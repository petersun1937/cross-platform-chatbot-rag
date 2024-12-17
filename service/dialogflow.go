package service

import (
	"context"
	"crossplatform_chatbot/utils"
	"fmt"
	"strings"

	document "crossplatform_chatbot/document_proc"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"cloud.google.com/go/dialogflow/apiv2/dialogflowpb"
)

// handleMessageDialogflow handles a message from the platform, sends it to Dialogflow for intent detection,
// retrieves the corresponding context using RAG, generates a response with OpenAI, and sends it back to the user.
func (s *Service) handleMessageDialogflow(chatID, message string) (string, string, []string, []float64, error) {
	//func (s *Service) HandleMessageDialogflow(platform bot.Platform, identifier interface{}, message string) (string, error) {
	//b := s.GetBot(botTag)
	//baseBot := b.Base()

	// Determine chat ID
	/*chatID, err := s.getSessionID(platform, identifier)
	if err != nil {
		fmt.Printf("Error getting session ID: %v\n", err)
		return "", fmt.Errorf("error getting session ID: %v", err)
	}*/

	// Detect intent using Dialogflow
	response, err := s.fetchDialogflowResponse(chatID, message)
	if err != nil {
		return "", "", nil, nil, fmt.Errorf("error detecting intent: %v", err)
	}

	intent := response.GetQueryResult().GetIntent().GetDisplayName()
	fmt.Printf("Detected intent: %s\n", intent)

	// Fetch document context
	context, topChunkIDs, topChunkScores, err := s.fetchDocumentContext(intent, message)
	if err != nil {
		return "", "", nil, nil, fmt.Errorf("error fetching document context: %v", err)
	}

	// Generate response using OpenAI with context
	prompt := fmt.Sprintf("Context:\n%s\nUser query: %s", context, message)
	finalResponse, err := s.client.GetResponse(prompt)
	if err != nil {
		return "", "", nil, nil, fmt.Errorf("error generating OpenAI response: %v", err)
	}

	return finalResponse, intent, topChunkIDs, topChunkScores, nil
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
func (s *Service) fetchDocumentContext(intent, userMessage string) (string, []string, []float64, error) {
	// Special case: Directly return an empty context for "Default Welcome Intent"
	if intent == "Default Welcome Intent" {
		return "", []string{}, []float64{}, nil
	}

	tags := mapTags(intent)
	if len(tags) > 0 {
		//topChunkIDs, topChunkScores, chunkScores, err := s.retrieveChunksByTags(tags, userMessage)
		return s.retrieveChunksByTags(tags, userMessage)
	}

	//topChunkIDs, topChunkScores, chunkScores, err := s.fallbackContext(userMessage)
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
func (s *Service) retrieveChunksByTags(tags []string, userMessage string) (string, []string, []float64, error) {
	documentChunks, err := s.repository.GetDocumentChunksByTags(tags)
	if err != nil {
		return "", nil, nil, fmt.Errorf("error retrieving document chunks: %v", err)
	}

	docIDToText := make(map[string]string)
	documentEmbeddings := make(map[string][]float64) // Map[ChunkID] -> Embeddin

	for _, chunk := range documentChunks {
		// Convert embedding from string to []float64
		floatSlice, err := utils.PostgresArrayToFloat64Slice(chunk.Embedding)
		if err != nil {
			fmt.Printf("Error parsing embedding for chunkID %s: %v\n", chunk.ChunkID, err)
			continue // Skip this chunk if there's an error
		}
		documentEmbeddings[chunk.ChunkID] = floatSlice
		docIDToText[chunk.ChunkID] = chunk.DocText
	}
	// Apply scoring using RetrieveTopNChunks
	topChunks, err := document.RetrieveTopNChunks(userMessage, documentEmbeddings, s.embConfig.NumTopChunks, docIDToText, s.embConfig.ScoreThreshold)
	if err != nil || len(topChunks) == 0 {
		fmt.Println("No relevant chunks found for the given tags.")
		return "", nil, nil, nil
	}

	/*var contextChunks []string
	for _, chunk := range documentEmbeddings {
		contextChunks = append(contextChunks, chunk.DocText)
	}*/

	// Prepare output: context, IDs, and scores
	var contextBuilder []string
	var topChunkIDs []string
	var topChunkScores []float64

	// Build context, separate IDs and scores TODO: build context later?
	for _, chunk := range topChunks {
		fmt.Printf("ChunkID: %s, Score: %.4f\n", chunk.ChunkID, chunk.Score)
		contextBuilder = append(contextBuilder, chunk.Text)
		topChunkIDs = append(topChunkIDs, chunk.ChunkID)
		topChunkScores = append(topChunkScores, chunk.Score)
	}

	context := strings.Join(contextBuilder, "\n")
	return context, topChunkIDs, topChunkScores, nil
}

// fallbackContext retrieves document chunks based on similarity to the user's message, functions as basic openAI mode.
func (s *Service) fallbackContext(userMessage string) (string, []string, []float64, error) {
	documentEmbeddings, chunkText, err := s.repository.FetchEmbeddings()
	if err != nil {
		return "", nil, nil, fmt.Errorf("error fetching embeddings: %v", err)
	}

	// topChunks, err := document.RetrieveTopNChunks(userMessage, documentEmbeddings, 3, chunkText, 0.75)
	// if err != nil || len(topChunks) == 0 {
	// 	return "", fmt.Errorf("no relevant chunks found: %v", err)
	// }
	topChunks, err := document.RetrieveTopNChunks(userMessage, documentEmbeddings, s.embConfig.NumTopChunks, chunkText, s.embConfig.ScoreThreshold)
	if err != nil || len(topChunks) == 0 {
		fmt.Printf("No relevant chunks found for message: %s\n", userMessage)
		return "", nil, nil, nil
	}

	// Prepare the context by joining the texts of the top chunks
	var contextBuilder []string
	var chunkIDs []string
	var chunkScores []float64

	for _, chunk := range topChunks {
		// Log the chunkID and score for debugging or analysis
		fmt.Printf("ChunkID: %s, Score: %.4f\n", chunk.ChunkID, chunk.Score)

		// Append the text of the chunk to the context
		contextBuilder = append(contextBuilder, chunk.Text)
		chunkIDs = append(chunkIDs, chunk.ChunkID)
		chunkScores = append(chunkScores, chunk.Score)
	}

	// Combine all top chunk texts into a single context string
	context := strings.Join(contextBuilder, "\n")

	//context := strings.Join(topChunks, "\n")
	fmt.Printf("Retrieved fallback context: %s\n", context)
	return context, chunkIDs, chunkScores, nil

	//return strings.Join(topChunks, "\n"), nil
}
