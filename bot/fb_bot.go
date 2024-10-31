package bot

import (
	"bytes"
	config "crossplatform_chatbot/configs"
	"crossplatform_chatbot/database"
	"crossplatform_chatbot/document"
	"crossplatform_chatbot/repository"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"cloud.google.com/go/dialogflow/apiv2/dialogflowpb"
)

type FbBot interface {
	HandleMessengerMessage(senderID, messageText string)
	Run() error
}

type fbBot struct {
	BaseBot
	//conf            config.BotConfig
	//ctx             context.Context
	//pageAccessToken string
}

// creates a new FbBot instance
func NewFBBot(conf config.BotConfig, database database.Database, dao repository.DAO) (*fbBot, error) {
	// Verify that the page access token is available
	if conf.FacebookPageToken == "" {
		return nil, errors.New("facebook Page Access Token is not provided")
	}

	return &fbBot{
		BaseBot: BaseBot{
			Platform: FACEBOOK,
			conf:     conf,
			database: database,
			dao:      dao,
		},
	}, nil
}

// // creates a new FbBot instance
// func NewFBBot(conf *config.Config, service *service.Service) (*fbBot, error) {
// 	// Verify that the page access token is available
// 	if conf.FacebookPageToken == "" {
// 		return nil, errors.New("facebook Page Access Token is not provided")
// 	}

// 	// Initialize the BaseBot structure
// 	baseBot := &BaseBot{
// 		Platform: FACEBOOK,
// 		Service:  service,
// 	}

// 	// Initialize and return the fbBot instance
// 	return &fbBot{
// 		BaseBot:         baseBot,
// 		conf:            conf.BotConfig,
// 		ctx:             context.Background(),
// 		pageAccessToken: conf.FacebookPageToken,
// 	}, nil
// }

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

// HandleMessengerMessage processes incoming messages and sends a response
func (b *fbBot) HandleMessengerMessage(senderID, messageText string) {
	// Trim whitespace and check for empty message
	if strings.TrimSpace(messageText) == "" {
		fmt.Printf("Empty or invalid message received from %s, ignoring...\n", senderID)
		return
	}

	// // Validate user and generate a token if necessary
	// token, err := b.validateAndGenerateToken(senderID)
	// if err != nil {
	// 	log.Printf("Error validating user: %s", err.Error())
	// 	return
	// }

	// // If a token is generated for a new user, send it to them
	// if token != nil {
	// 	b.sendResponse(senderID, "Welcome! Your access token is: "+*token)
	// } else {
	// Process the user's message if no token is sent
	b.processUserMessage(senderID, messageText)
	//}
}

// validateAndGenerateToken checks if the user exists and generates a token if not
// func (b *fbBot) validateAndGenerateToken(userID string) (*string, error) {
// 	// Retrieve user profile information from Facebook
// 	userProfile, err := b.getUserProfile(userID)
// 	if err != nil {
// 		return nil, fmt.Errorf("error fetching user profile: %w", err)
// 	}

// 	// Check if the user exists in the database
// 	var dbUser models.User
// 	err = b.Service.GetDB().Where("user_id = ? AND deleted_at IS NULL", userID).First(&dbUser).Error
// 	if err != nil {
// 		// If user does not exist, create a new user
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			dbUser = models.User{
// 				UserID:       userID,
// 				UserName:     userProfile.FirstName + " " + userProfile.LastName, // Combine first and last name
// 				FirstName:    userProfile.FirstName,
// 				LastName:     userProfile.LastName,
// 				LanguageCode: "", // Facebook doesn't provide language directly
// 			}

// 			// Create the new user record in the database
// 			if err := b.Service.GetDB().Create(&dbUser).Error; err != nil {
// 				return nil, fmt.Errorf("error creating user: %w", err)
// 			}

// 			// Generate a JWT token using the service's ValidateUser method
// 			token, err := b.Service.ValidateUser(userID, service.ValidateUserReq{
// 				FirstName:    userProfile.FirstName,
// 				LastName:     userProfile.LastName,
// 				UserName:     "", // Facebook doesn’t provide username directly
// 				LanguageCode: "", // Facebook doesn’t provide language directly
// 			})
// 			if err != nil {
// 				return nil, fmt.Errorf("error generating JWT: %w", err)
// 			}

// 			return token, nil // Return the generated token
// 		}
// 		return nil, fmt.Errorf("error retrieving user: %w", err)
// 	}

// 	return nil, nil // User already exists, no token generation needed
// }

// // getUserProfile retrieves the user profile information from Facebook
// func (b *fbBot) getUserProfile(userID string) (*service.UserProfile, error) {
// 	url := fmt.Sprintf("https://graph.facebook.com/%s?fields=first_name,last_name&access_token=%s", userID, b.conf.FacebookPageToken) // TODO move url to env

// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, fmt.Errorf("error fetching user profile: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("invalid response from Facebook, status code: %d", resp.StatusCode)
// 	}

// 	var profile service.UserProfile
// 	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
// 		return nil, fmt.Errorf("error decoding profile response: %w", err)
// 	}

// 	return &profile, nil
// }

// processUserMessage processes the user message and responds accordingly
func (b *fbBot) processUserMessage(senderID, text string) {
	fmt.Printf("Received message from %s: %s \n", senderID, text)

	var response string
	var err error

	// Check if the message is a command (starts with "/")
	if strings.HasPrefix(text, "/") {
		response = handleCommand(text)
		/*response, err = handleCommand(senderID, text, b)
		if err != nil {
			fmt.Printf("An error occurred: %s \n", err.Error())
			response = "An error occurred while processing your command."
		}*/
	} else if screaming && len(text) > 0 {
		// Check for a "screaming" mode if applicable (uppercase response)
		response = strings.ToUpper(text)
	} else {
		// Fetch document embeddings and try to match based on similarity
		documentEmbeddings, chunkText, err := b.BaseBot.dao.FetchEmbeddings()
		//documentEmbeddings, chunkText, err := b.Service.GetAllDocumentEmbeddings()
		if err != nil {
			fmt.Printf("Error retrieving document embeddings: %v", err)
			response = "Error retrieving document embeddings."
		} else if useOpenAI {
			// Perform similarity matching with the user's message
			topChunks, err := document.RetrieveTopNChunks(text, documentEmbeddings, 10, chunkText, 0.7) // Retrieve top 3 relevant chunks thresholded by score of 0.7
			if err != nil {
				fmt.Printf("Error retrieving document chunks: %v", err)
				response = "Error retrieving related document information."
			} else if len(topChunks) > 0 {
				// If there are similar chunks found, provide them as context for GPT
				context := strings.Join(topChunks, "\n")
				gptPrompt := fmt.Sprintf("Context:\n%s\nUser query: %s", context, text)

				// Call GPT with the context and user query
				response, err = GetOpenAIResponse(gptPrompt)
				if err != nil {
					response = fmt.Sprintf("OpenAI Error: %v", err)
				}
			} else {
				// If no relevant document found, fallback to OpenAI response
				response, err = GetOpenAIResponse(text)
				if err != nil {
					response = fmt.Sprintf("OpenAI Error: %v", err)
				}
			}
		} else {
			// Use Dialogflow if OpenAI is not enabled
			handleMessageDialogflow(FACEBOOK, senderID, text, b)
			return
		}
	}

	// Send the response if it's not empty
	if response != "" {
		err = b.sendResponse(senderID, response)
		if err != nil {
			fmt.Printf("Error sending response: %s\n", err.Error())
		}
	}

}

// handleDialogflowResponse processes and sends the Dialogflow response to the appropriate platform
func (b *fbBot) handleDialogflowResponse(response *dialogflowpb.DetectIntentResponse, identifier interface{}) error {

	// Check if the ID (identifier) is a string (which would be the sender ID for Facebook)
	_, ok := identifier.(string)
	if !ok {
		return fmt.Errorf("invalid Facebook message identifier")
	}

	// Iterate over the fulfillment messages returned by Dialogflow
	for _, msg := range response.QueryResult.FulfillmentMessages {
		if text := msg.GetText(); text != nil {
			// Send the response message to the user on Facebook Messenger
			return b.sendResponse(identifier, text.Text[0])
		}
	}

	return nil
}

// sendMessage sends a message to the specified user on Messenger
func (b *fbBot) sendResponse(senderID interface{}, messageText string) error {
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

/*
func (b *fbBot) sendMenu(identifier interface{}) error {
	if senderID, ok := identifier.(string); ok {
		return b.sendMessengerMenu(senderID)
	} else {
		return fmt.Errorf("invalid identifier type for LINE platform")
	}
}

// sendMessengerMenu sends a menu with buttons to the user
func (b *fbBot) sendMessengerMenu(senderID string) error {
	// Define the URL for the Facebook Graph API
	url := "https://graph.facebook.com/v12.0/me/messages?access_token=" + b.pageAccessToken

	// Define the menu with buttons
	menuPayload := map[string]interface{}{
		"recipient": map[string]interface{}{
			"id": senderID,
		},
		"message": map[string]interface{}{
			"attachment": map[string]interface{}{
				"type": "template",
				"payload": map[string]interface{}{
					"template_type": "button",
					"text":          "Choose an option:",
					"buttons": []map[string]interface{}{
						{
							"type":    "postback",
							"title":   "Menu 1",
							"payload": "MENU_1_PAYLOAD",
						},
						{
							"type":    "postback",
							"title":   "Menu 2",
							"payload": "MENU_2_PAYLOAD",
						},
						{
							"type":  "web_url",
							"title": "Visit Website",
							"url":   "https://www.example.com",
						},
					},
				},
			},
		},
	}

	// Marshal the payload to JSON
	menuBody, err := json.Marshal(menuPayload)
	if err != nil {
		return fmt.Errorf("error marshaling menu payload: %w", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(menuBody))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Send the request using the HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending menu request: %w", err)
	}
	defer resp.Body.Close()

	fmt.Printf("Menu sent successfully to %s\n", senderID)
	return nil
}*/
