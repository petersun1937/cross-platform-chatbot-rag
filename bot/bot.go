package bot

import (
	config "crossplatform_chatbot/configs"
	"crossplatform_chatbot/database"
	"crossplatform_chatbot/repository"

	"cloud.google.com/go/dialogflow/apiv2/dialogflowpb"
)

// action, public function
type Bot interface {
	Run() error

	//sendMenu(identifier interface{}) error
	sendResponse(identifier interface{}, message string) error
	handleDialogflowResponse(response *dialogflowpb.DetectIntentResponse, identifier interface{}) error
}

type BaseBot struct {
	Platform Platform
	// Service  *service.Service
	conf     config.BotConfig
	database database.Database
	dao      repository.DAO
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
