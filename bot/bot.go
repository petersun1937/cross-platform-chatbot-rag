package bot

import (
	"crossplatform_chatbot/ai_clients"
	config "crossplatform_chatbot/configs"
	"crossplatform_chatbot/database"
	"crossplatform_chatbot/repository"
)

// action, public function
type Bot interface {
	Run() error

	//sendMenu(identifier interface{}) error
	SendReply(identifier interface{}, message string) error
	//handleDialogflowResponse(response *dialogflowpb.DetectIntentResponse, identifier interface{}) error
	Base() *BaseBot // Make BaseBot accessible
	Platform() Platform
}

type BaseBot struct {
	platform Platform
	// Service  *service.Service
	conf      *config.BotConfig
	database  database.Database
	dao       repository.DAO
	aiClients ai_clients.AIClients
	//openAIclient  *openai.Client
	//mistralClient *mistral.Client
	embConfig config.EmbeddingConfig
}

func (b *BaseBot) Base() *BaseBot {
	return b
}

func (b *BaseBot) Platform() Platform {
	return b.platform
}

// define platforms
type Platform int

const (
	LINE Platform = iota
	TELEGRAM
	FACEBOOK
	INSTAGRAM
	GENERAL
)
