package bot

import (
	"bytes"
	"crossplatform_chatbot/ai_clients"
	config "crossplatform_chatbot/configs"
	"crossplatform_chatbot/database"
	"crossplatform_chatbot/repository"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type FbBot interface {
	Run() error
}

type fbBot struct {
	BaseBot
}

// creates a new FbBot instance
func NewFBBot(conf *config.BotConfig, embconf config.EmbeddingConfig, aiClients ai_clients.AIClients, database database.Database, dao repository.DAO) (*fbBot, error) {
	// Verify that the page access token is available
	if conf.FacebookPageToken == "" {
		return nil, errors.New("facebook Page Access Token is not provided")
	}

	return &fbBot{
		BaseBot: BaseBot{
			platform:  FACEBOOK,
			conf:      conf,
			database:  database,
			dao:       dao,
			aiClients: aiClients,
			//openAIclient:  openaiClient,
			//mistralClient: mistralClient,
			embConfig: embconf,
		},
	}, nil
}

// Run initializes and starts the Facebook bot with webhook
func (b *fbBot) Run() error {
	if b.conf.FacebookPageToken == "" {
		return errors.New("page access token is missing")
	}

	// No bot instance for Messenger

	// webhook confirmation
	fmt.Println("Facebook Messenger bot is running with webhook!")

	return nil
}

// MessengerEvent defines the structure of incoming events from Facebook Messenger
type MessengerEvent struct {
	Object string `json:"object"`
	Entry  []struct {
		ID        string `json:"id"`
		Time      int64  `json:"time"`
		Messaging []struct {
			Sender struct {
				ID string `json:"id"`
			} `json:"sender"`
			Recipient struct {
				ID string `json:"id"`
			} `json:"recipient"`
			Timestamp int64 `json:"timestamp"`
			Message   struct {
				Mid  string `json:"mid"`
				Text string `json:"text"`
			} `json:"message"`
		} `json:"messaging"`
	} `json:"entry"`
}

// sendMessage sends a message to the specified user on Messenger
func (b *fbBot) SendReply(senderID interface{}, messageText string) error {
	//conf := config.GetConfig()
	url := b.conf.FacebookAPIURL + "/messages?access_token=" + b.conf.FacebookPageToken

	// Create the message payload
	messageData := map[string]interface{}{
		"recipient": map[string]string{"id": senderID.(string)},
		"message":   map[string]string{"text": messageText},
	}

	// Marshal the payload to JSON
	messageBody, err := json.Marshal(messageData)
	if err != nil {
		return fmt.Errorf("error marshaling message: %w", err)
	}

	// Create HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(messageBody))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// Set Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending response: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("Message sent successfully to %s", senderID)
	return nil
}
