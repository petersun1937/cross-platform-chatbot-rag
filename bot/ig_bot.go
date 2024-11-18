package bot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	config "crossplatform_chatbot/configs"
	"crossplatform_chatbot/database"
	document "crossplatform_chatbot/document_proc"
	openai "crossplatform_chatbot/openai"
	"crossplatform_chatbot/repository"

	"cloud.google.com/go/dialogflow/apiv2/dialogflowpb"
)

type IgBot interface {
	Run() error
	//HandleInstagramWebhook(c *gin.Context, igBot IgBot)
	HandleInstagramMessage(senderID, messageText string)
	//sendResponse(identifier interface{}, response string) error
	//setWebhook(webhookURL string) error
}

type igBot struct {
	BaseBot
	// conf config.BotConfig
	// ctx             context.Context
	// pageAccessToken string
	//openAIclient    *openai.Client
}

// creates a new IGBot instance
func NewIGBot(conf config.BotConfig, database database.Database, embconf config.EmbeddingConfig, dao repository.DAO) (*igBot, error) {
	// Verify that the page access token is available
	if conf.InstagramPageToken == "" {
		return nil, errors.New(" Instagram Page Access Token is not provided")
	}

	return &igBot{
		BaseBot: BaseBot{
			Platform:     INSTAGRAM,
			conf:         conf,
			database:     database,
			dao:          dao,
			openAIclient: openai.NewClient(),
			embConfig:    embconf,
		},
	}, nil
}

// // creates a new Instagram bot instance
// func NewIGBot(conf *config.Config, service *service.Service) (*igBot, error) {
// 	// Verify that the page access token is available
// 	if conf.InstagramPageToken == "" {
// 		return nil, errors.New(" Instagram Page Access Token is not provided")
// 	}

// 	// Initialize the BaseBot structure
// 	baseBot := &BaseBot{
// 		Platform: INSTAGRAM,
// 		Service:  service,
// 	}

// 	// Initialize and return the IgBot instance
// 	return &igBot{
// 		BaseBot:         baseBot,
// 		conf:            conf.BotConfig,
// 		ctx:             context.Background(),
// 		pageAccessToken: conf.InstagramPageToken,
// 		//openAIclient    *openai.Client,
// 	}, nil
// }

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

func (b *igBot) HandleInstagramMessage(senderID, messageText string) {
	// Trim whitespace and check for empty message
	if strings.TrimSpace(messageText) == "" {
		fmt.Printf("Empty message received from user %s, ignoring...\n", senderID)
		return
	}

	// Log the message
	log.Printf("Instagram message received from %s: %s\n", senderID, messageText)

	// Validate user and generate a token if necessary
	/*token, err := b.validateAndGenerateToken(senderID)
	if err != nil {
		log.Printf("Error validating user: %s", err.Error())
		return
	}*/

	// If a token is generated for a new user, send it to them
	//if token != nil {
	//	b.sendResponse(senderID, "Welcome! Your access token is: "+*token)
	//} else {
	// Process the user's message if no token is sent
	b.processUserMessage(senderID, messageText)
	//}
}

// // validateAndGenerateToken checks if the user exists and generates a token if not
// func (b *igBot) validateAndGenerateToken(userID string) (*string, error) {
// 	// Retrieve user profile information from Instagram (similar to Facebook logic)
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
// 				LanguageCode: "", // Instagram doesn't provide language directly
// 			}

// 			// Create the new user record in the database
// 			if err := b.Service.GetDB().Create(&dbUser).Error; err != nil {
// 				return nil, fmt.Errorf("error creating user: %w", err)
// 			}

// 			// Generate a JWT token using the service's ValidateUser method
// 			token, err := b.Service.ValidateUser(userID, service.ValidateUserReq{
// 				FirstName:    userProfile.FirstName,
// 				LastName:     userProfile.LastName,
// 				UserName:     "", // Instagram doesn’t provide username directly
// 				LanguageCode: "", // Instagram doesn’t provide language directly
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

// // getUserProfile retrieves the user profile information from Instagram (similar to Facebook)
// func (b *igBot) getUserProfile(userID string) (*service.UserProfile, error) {
// 	// Use Instagram API to fetch user profile details
// 	url := fmt.Sprintf("https://graph.instagram.com/%s?fields=first_name,last_name&access_token=%s", userID, b.pageAccessToken)

// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, fmt.Errorf("error fetching user profile: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	/*if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("invalid response from Instagram, status code: %d", resp.StatusCode)
// 	}*/

// 	if resp.StatusCode != http.StatusOK {
// 		// Print the response body for debugging
// 		bodyBytes, _ := io.ReadAll(resp.Body)
// 		bodyString := string(bodyBytes)
// 		return nil, fmt.Errorf("invalid response from Instagram, status code: %d, response: %s", resp.StatusCode, bodyString)
// 	}

// 	var profile service.UserProfile
// 	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
// 		return nil, fmt.Errorf("error decoding profile response: %w", err)
// 	}

// 	return &profile, nil
// }

// sendResponse sends a message to the specified user on Instagram
func (b *igBot) sendResponse(senderID interface{}, messageText string) error {

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

// To processes incoming attachments (e.g., images or videos)
/*func (b *igBot) HandleInstagramAttachment(senderID, attachmentURL string) {
	log.Printf("Received attachment from Instagram user %s: %s\n", senderID, attachmentURL)

	// You can add custom logic to handle the attachment, for now, just send a response
	responseText := fmt.Sprintf("Thanks for the attachment! You sent: %s", attachmentURL)
	err := b.sendResponse(senderID, responseText)
	if err != nil {
		log.Printf("Error sending attachment response: %v\n", err)
	}
}*/

// handleDialogflowResponse processes and sends the Dialogflow response to the Instagram platform
func (b *igBot) handleDialogflowResponse(response *dialogflowpb.DetectIntentResponse, identifier interface{}) error {

	// Check if the senderID (identifier) is a string (which would be the sender ID for Instagram)
	senderID, ok := identifier.(string)
	if !ok {
		return fmt.Errorf("invalid Instagram message identifier")
	}

	// Iterate over the fulfillment messages returned by Dialogflow
	for _, msg := range response.QueryResult.FulfillmentMessages {
		if text := msg.GetText(); text != nil {
			// Send the response message to the user on Instagram
			return b.sendResponse(senderID, text.Text[0])
		}
	}

	return nil
}

// processUserMessage processes the user's message and responds accordingly
func (b *igBot) processUserMessage(senderID, text string) {
	fmt.Printf("Received message from %s: %s \n", senderID, text)

	var response string
	var err error

	// Check if the message is a command (starts with "/")
	if strings.HasPrefix(text, "/") {
		response = handleCommand(text)
	} else if screaming && len(text) > 0 {
		// Handle "screaming" mode, where responses are in uppercase
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
				response, err = b.BaseBot.GetOpenAIResponse(gptPrompt)
				if err != nil {
					response = fmt.Sprintf("OpenAI Error: %v", err)
				}
			} else {
				// If no relevant document found, fallback to OpenAI response
				response, err = b.BaseBot.GetOpenAIResponse(text)
				if err != nil {
					response = fmt.Sprintf("OpenAI Error: %v", err)
				}
			}
		} else {
			// Use Dialogflow if OpenAI is not enabled
			b.BaseBot.handleMessageDialogflow(INSTAGRAM, senderID, text, b)
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
