package service

import (
	document "crossplatform_chatbot/document_proc"
	"fmt"
	"log"
	"strings"
)

func (s *Service) processUserMessage(chatID, message, botTag string) (string, string, []string, []float64, error) { //TODO add username?
	fmt.Printf("Received message: %s from %s \n", message, botTag)
	fmt.Printf("Chat ID: %s\n", chatID)

	b := s.GetBot(botTag)
	baseBot := b.Base()

	// Get the corresponding platform from tag
	/*platform, err := getPlatformFromBotTag(botTag)
	if err != nil {
		return "Error finding bot tag.", "", nil, nil, err
	}*/

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

		// Fetch conversation history from Redis
		history, err := s.getConversationHistory(chatID, 5) // Retrieve the last 5 messages TODO
		if err != nil {
			log.Printf("Error retrieving conversation history: %v", err)
			history = "" // Default to no history
		}

		if s.botConfig.UseOpenAI {
			// Retrieve top relevant chunks.
			topChunks, err := document.RetrieveTopNChunks(message, documentEmbeddings, s.openaiClient, s.embConfig.NumTopChunks, chunkText, s.embConfig.ScoreThreshold)
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
				//prompt := fmt.Sprintf("Context:\n%s\nUser query: %s", context, message)
				prompt := fmt.Sprintf("Conversation history:\n%s\n\nContext:\n%s\nUser query: %s", history, context, message)

				response, err = baseBot.GetOpenAIResponse(prompt)
				if err != nil {
					return fmt.Sprintf("OpenAI Error: %v", err), "", nil, nil, err
				}
			} else {
				// Fallback to OpenAI response with history but without context.
				prompt := fmt.Sprintf("Conversation history:\n%s\nUser query: %s", history, message)
				response, err = baseBot.GetOpenAIResponse(prompt)
				if err != nil {
					return fmt.Sprintf("OpenAI Error: %v", err), "", nil, nil, err
				}
			}
		} else {
			// Fallback to dialogflow or another approach.
			//response, err = s.HandleMessageDialogflow(sessionID, message)

			response, intent, topChunkIDs, topChunkScores, err = s.handleMessageDialogflow(chatID, message, history) // sessionID passed down from outside
			if err != nil {
				return "Error processing with Dialogflow.", "", nil, nil, err
			}

			// Send response back to the platform
			// if err := b.SendResponse(sessionID, response); err != nil {
			// 	fmt.Printf("Error sending response: %v\n", err)
			// }

		}

	}

	err := s.saveConversation(chatID, message, response)
	if err != nil {
		return "Error saving to Redis.", "", nil, nil, err
	}

	return response, intent, topChunkIDs, topChunkScores, nil
}
