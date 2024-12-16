package service

import (
	"crossplatform_chatbot/bot"
	"crossplatform_chatbot/models"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/line/line-bot-sdk-go/linebot"
)

// HandleLine processes incoming requests from the LINE platform.
func (s *Service) HandleLine(req *http.Request) error {
	b := s.GetBot("line")
	lineBot, exist := b.(bot.LineBot)
	if !exist {
		return errors.New("line bot not found")
	}

	// Parse the incoming request from the Line platform and extract the events
	events, err := lineBot.ParseRequest(req)
	if err != nil {
		return err
	}

	// Loop through each event received from the Line platform
	for _, event := range events {
		// Check if the event is a message event
		if event.Type == linebot.EventTypeMessage {
			// Switch on the type of message (could be text, image, video, audio etc., only support text for now)
			switch message := event.Message.(type) {
			// If the message is a text message, process it using handleLineMessage
			case *linebot.TextMessage:

				// Retrieve and validate user profile
				userProfile, err := lineBot.GetUserProfile(event.Source.UserID)
				if err != nil {
					return fmt.Errorf("error fetching user profile: %v", err)

				}

				// Ensure user exists in the database
				userExists, err := lineBot.ValidateUser(userProfile, event.Source.UserID)
				if err != nil {
					return fmt.Errorf("error validating user: %v", err)
				}

				// userExists, err := b.validateAndGenerateToken(userProfile, event, event.Source.UserID)
				// if err != nil {
				// 	return fmt.Printf("Error ensuring user exists: %v\n", err)
				// }

				// If user didn't exist, a welcome message was sent, so return
				if !userExists {
					return fmt.Errorf("user does not exist in the database")
				}

				//lineBot.HandleLineMessage(event, message)
				//userID := event.Source.UserID
				chatID, err := s.getChatID(bot.LINE, event)
				if err != nil {
					fmt.Printf("Error getting chat ID: %v\n", err)
					return fmt.Errorf("error getting chat ID: %v", err)
				}
				response, err := s.processUserMessage(chatID, message.Text, "line")
				if err != nil {
					return fmt.Errorf("error processing user message: %w", err)
				}
				err = b.SendReply(event, response)
				if err != nil {
					return fmt.Errorf("error occurred while sending the response: %s", err.Error())
				}
			}
		}
	}

	return nil
}

// HandleTelegram processes incoming updates from Telegram, including documents and messages.
func (s *Service) HandleTelegram(update tgbotapi.Update) error {
	b := s.GetBot("telegram")
	tgBot, exists := b.(bot.TgBot)
	if !exists {
		return errors.New(" Telegram bot not found")
	}

	//tgBot.HandleTelegramUpdate(update)

	chatID := strconv.FormatInt(update.Message.Chat.ID, 10)

	if update.Message != nil {
		if update.Message.Document != nil {

			// get filename, fileURL, fileID
			fileID, fileURL, filename, err := tgBot.GetDocFile(update)
			if err != nil {
				return fmt.Errorf("error getting file:  %w", err)
			}

			// If the message contains a document, handle the document upload
			err = s.HandleDocumentUpload(filename, fileID, fileURL)
			if err != nil {
				b.SendReply(update.Message, "Error handling document: "+err.Error())
				return fmt.Errorf("error handling the document:  %w", err)
			}

			//b.SendTelegramMessage(update.Message.Chat.ID, "Document processed and stored in chunks for future queries.")

		} else {
			// Otherwise, handle regular text messages
			user := update.Message.From
			if user == nil {
				return fmt.Errorf("error getting user")
			}

			// token, err := b.validateAndGenerateToken(userIDStr, user, message)
			userExists, err := tgBot.ValidateUser(user, update.Message)
			if err != nil {
				return fmt.Errorf("error validating user: %s", err.Error())
			}

			if !userExists {
				// User was just created and welcomed, no further action needed.
				return nil
			}

			//tgBot.HandleTgMessage(update.Message)
			// Process the message and generate a response using the service layer.
			response, err := s.processUserMessage(chatID, update.Message.Text, "telegram")
			if err != nil {
				return fmt.Errorf("error processing user message: %w", err)
			}

			err = b.SendReply(update.Message, response)
			if err != nil {
				return fmt.Errorf("error occurred while sending the response: %s", err.Error())
			}
		}
	}

	return nil
}

// HandleMessenger processes incoming events from Facebook Messenger.
func (s *Service) HandleMessenger(event bot.MessengerEvent) error {
	b := s.GetBot("facebook")
	// fbBot, exists := b.(bot.FbBot)
	// if !exists {
	// 	return errors.New(" Messenger bot not found")
	// }

	for _, entry := range event.Entry {
		for _, msg := range entry.Messaging {
			senderID := msg.Sender.ID
			if messageText := strings.TrimSpace(msg.Message.Text); messageText != "" {
				//fbBot.HandleMessengerMessage(senderID, messageText)
				response, err := s.processUserMessage(senderID, messageText, "facebook")
				if err != nil {
					return fmt.Errorf("error processing user message: %w", err)
				}
				fmt.Printf("Sent message %s \n", response)
				err = b.SendReply(senderID, response)
				if err != nil {
					//c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while sending the response"})
					return fmt.Errorf("error occurred while sending the response: %s", err.Error())
				}

			}
		}
	}
	return nil
}

// HandleInstagram processes incoming events from Instagram.
func (s *Service) HandleInstagram(event bot.InstagramEvent) error {
	b := s.GetBot("instagram")
	// igBot, exists := s.GetBot("instagram").(bot.IgBot)
	// if !exists {
	// 	return errors.New(" Instagram bot not found")
	// }

	for _, entry := range event.Entry {
		for _, msg := range entry.Messaging {
			senderID := msg.Sender.ID
			if messageText := strings.TrimSpace(msg.Message.Text); messageText != "" {
				//igBot.HandleInstagramMessage(senderID, messageText)
				response, err := s.processUserMessage(senderID, messageText, "facebook")
				if err != nil {
					return fmt.Errorf("error processing user message: %w", err)
				}
				fmt.Printf("Sent message %s \n", response)
				err = b.SendReply(senderID, response)
				if err != nil {
					//c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while sending the response"})
					return fmt.Errorf("error occurred while sending the response: %s", err.Error())
				}
			}
		}
	}
	return nil
}

// HandleGeneral processes requests from the frontend for the general bot.
func (s *Service) HandleGeneral(req models.GeneralRequest) (string, error) {

	// Process the message and generate a response using the service layer.
	response, err := s.processUserMessage(req.SessionID, req.Message, "general")
	if err != nil {
		return "", fmt.Errorf("error processing user message: %w", err)
	}

	return response, nil
}

// getChatID returns a chat ID with a given platform
func (s *Service) getChatID(platform bot.Platform, identifier interface{}) (string, error) {
	switch platform {
	case bot.LINE:
		if event, ok := identifier.(*linebot.Event); ok {
			return event.Source.UserID, nil
		}
		return "", fmt.Errorf("invalid LINE event identifier")
	case bot.TELEGRAM:
		if message, ok := identifier.(*tgbotapi.Message); ok {
			return strconv.FormatInt(message.Chat.ID, 10), nil
		}
		return "", fmt.Errorf("invalid Telegram message identifier")
	case bot.FACEBOOK:
		if recipientID, ok := identifier.(string); ok {
			return recipientID, nil
		}
		return "", fmt.Errorf("invalid Messenger recipient identifier")
	case bot.GENERAL:
		if sessionID, ok := identifier.(string); ok {
			return sessionID, nil
		}
		return "", fmt.Errorf("invalid session identifier")
	default:
		return "", fmt.Errorf("unsupported platform")
	}
}

// func (s *Service) GetBotPlatform(botTag string) (bot.Bot, bot.Platform, error) {
// 	bot := s.GetBot(botTag)
// 	if bot == nil {
// 		return nil, 0, fmt.Errorf("bot not found for tag: %s", botTag)
// 	}
// 	return bot, bot.Platform(), nil
// }
