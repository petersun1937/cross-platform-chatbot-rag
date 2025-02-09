package bot

import (
	"crossplatform_chatbot/ai_clients"
	config "crossplatform_chatbot/configs"
	"crossplatform_chatbot/database"
	"crossplatform_chatbot/models"
	"crossplatform_chatbot/repository"
	"errors"
	"fmt"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
	"gorm.io/gorm"
)

type LineBot interface {
	Run() error
	ParseRequest(req *http.Request) ([]*linebot.Event, error)
	//HandleLineMessage(event *linebot.Event, message *linebot.TextMessage)
	GetUserProfile(userID string) (*linebot.UserProfileResponse, error)
	ValidateUser(userProfile *linebot.UserProfileResponse, userID string) (bool, error)
	//sendResponse(identifier interface{}, response string) error
}

type lineBot struct {
	BaseBot
	lineClient *linebot.Client
}

func NewLineBot(conf *config.BotConfig, embconf config.EmbeddingConfig, aiClients ai_clients.AIClients, database database.Database, dao repository.DAO) (*lineBot, error) {
	lineClient, err := linebot.New(conf.LineChannelSecret, conf.LineChannelToken)
	if err != nil {
		return nil, err
	}

	return &lineBot{
		BaseBot: BaseBot{
			platform:  LINE,
			conf:      conf,
			database:  database,
			dao:       dao,
			aiClients: aiClients,
			embConfig: embconf,
		},
		lineClient: lineClient,
	}, nil
}

func (b *lineBot) Run() error {

	// Start the bot with webhook
	fmt.Println("Line bot is running with webhook!")

	return nil
}

// Get user profile from Line
func (b *lineBot) GetUserProfile(userID string) (*linebot.UserProfileResponse, error) {
	userProfile, err := b.lineClient.GetProfile(userID).Do()
	if err != nil {
		return nil, err
	}
	return userProfile, nil
}

// validateUser checks if the user exists in the database and creates a new record if not.
func (b *lineBot) ValidateUser(userProfile *linebot.UserProfileResponse, userID string) (bool, error) {
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

// Check identifier and send message via LINE
func (b *lineBot) SendReply(identifier interface{}, response string) error {
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
