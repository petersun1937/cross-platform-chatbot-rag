package bot

import (
	config "crossplatform_chatbot/configs"
	"crossplatform_chatbot/database"
	openai "crossplatform_chatbot/openai"
	"crossplatform_chatbot/repository"
)

// action, public function
type Bot interface {
	Run() error
	SendReply(identifier interface{}, message string) error
	Base() *BaseBot // Make BaseBot accessible
	Platform() Platform
}

type BaseBot struct {
	platform Platform
	// Service  *service.Service
	conf         *config.BotConfig
	database     database.Database
	dao          repository.DAO
	openAIclient *openai.Client
	embConfig    config.EmbeddingConfig
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
