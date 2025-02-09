package bot

import (
	"bytes"
	"crossplatform_chatbot/ai_clients"
	config "crossplatform_chatbot/configs"
	"crossplatform_chatbot/database"
	"crossplatform_chatbot/models"
	"crossplatform_chatbot/repository"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

type TgBot interface {
	Run() error
	GetDocFile(update tgbotapi.Update) (string, string, string, error)
	ValidateUser(user *tgbotapi.User, message *tgbotapi.Message) (bool, error)
}

type tgBot struct {
	BaseBot
	botApi *tgbotapi.BotAPI
}

// creates a new TGBot instance
func NewTGBot(botconf *config.BotConfig, embconf config.EmbeddingConfig, aiClients ai_clients.AIClients, database database.Database, dao repository.DAO) (*tgBot, error) {
	// Attempt to create a new Telegram bot using the provided token
	botApi, err := tgbotapi.NewBotAPI(botconf.TelegramBotToken)
	if err != nil {
		return nil, err
	}
	// Ensure botApi is not nil before proceeding
	if botApi == nil {
		return nil, errors.New("telegram Bot API is nil")
	}

	return &tgBot{
		BaseBot: BaseBot{
			platform:  TELEGRAM,
			conf:      botconf,
			database:  database,
			dao:       dao,
			aiClients: aiClients,
			//openAIclient:  openaiClient,
			//mistralClient: mistralClient,
			embConfig: embconf,
		},
		botApi: botApi,
		//openAIclient: openai.NewClient(),
	}, nil
}

// SetWebhook sets the webhook for Telegram bot
func (b *tgBot) setWebhook(webhookURL string) error {
	webhookConfig, err := tgbotapi.NewWebhook(webhookURL)
	if err != nil {
		return fmt.Errorf("error creating webhook config: %w", err)
	}

	// Send the request to set the webhook
	_, err = b.botApi.Request(webhookConfig)
	if err != nil {
		return fmt.Errorf("error setting webhook: %w", err)
	}

	return nil
}

func (b *tgBot) Run() error {
	botApi, err := tgbotapi.NewBotAPI(b.conf.TelegramBotToken) // create new BotAPI instance using the token
	// utils.TgBot: global variable (defined in the utils package) that holds the reference to the bot instance.
	if err != nil {
		return err
	}

	b.botApi = botApi

	// Start the bot with webhook
	fmt.Println("Telegram bot is running with webhook!")

	// // Use go routine to continuously process received updates from the updates channel
	// go b.receiveUpdates(b.ctx, updates)
	return b.setWebhook(b.BaseBot.conf.TelegramWebhookURL)
}

// validateUser checks if the user exists in the database and creates a new record if not.
func (b *tgBot) ValidateUser(user *tgbotapi.User, message *tgbotapi.Message) (bool, error) {
	var dbUser models.User

	userIDStr := strconv.FormatInt(user.ID, 10)
	fmt.Printf("User ID: %s \n", userIDStr)

	// Check if the user exists in the database.
	err := b.database.GetDB().Where("user_id = ? AND deleted_at IS NULL", userIDStr).First(&dbUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// User does not exist; create a new record.
			dbUser = models.User{
				UserID:       userIDStr,
				UserName:     user.UserName,
				FirstName:    user.FirstName,
				LastName:     user.LastName,
				LanguageCode: user.LanguageCode,
			}

			// Save the new user to the database.
			if err := b.database.GetDB().Create(&dbUser).Error; err != nil {
				return false, fmt.Errorf("error creating user: %w", err)
			}

			// Send a welcome message to the new user.
			welcomeMessage := fmt.Sprintf("Welcome, %s!", user.UserName)
			if err := b.SendReply(message, welcomeMessage); err != nil {
				return false, fmt.Errorf("error sending welcome message: %w", err)
			}

			return false, nil // User was just created and welcomed.
		}
		return false, fmt.Errorf("error retrieving user: %w", err)
	}
	return true, nil // User already exists.
}

// validateAndGenerateToken checks if the user exists in the database and generates a token if not
// func (b *tgBot) validateAndGenerateToken(userIDStr string, user *tgbotapi.User) (*string, error) {
// 	// Check if the user exists in the database
// 	var dbUser models.User
// 	err := b.Service.GetDB().Where("user_id = ? AND deleted_at IS NULL", userIDStr).First(&dbUser).Error
// 	if err != nil {
// 		// If user does not exist, create a new user
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			dbUser = models.User{
// 				UserID:       userIDStr,
// 				UserName:     user.UserName,
// 				FirstName:    user.FirstName,
// 				LastName:     user.LastName,
// 				LanguageCode: user.LanguageCode,
// 			}

// 			// Create the new user record in the database
// 			if err := b.Service.GetDB().Create(&dbUser).Error; err != nil {
// 				return nil, fmt.Errorf("error creating user: %w", err)
// 			}

// 			// Generate a JWT token using the service's ValidateUser method
// 			token, err := b.Service.ValidateUser(userIDStr, service.ValidateUserReq{
// 				FirstName:    user.FirstName,
// 				LastName:     user.LastName,
// 				UserName:     user.UserName,
// 				LanguageCode: user.LanguageCode,
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

// Process user messages and respond accordingly
// func (b *tgBot) processUserMessage(message *tgbotapi.Message, firstName, text string) { //chatID int64
// 	chatID := message.Chat.ID

// 	fmt.Printf("Received message from %s: %s \n", firstName, text)
// 	fmt.Printf("Chat ID: %d \n", chatID)

// 	var response string
// 	//var err error

// 	if strings.HasPrefix(text, "/") {
// 		response = b.BaseBot.HandleCommand(text)
// 		/*response, err = handleCommand(chatID, text, b)
// 		if err != nil {
// 			fmt.Printf("An error occurred: %s \n", err.Error())
// 			response = "An error occurred while processing your command."
// 		}*/
// 	} else if b.conf.Screaming && len(text) > 0 {
// 		response = strings.ToUpper(text)
// 	} else {
// 		// Get all document embeddings
// 		documentEmbeddings, chunkText, err := b.BaseBot.dao.FetchEmbeddings()
// 		//documentEmbeddings, chunkText, err := b.Service.GetAllDocumentEmbeddings()
// 		if err != nil {
// 			fmt.Printf("Error retrieving document embeddings: %v", err)
// 			response = "Error retrieving document embeddings."
// 		} else if b.conf.UseOpenAI {
// 			//conf := config.GetConfig()
// 			// Perform similarity matching with the user's message when OpenAI is enabled
// 			topChunksText, err := document.RetrieveTopNChunks(text, documentEmbeddings, b.embConfig.NumTopChunks, chunkText, b.embConfig.ScoreThreshold) // Returns maximum N chunks with similarity threshold
// 			if err != nil {
// 				fmt.Printf("Error retrieving document chunks: %v", err)
// 				response = "Error retrieving related document information."
// 			} else if len(topChunksText) > 0 {
// 				// If there are similar chunks found, respond with those
// 				//response = fmt.Sprintf("Found related information:\n%s", strings.Join(topChunksText, "\n"))

// 				// If there are similar chunks found, provide them as context for GPT
// 				context := strings.Join(topChunksText, "\n")
// 				gptPrompt := fmt.Sprintf("Context:\n%s\nUser query: %s", context, text)

// 				// Call GPT with the context and user query
// 				response, err = b.BaseBot.GetOpenAIResponse(gptPrompt)
// 				if err != nil {
// 					response = fmt.Sprintf("OpenAI Error: %v", err)
// 				} /*else {
// 					response = fmt.Sprintf("Found related information based on context:\n%s", response)
// 				}*/
// 			} else {
// 				// If no relevant document found, fallback to OpenAI response
// 				response, err = b.BaseBot.GetOpenAIResponse(text)
// 				//response = fmt.Sprintf("Found related information:\n%s", strings.Join(topChunksText, "\n"))
// 				if err != nil {
// 					response = fmt.Sprintf("OpenAI Error: %v", err)
// 				}
// 			}
// 		} else {
// 			// Fall back to Dialogflow if OpenAI is not enabled
// 			//b.BaseBot.handleMessageDialogflow(TELEGRAM, message, text, b) //TODO
// 			return
// 		}
// 	}

// 	/*if strings.HasPrefix(text, "/") {
// 		response, err = handleCommand(chatID, text, b)
// 		if err != nil {
// 			fmt.Printf("An error occurred: %s \n", err.Error())
// 			response = "An error occurred while processing your command."
// 		}
// 	} else if screaming && len(text) > 0 {
// 		response = strings.ToUpper(text)
// 	} else if useOpenAI {
// 		// Call OpenAI to get the response
// 		response, err = GetOpenAIResponse(text)
// 		if err != nil {
// 			response = fmt.Sprintf("OpenAI Error: %v", err)
// 		}
// 	} else {
// 		handleMessageDialogflow(TELEGRAM, message, text, b)
// 		return
// 	}*/

// 	if response != "" {
// 		b.sendTelegramMessage(chatID, response)
// 	}
// }

// handleDialogflowResponse processes and sends the Dialogflow response to the appropriate platform
// func (b *tgBot) handleDialogflowResponse(response *dialogflowpb.DetectIntentResponse, identifier interface{}) error {

// 	// Send the response to respective platform
// 	// by iterating over the fulfillment messages returned by Dialogflow and processes any text messages.
// 	for _, msg := range response.QueryResult.FulfillmentMessages {
// 		if _, ok := identifier.(*tgbotapi.Message); ok {
// 			if text := msg.GetText(); text != nil {
// 				return b.SendReply(identifier, text.Text[0])
// 			}
// 		}
// 	}
// 	return fmt.Errorf("invalid Telegram message identifier")
// }

// Check identifier and send message via Telegram
func (b *tgBot) SendReply(identifier interface{}, response string) error {
	if message, ok := identifier.(*tgbotapi.Message); ok { // Assertion to check if identifier is of type tgbotapi.Message
		return b.sendTelegramMessage(message.Chat.ID, response)
	} else {
		return fmt.Errorf("invalid identifier for Telegram platform")
	}
}

// Send a message via Telegram (TG requires manual construction of an HTTP request)
func (b *tgBot) sendTelegramMessage(chatID int64, messageText string) error {
	// Use the Telegram API URL from the config
	url := b.conf.TelegramAPIURL + b.conf.TelegramBotToken + "/sendMessage"
	//conf := config.GetConfig()
	//url := conf.TelegramAPIURL + b.token + "/sendMessage"

	// Create the message payload
	message := map[string]interface{}{
		"chat_id": chatID,
		"text":    messageText,
	}

	// Marshal the message payload to JSON
	jsonMessage, _ := json.Marshal(message)

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonMessage))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Send the request using the HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending response: %w", err)
	}
	defer resp.Body.Close()

	// Log the response (can be removed if not needed)
	log.Printf("Response sent to chat ID %d", chatID)

	return nil
}

// GetDocFile retrieves relevant file information of the document
func (b *tgBot) GetDocFile(update tgbotapi.Update) (string, string, string, error) {

	fileID := update.Message.Document.FileID
	filename := update.Message.Document.FileName
	fileURL, err := b.botApi.GetFileDirectURL(fileID)
	if err != nil {
		//b.sendTelegramMessage(update.Message.Chat.ID, "Error getting file: "+err.Error())
		return "", "", "", err
	}
	return fileID, fileURL, filename, nil
}
