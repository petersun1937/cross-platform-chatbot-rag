package service

import (
	"crossplatform_chatbot/bot"
	document "crossplatform_chatbot/document_proc"
	"fmt"
	"strings"
)

func (s *Service) processUserMessage(chatID, message, botTag string) (string, error) { //TODO add username?
	fmt.Printf("Received message: %s from %s \n", message, botTag)
	fmt.Printf("Chat ID: %s\n", chatID)

	b := s.GetBot(botTag)
	baseBot := b.Base()

	// Get the corresponding platform from tag
	platform, err := getPlatformFromBotTag(botTag)
	if err != nil {
		return "Error finding bot tag.", err
	}

	var response string

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
			return "Error retrieving document embeddings.", err
		}

		if s.botConfig.UseOpenAI {
			// Retrieve top relevant chunks.
			topChunks, err := document.RetrieveTopNChunks(message, documentEmbeddings, s.embConfig.NumTopChunks, chunkText, s.embConfig.ScoreThreshold)
			if err != nil {
				return "Error retrieving related document information.", err
			}

			if len(topChunks) > 0 {
				// Use chunks as context for OpenAI.
				context := strings.Join(topChunks, "\n")
				prompt := fmt.Sprintf("Context:\n%s\nUser query: %s", context, message)
				response, err = baseBot.GetOpenAIResponse(prompt)
				if err != nil {
					return fmt.Sprintf("OpenAI Error: %v", err), nil
				}
			} else {
				// Fallback to OpenAI response without context.
				response, err = baseBot.GetOpenAIResponse(message)
				if err != nil {
					return fmt.Sprintf("OpenAI Error: %v", err), nil
				}
			}
		} else {
			// Fallback to dialogflow or another approach.
			//response, err = s.HandleMessageDialogflow(sessionID, message)

			response, err = s.HandleMessageDialogflow(platform, chatID, message) // sessionID passed down from outside
			if err != nil {
				return "Error processing with Dialogflow.", err
			}

			// Send response back to the platform
			// if err := b.SendResponse(sessionID, response); err != nil {
			// 	fmt.Printf("Error sending response: %v\n", err)
			// }

		}

	}

	return response, nil
}

func getPlatformFromBotTag(botTag string) (bot.Platform, error) {
	switch botTag {
	case "line":
		return bot.LINE, nil
	case "telegram":
		return bot.TELEGRAM, nil
	case "facebook":
		return bot.FACEBOOK, nil
	case "instagram":
		return bot.INSTAGRAM, nil
	case "general":
		return bot.GENERAL, nil
	default:
		return 0, fmt.Errorf("unknown bot tag: %s", botTag)
	}
}
