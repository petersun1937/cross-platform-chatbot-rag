package bot

import (
	"crossplatform_chatbot/ai_clients"
	config "crossplatform_chatbot/configs"
	"crossplatform_chatbot/database"
	"crossplatform_chatbot/repository"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GeneralBot interface {
	Run() error
	StoreContext(sessionID string, c *gin.Context)
	//SetWebhook(webhookURL string) error
}

type generalBot struct {
	BaseBot
}

// creates a new GeneralBot instance
func NewGeneralBot(botconf *config.BotConfig, embconf config.EmbeddingConfig, aiClients ai_clients.AIClients, database database.Database, dao repository.DAO) (*generalBot, error) {

	return &generalBot{
		BaseBot: BaseBot{
			platform:  GENERAL,
			conf:      botconf,
			database:  database,
			dao:       dao,
			aiClients: aiClients,
			embConfig: embconf,
		},
	}, nil
}

func (b *generalBot) Run() error {
	// Implement logic for running the bot
	fmt.Println("General bot is running...")
	return nil
}

func (b *generalBot) SendReply(identifier interface{}, response string) error {
	// Perform type assertion to convert identifier to string
	if sessionID, ok := identifier.(string); ok {
		// Retrieve context using the sessionID
		c, err := getContext(sessionID)
		if err != nil {
			return fmt.Errorf("failed to retrieve context for sessionID: %s, error: %w", sessionID, err)
		}
		// Call sendFrontendMessage using the retrieved context
		return b.sendFrontendMessage(c, response)
	}
	return fmt.Errorf("invalid identifier type, expected string")
}

func (b *generalBot) sendFrontendMessage(c *gin.Context, message string) error {
	if c == nil {
		return fmt.Errorf("gin context is nil")
	}

	// Merge the response message into extraData
	responseData := gin.H{
		"response": message,
	}

	// Merge extraData into responseData
	/*for key, value := range extraData {
		responseData[key] = value
	}*/

	// Send the combined response
	c.JSON(http.StatusOK, responseData)
	return nil
}

var sessionContextMap = make(map[string]*gin.Context)

// StoreContext stores the context in sessionContextMap using the session ID
func (b *generalBot) StoreContext(sessionID string, c *gin.Context) {
	sessionContextMap[sessionID] = c
}

// Retrieve the context using sessionID when you need to send a response
func getContext(sessionID string) (*gin.Context, error) {
	if context, ok := sessionContextMap[sessionID]; ok {
		return context, nil
	}
	return nil, fmt.Errorf("no context found for session ID %s", sessionID)
}
