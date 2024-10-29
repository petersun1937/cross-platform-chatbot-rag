package bot

import (
	config "crossplatform_chatbot/configs"
	"crossplatform_chatbot/database"
	"crossplatform_chatbot/document"
	"crossplatform_chatbot/models"
	"crossplatform_chatbot/repository"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"cloud.google.com/go/dialogflow/apiv2/dialogflowpb"
	"github.com/line/line-bot-sdk-go/linebot"
	"gorm.io/gorm"
)

type LineBot interface {
	Run() error
	ParseRequest(req *http.Request) ([]*linebot.Event, error)
	HandleLineMessage(event *linebot.Event, message *linebot.TextMessage)
}

type lineBot struct {
	BaseBot
	//conf config.BotConfig
	//secret     string
	//token      string
	lineClient *linebot.Client
	//service    *service.Service
}

func NewLineBot(conf config.BotConfig, database database.Database, dao repository.DAO) (*lineBot, error) {
	lineClient, err := linebot.New(conf.LineChannelSecret, conf.LineChannelToken)
	if err != nil {
		return nil, err
	}

	return &lineBot{
		BaseBot: BaseBot{
			Platform: LINE,
			conf:     conf,
			database: database,
			dao:      dao,
		},
		lineClient: lineClient,
	}, nil
}

// func NewLineBot(conf *config.Config, service *service.Service) (*lineBot, error) {
// 	lineClient, err := linebot.New(conf.LineChannelSecret, conf.LineChannelToken)
// 	if err != nil {
// 		return nil, err
// 	}

// 	baseBot := &BaseBot{
// 		Platform: LINE,
// 		Service:  service,
// 	}

// 	return &lineBot{
// 		BaseBot: baseBot,
// 		conf:    conf.BotConfig,
// 		//secret:     conf.LineChannelSecret,
// 		//token:      conf.LineChannelToken,
// 		lineClient: lineClient,
// 		service:    service,
// 	}, nil
// 	/*return &lineBot{
// 		secret: conf.GetLineSecret(),
// 		token:  conf.GetLineToken(),
// 	}*/
// }

func (b *lineBot) Run() error {
	// // Initialize Linebot
	// //lineClient, err := linebot.New(b.secret, b.token)
	// lineClient, err := linebot.New(b.conf.LineChannelSecret, b.conf.LineChannelToken) // create new BotAPI instance using the channel token and secret
	// if err != nil {
	// 	return err
	// }

	// b.lineClient = lineClient

	// Start the bot with webhook
	fmt.Println("Line bot is running with webhook!")

	return nil
}

func (b *lineBot) HandleLineMessage(event *linebot.Event, message *linebot.TextMessage) {

	// Retrieve and validate user profile
	userProfile, err := b.getUserProfile(event.Source.UserID)
	if err != nil {
		fmt.Printf("Error fetching user profile: %v\n", err)
		return
	}

	// Ensure user exists in the database
	userExists, err := b.validateUser(userProfile, event.Source.UserID)
	if err != nil {
		fmt.Printf("Error ensuring user exists: %v\n", err)
		return
	}
	// userExists, err := b.validateAndGenerateToken(userProfile, event, event.Source.UserID)
	// if err != nil {
	// 	fmt.Printf("Error ensuring user exists: %v\n", err)
	// 	return
	// }

	// If user didn't exist, a welcome message was sent, so return
	if !userExists {
		return
	}

	// Process the user's message
	b.processUserMessage(event, message.Text)
	// response, err := b.processUserMessage(event, message.Text)
	// if err != nil {
	// 	fmt.Printf("Error processing user message: %v\n", err)
	// 	return
	// }

	// Send the response if it's not empty
	// if response != "" {
	// 	if err := b.sendLineMessage(event.ReplyToken, response); err != nil {
	// 		fmt.Printf("Error sending response message: %v\n", err)
	// 	}
	// }
}

// Get user profile from Line
func (b *lineBot) getUserProfile(userID string) (*linebot.UserProfileResponse, error) {
	userProfile, err := b.lineClient.GetProfile(userID).Do()
	if err != nil {
		return nil, err
	}
	return userProfile, nil
}

// validateUser checks if the user exists in the database and creates a new record if not.
func (b *lineBot) validateUser(userProfile *linebot.UserProfileResponse, userID string) (bool, error) {
	var dbUser models.User

	// Check if the user exists in the database.
	err := b.database.GetDB().Where("user_id = ? AND deleted_at IS NULL", userID).First(&dbUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// User does not exist; create a new record.
			dbUser = models.User{
				UserID:       userID,
				UserName:     userProfile.DisplayName,
				FirstName:    "", // LINE doesn't provide first and last names
				LastName:     "",
				LanguageCode: userProfile.Language,
			}

			// Save the new user to the database.
			if err := b.database.GetDB().Create(&dbUser).Error; err != nil {
				return false, fmt.Errorf("error creating user: %w", err)
			}

			// Send a welcome message to the new user.
			welcomeMessage := fmt.Sprintf("Welcome, %s!", userProfile.DisplayName)
			if err := b.sendLineMessage(userID, welcomeMessage); err != nil {
				return false, fmt.Errorf("error sending welcome message: %w", err)
			}

			return false, nil // User was just created and welcomed.
		}
		return false, fmt.Errorf("error retrieving user: %w", err)
	}
	return true, nil // User already exists.
}

// validateAndGenerateToken checks if the user exists and generates a token if not
// func (b *lineBot) validateAndGenerateToken(userProfile *linebot.UserProfileResponse, event *linebot.Event, userID string) (bool, error) {
// 	var dbUser models.User
// 	// err := b.service.GetDB().Where("user_id = ? AND deleted_at IS NULL", userID).First(&dbUser).Error
// 	err := b.database.GetDB().Where("user_id = ? AND deleted_at IS NULL", userID).First(&dbUser).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			dbUser = models.User{
// 				UserID:       userID,
// 				UserName:     userProfile.DisplayName,
// 				FirstName:    "", // LINE doesn't provide first and last names
// 				LastName:     "",
// 				LanguageCode: userProfile.Language,
// 			}

// 			// Create the new user record in the database
// 			if err := b.database.GetDB().Create(&dbUser).Error; err != nil {
// 				//if err := b.service.GetDB().Create(&dbUser).Error; err != nil {
// 				return false, fmt.Errorf("error creating user: %w", err)
// 			}

// 			// Generate a JWT token using the service's ValidateUser method
// 			token, err := b.service.ValidateUser(userID, service.ValidateUserReq{
// 				FirstName:    "", // LINE doesn't provide first and last names
// 				LastName:     "", // LINE doesn't provide first and last names
// 				UserName:     userProfile.DisplayName,
// 				LanguageCode: userProfile.Language,
// 			})
// 			if err != nil {
// 				return false, fmt.Errorf("error generating JWT: %w", err)
// 			}

// 			// Send welcome message with the token
// 			if err := b.sendLineMessage(event.ReplyToken, "Welcome! Your access token is: "+*token); err != nil {
// 				return false, fmt.Errorf("error sending token message: %w", err)
// 			}

// 			return false, nil // User was just created and welcomed
// 		}
// 		return false, fmt.Errorf("error retrieving user: %w", err)
// 	}
// 	return true, nil // User already existed
// }

// Process the user's message (commands or Dialogflow)
func (b *lineBot) processUserMessage(event *linebot.Event, text string) {
	userID := event.Source.UserID
	fmt.Printf("Received message from %s: %s \n", userID, text)

	var response string
	var err error

	// Check if the message is a command (starts with "/")
	if strings.HasPrefix(text, "/") {
		response = handleCommand(text)
		/*response, err = handleCommand(userID, text, b)
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
		//documentEmbeddings, chunkText, err := b.database.GetAllDocumentEmbeddings()
		if err != nil {
			fmt.Printf("Error retrieving document embeddings: %v", err)
			response = "Error retrieving document embeddings."
		} else if useOpenAI {
			// Perform similarity matching with the user's message
			topChunks, err := document.RetrieveTopNChunks(text, documentEmbeddings, 10, chunkText, 0.7) // Top 10 relevant chunks, threshold score 0.7
			if err != nil {
				fmt.Printf("Error retrieving document chunks: %v", err)
				response = "Error retrieving related document information."
			} else if len(topChunks) > 0 {
				// If similar chunks are found, provide them as context for GPT
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
			handleMessageDialogflow(LINE, event, text, b)
			return
		}
	}

	// Send the response if it's not empty
	if response != "" {
		err = b.sendResponse(event, response)
		if err != nil {
			fmt.Printf("Error sending response: %s\n", err.Error())
		}
	}
}

// handleDialogflowResponse processes and sends the Dialogflow response to the appropriate platform
func (b *lineBot) handleDialogflowResponse(response *dialogflowpb.DetectIntentResponse, identifier interface{}) error {

	// Send the response to respective platform
	// by iterating over the fulfillment messages returned by Dialogflow and processes any text messages.
	for _, msg := range response.QueryResult.FulfillmentMessages {
		if _, ok := identifier.(*linebot.Event); ok {
			if text := msg.GetText(); text != nil {
				return b.sendResponse(identifier, text.Text[0])
			}
		}
	}
	return fmt.Errorf("invalid LINE event identifier")
}

// Check identifier and send message via LINE
func (b *lineBot) sendResponse(identifier interface{}, response string) error {
	if event, ok := identifier.(*linebot.Event); ok { // Assertion to check if identifier is of type linebot.Event
		return b.sendLineMessage(event.ReplyToken, response)
	} else {
		return fmt.Errorf("invalid identifier for LINE platform")
	}
}

// Send a message via LINE
func (b *lineBot) sendLineMessage(replyToken string, messageText string) error {
	// Create the message payload
	replyMessage := linebot.NewTextMessage(messageText)

	// Send the message
	_, err := b.lineClient.ReplyMessage(replyToken, replyMessage).Do()
	if err != nil {
		return fmt.Errorf("error sending LINE message: %w", err)
	}
	return nil
}

func (b *lineBot) ParseRequest(req *http.Request) ([]*linebot.Event, error) {
	return b.lineClient.ParseRequest(req)
}

/*func (b *lineBot) sendMenu(identifier interface{}) error {
	if event, ok := identifier.(*linebot.Event); ok {
		return b.sendLineMenu(event.ReplyToken)
	} else {
		return fmt.Errorf("invalid identifier type for LINE platform")
	}
}

func (b *lineBot) sendLineMenu(replyToken string) error {
	// Define the LINE menu
	firstMenu := linebot.NewTextMessage("Here's the LINE menu:")
	actions := linebot.NewURIAction("Visit website", "https://example.com")
	template := linebot.NewButtonsTemplate("", "Menu", "Select an option:", actions)
	message := linebot.NewTemplateMessage("Menu", template)

	// Send the menu message to the user
	_, err := b.lineClient.ReplyMessage(replyToken, firstMenu, message).Do()
	if err != nil {
		return fmt.Errorf("error sending LINE menu: %w", err)
	}
	return nil
}*/
