package service

import (
	"crossplatform_chatbot/bot"
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

		if !s.botConfig.UseDialogflow {
			// Retrieve top relevant chunks.
			topChunks, err := document.RetrieveTopNChunks(message, documentEmbeddings, s.aiClients.OpenAI, s.embConfig.NumTopChunks, chunkText, s.embConfig.ScoreThreshold)
			if err != nil {
				return "Error retrieving related document information.", "", nil, nil, err
			}
			prompt := ""

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
				prompt = fmt.Sprintf("Conversation history:\n%s\n\nContext:\n%s\nUser query: %s", history, context, message)

			} else {
				// Fallback to OpenAI response with history but without context.
				prompt = fmt.Sprintf("Conversation history:\n%s\nUser query: %s", history, message)
			}
			//response, err = baseBot.GetOpenAIResponse(prompt)
			response, err = s.generateResponse(prompt, baseBot)
			if err != nil {
				return fmt.Sprintf("Error: %v", err), "", nil, nil, err
			}
		} else {
			// Fallback to dialogflow or another approach.
			//response, err = s.HandleMessageDialogflow(sessionID, message)
			prompt := ""
			prompt, intent, topChunkIDs, topChunkScores, err = s.handleMessageDialogflow(chatID, message, history)
			if err != nil {
				return "Error processing with Dialogflow.", "", nil, nil, err
			}

			//response, err = baseBot.GetOpenAIResponse(prompt)
			response, err = s.generateResponse(prompt, baseBot)
			if err != nil {
				return "", "", nil, nil, fmt.Errorf("error generating response: %v", err)
			}
		}

	}

	err := s.saveConversation(chatID, message, response)
	if err != nil {
		return "Error saving to Redis.", "", nil, nil, err
	}

	return response, intent, topChunkIDs, topChunkScores, nil
}

func (s *Service) generateResponse(prompt string, b *bot.BaseBot) (string, error) {
	if s.botConfig.UseOpenAI {
		return b.GetOpenAIResponse(prompt)
		//return s.aiClients.OpenAI.GetResponse(prompt)
	} else if s.botConfig.UseMistral {
		return b.GetMistralResponse(prompt)
		//return s.aiClients.Mistral.GetResponse(prompt)
	} else if s.botConfig.UseMETA {
		return b.GetTogetherAIResponse(prompt)
		//return s.aiClients.TogetherAI.GetResponse(prompt)
	} /*else if s.botConfig.UseHuggingFace {
		return s.aiClients.HuggingFace.GetResponse(prompt)
	}*/
	return "", fmt.Errorf("error: No AI provider is enabled in the configuration")
}
