package service

import (
	"crossplatform_chatbot/bot"
	"crossplatform_chatbot/models"
	"errors"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/line/line-bot-sdk-go/linebot"
)

// HandleLine processes incoming requests from the LINE platform.
func (s *Service) HandleLine(req *http.Request) error {
	lineBot, exist := s.GetBot("line").(bot.LineBot)
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
				lineBot.HandleLineMessage(event, message)
			}
		}
	}

	return nil
}

// HandleTelegram processes incoming updates from Telegram.
func (s *Service) HandleTelegram(update tgbotapi.Update) error {
	tgBot, exists := s.GetBot("telegram").(bot.TgBot)
	if !exists {
		return errors.New(" Telegram bot not found")
	}

	tgBot.HandleTelegramUpdate(update)
	return nil
}

// HandleMessenger processes incoming events from Facebook Messenger.
func (s *Service) HandleMessenger(event bot.MessengerEvent) error {
	fbBot, exists := s.GetBot("facebook").(bot.FbBot)
	if !exists {
		return errors.New(" Messenger bot not found")
	}

	for _, entry := range event.Entry {
		for _, msg := range entry.Messaging {
			senderID := msg.Sender.ID
			if messageText := strings.TrimSpace(msg.Message.Text); messageText != "" {
				fbBot.HandleMessengerMessage(senderID, messageText)
			}
		}
	}
	return nil
}

// HandleInstagram processes incoming events from Instagram.
func (s *Service) HandleInstagram(event bot.InstagramEvent) error {
	igBot, exists := s.GetBot("instagram").(bot.IgBot)
	if !exists {
		return errors.New(" Instagram bot not found")
	}

	for _, entry := range event.Entry {
		for _, msg := range entry.Messaging {
			senderID := msg.Sender.ID
			if messageText := strings.TrimSpace(msg.Message.Text); messageText != "" {
				igBot.HandleInstagramMessage(senderID, messageText)
			}
		}
	}
	return nil
}

// HandleGeneral processes requests from the frontend for the general bot.
func (s *Service) HandleGeneral(req models.GeneralRequest) error {

	genBot, exists := s.GetBot("general").(bot.GeneralBot)
	if !exists {
		return errors.New("general bot not found")
	}

	// Delegate the handling of the message to the general bot.
	genBot.HandleGeneralMessage(req.SessionID, req.Message)

	return nil
}
