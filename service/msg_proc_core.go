package service

import (
	document "crossplatform_chatbot/document_proc"
	"fmt"
	"strings"
)

func (s *Service) processUserMessage(chatID, message, botTag string) (string, string, []string, []float64, error) { //TODO add username?
	fmt.Printf("Received message: %s from %s \n", message, botTag)
	fmt.Printf("Chat ID: %s\n", chatID)

	b := s.GetBot(botTag)
	baseBot := b.Base()

	var response string
	var intent string
	var topChunkIDs []string
	var topChunkScores []float64

	if strings.HasPrefix(message, "/") {
		// Handle commands.
		response = baseBot.HandleCommand(message)
	} else if s.botConfig.Screaming && len(message) > 0 {
		// Example of simple transformation.
		response = strings.ToUpper(message)
	} else {
		// Fetch embeddings and process the message.
		documentEmbeddings, chunkText, err := s.repository.FetchEmbeddings()
		if err != nil {
			return "Error retrieving document embeddings.", "", nil, nil, err
		}

		if s.botConfig.UseOpenAI {
			// Retrieve top relevant chunks.
			topChunks, err := document.RetrieveTopNChunks(message, documentEmbeddings, s.embConfig.NumTopChunks, chunkText, s.embConfig.ScoreThreshold)
			if err != nil {
				return "Error retrieving related document information.", "", nil, nil, err
			}

			if len(topChunks) > 0 {
				// Extract chunk IDs and scores
				var contextBuilder []string
				for _, chunk := range topChunks {
					contextBuilder = append(contextBuilder, chunk.Text)
					topChunkIDs = append(topChunkIDs, chunk.ChunkID)
					topChunkScores = append(topChunkScores, chunk.Score)
				}

				// Use chunks as context for OpenAI.
				context := strings.Join(contextBuilder, "\n")
				prompt := fmt.Sprintf("Context:\n%s\nUser query: %s", context, message)

				response, err = baseBot.GetOpenAIResponse(prompt)
				if err != nil {
					return fmt.Sprintf("OpenAI Error: %v", err), "", nil, nil, err
				}
			} else {
				// Fallback to OpenAI response without context.
				response, err = baseBot.GetOpenAIResponse(message)
				if err != nil {
					return fmt.Sprintf("OpenAI Error: %v", err), "", nil, nil, err
				}
			}
		} else {
			// Fallback to dialogflow or another approach.
			//response, err = s.HandleMessageDialogflow(sessionID, message)

			response, intent, topChunkIDs, topChunkScores, err = s.handleMessageDialogflow(chatID, message) // sessionID passed down from outside
			if err != nil {
				return "Error processing with Dialogflow.", "", nil, nil, err
			}

			// Send response back to the platform
			// if err := b.SendResponse(sessionID, response); err != nil {
			// 	fmt.Printf("Error sending response: %v\n", err)
			// }

		}

	}

	return response, intent, topChunkIDs, topChunkScores, nil
}
