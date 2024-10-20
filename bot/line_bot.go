package bot

import (
	config "crossplatform_chatbot/configs"
	"crossplatform_chatbot/models"
	"crossplatform_chatbot/service"
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
	*BaseBot
	conf config.BotConfig
	//secret     string
	//token      string
	lineClient *linebot.Client
	service    *service.Service
}

func NewLineBot(conf *config.Config, service *service.Service) (*lineBot, error) {
	lineClient, err := linebot.New(conf.LineChannelSecret, conf.LineChannelToken)
	if err != nil {
		return nil, err
	}

	baseBot := &BaseBot{
		Platform: LINE,
		Service:  service,
	}

	return &lineBot{
		BaseBot: baseBot,
		conf:    conf.BotConfig,
		//secret:     conf.LineChannelSecret,
		//token:      conf.LineChannelToken,
		lineClient: lineClient,
		service:    service,
	}, nil
	/*return &lineBot{
		secret: conf.GetLineSecret(),
		token:  conf.GetLineToken(),
	}*/
}

func (b *lineBot) Run() error {
	// Initialize Linebot
	//lineClient, err := linebot.New(b.secret, b.token)
	lineClient, err := linebot.New(b.conf.LineChannelSecret, b.conf.LineChannelToken) // create new BotAPI instance using the channel token and secret
	if err != nil {
		return err
	}

	b.lineClient = lineClient

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
	userExists, err := b.validateAndGenerateToken(userProfile, event, event.Source.UserID)
	if err != nil {
		fmt.Printf("Error ensuring user exists: %v\n", err)
		return
	}

	// If user didn't exist, a welcome message was sent, so return
	if !userExists {
		return
	}

	// Process the user's message
	response, err := b.processUserMessage(event, message.Text)
	if err != nil {
		fmt.Printf("Error processing user message: %v\n", err)
		return
	}

	// Send the response if it's not empty
	if response != "" {
		if err := b.sendLineMessage(event.ReplyToken, response); err != nil {
			fmt.Printf("Error sending response message: %v\n", err)
		}
	}
}

// Get user profile from Line
func (b *lineBot) getUserProfile(userID string) (*linebot.UserProfileResponse, error) {
	userProfile, err := b.lineClient.GetProfile(userID).Do()
	if err != nil {
		return nil, err
	}
	return userProfile, nil
}

// validateAndGenerateToken checks if the user exists and generates a token if not
func (b *lineBot) validateAndGenerateToken(userProfile *linebot.UserProfileResponse, event *linebot.Event, userID string) (bool, error) {
	var dbUser models.User
	err := b.service.GetDB().Where("user_id = ? AND deleted_at IS NULL", userID).First(&dbUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			dbUser = models.User{
				UserID:       userID,
				UserName:     userProfile.DisplayName,
				FirstName:    "", // LINE doesn't provide first and last names
				LastName:     "",
				LanguageCode: userProfile.Language,
			}

			// Create the new user record in the database
			if err := b.service.GetDB().Create(&dbUser).Error; err != nil {
				return false, fmt.Errorf("error creating user: %w", err)
			}

			// Generate a JWT token using the service's ValidateUser method
			token, err := b.service.ValidateUser(userID, service.ValidateUserReq{
				FirstName:    "", // LINE doesn't provide first and last names
				LastName:     "", // LINE doesn't provide first and last names
				UserName:     userProfile.DisplayName,
				LanguageCode: userProfile.Language,
			})
			if err != nil {
				return false, fmt.Errorf("error generating JWT: %w", err)
			}

			// Send welcome message with the token
			if err := b.sendLineMessage(event.ReplyToken, "Welcome! Your access token is: "+*token); err != nil {
				return false, fmt.Errorf("error sending token message: %w", err)
			}

			return false, nil // User was just created and welcomed
		}
		return false, fmt.Errorf("error retrieving user: %w", err)
	}
	return true, nil // User already existed
}

// Process the user's message (commands or Dialogflow)
func (b *lineBot) processUserMessage(event *linebot.Event, text string) (string, error) {
	var response string
	var err error

	if strings.HasPrefix(text, "/") {
		response = handleCommand(text)
		/*response, err = handleCommand(event, text, b)
		if err != nil {
			return "An error occurred while processing your command.", err
		}*/
	} else if screaming && len(text) > 0 {
		response = strings.ToUpper(text)
	} else if useOpenAI {
		response, err = GetOpenAIResponse(text)
		if err != nil {
			response = "Error contacting OpenAI."
		}
	} else {
		handleMessageDialogflow(LINE, event, text, b)
		return "", nil
	}

	return response, nil
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
