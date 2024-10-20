package handlers

import (
	"crossplatform_chatbot/bot"
	config "crossplatform_chatbot/configs"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/line/line-bot-sdk-go/linebot"
)

// HandleLineWebhook handles incoming POST requests from the Line platform
func HandleLineWebhook(c *gin.Context, lineBot bot.LineBot) {
	// Parse the incoming request from the Line platform and extract the events
	events, err := lineBot.ParseRequest(c.Request)
	if err != nil {
		// If the request has an invalid signature, return a 400 Bad Request error
		if err == linebot.ErrInvalidSignature {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid signature"})
			return
		}
		// If there is any other error during parsing, return a 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse request"})
		return
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
	c.Status(http.StatusOK)
}

// HandleTelegramWebhook handles incoming POST requests from Telegram
func HandleTelegramWebhook(c *gin.Context, tgBot bot.TgBot) {
	// Parse the incoming request from Telegram and extract the update
	var update tgbotapi.Update
	if err := c.ShouldBindJSON(&update); err != nil {
		// If there's an error parsing the request, return a 400 Bad Request error
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind request"})
		return
	}

	// Handle the update
	tgBot.HandleTelegramUpdate(update)
	c.Status(http.StatusOK)
}

// HandleMessengerWebhook handles incoming POST requests from Facebook Messenger
func HandleMessengerWebhook(c *gin.Context, fbBot bot.FbBot) {
	var event bot.MessengerEvent // Use the struct from the bot package

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse request"})
		return
	}

	for _, entry := range event.Entry {
		for _, msg := range entry.Messaging {
			senderID := msg.Sender.ID
			//messageText := msg.Message.Text

			//fbBot.HandleMessengerMessage(senderID, messageText)

			// Check if the message text is non-empty
			if strings.TrimSpace(msg.Message.Text) != "" {
				messageText := msg.Message.Text
				fbBot.HandleMessengerMessage(senderID, messageText)
			} else {
				fmt.Printf("Non-text or empty message received from %s, skipping...\n", senderID)
			}
		}
	}

	c.Status(http.StatusOK)
}

// VerifyMessengerWebhook verifies the webhook for Facebook Messenger (handles GET request)
func VerifyMessengerWebhook(c *gin.Context) {
	// Verify token from environment or configuration
	//verifyToken := os.Getenv("VERIFY_TOKEN")
	conf := config.GetConfig() //TODO: should config be loaded here?
	verifyToken := conf.FacebookVerifyToken

	// Check if the verify token matches
	if c.Query("hub.verify_token") == verifyToken {
		c.String(http.StatusOK, c.Query("hub.challenge"))
	} else {
		c.String(http.StatusForbidden, "Invalid verification token")
	}
}

// HandleInstagramWebhook handles incoming POST requests from Instagram
func HandleInstagramWebhook(c *gin.Context, igBot bot.IgBot) {
	var event bot.InstagramEvent // Use the struct from the bot package for Instagram

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse request"})
		return
	}

	for _, entry := range event.Entry {
		for _, msg := range entry.Messaging {
			senderID := msg.Sender.ID
			//senderID := msg.Recipient.ID

			// Check if the message text is non-empty
			if strings.TrimSpace(msg.Message.Text) != "" {
				messageText := msg.Message.Text
				igBot.HandleInstagramMessage(senderID, messageText)
			} else {
				fmt.Printf("Non-text or empty message received from %s, skipping...\n", senderID)
			}
		}
	}

	c.Status(http.StatusOK)
}

// VerifyInstagramWebhook verifies the webhook for Instagram Messaging (handles GET request)
func VerifyInstagramWebhook(c *gin.Context) {
	// Load verification token from configuration or environment
	conf := config.GetConfig()
	verifyToken := conf.InstagramVerifyToken // Use Instagram-specific verify token

	// Check if the verify token matches
	if c.Query("hub.verify_token") == verifyToken {
		c.String(http.StatusOK, c.Query("hub.challenge"))
	} else {
		c.String(http.StatusForbidden, "Invalid verification token")
	}
}
