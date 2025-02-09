package bot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"crossplatform_chatbot/ai_clients"
	config "crossplatform_chatbot/configs"
	"crossplatform_chatbot/database"
	"crossplatform_chatbot/repository"
)

type IgBot interface {
	Run() error
}

type igBot struct {
	BaseBot
}

// creates a new IGBot instance
func NewIGBot(conf *config.BotConfig, embconf config.EmbeddingConfig, aiClients ai_clients.AIClients, database database.Database, dao repository.DAO) (*igBot, error) {
	// Verify that the page access token is available
	if conf.InstagramPageToken == "" {
		return nil, errors.New(" Instagram Page Access Token is not provided")
	}

	return &igBot{
		BaseBot: BaseBot{
			platform:  INSTAGRAM,
			conf:      conf,
			database:  database,
			dao:       dao,
			aiClients: aiClients,
			embConfig: embconf,
		},
	}, nil
}

// Run initializes and starts the Instagram bot with webhook
func (b *igBot) Run() error {
	if b.conf.InstagramPageToken == "" {
		return errors.New(" Instagram page access token is missing")
	}

	// Webhook confirmation (for Instagram API)
	fmt.Println("Instagram bot is running with webhook!")
	return nil
}

// The structure of incoming events from Instagram Messaging API
type InstagramEvent struct {
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

// sendResponse sends a message to the specified user on Instagram
func (b *igBot) SendReply(senderID interface{}, messageText string) error {

	//conf := config.GetConfig()
	url := fmt.Sprintf("https://graph.facebook.com/v17.0/me/messages?access_token=%s", b.conf.InstagramPageToken)

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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error response from Instagram API: %s", resp.Status)
	}

	log.Printf("Message sent successfully to %s", senderID.(string))
	return nil
}
