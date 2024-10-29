package service

import (
	"crossplatform_chatbot/bot"
	config "crossplatform_chatbot/configs"
	"crossplatform_chatbot/database"
	"crossplatform_chatbot/repository"
	"fmt"
	"log"
)

type Service struct {
	bots     map[string]bot.Bot
	database database.Database // TODO
	//repository *repository.Repository
}

func NewService(botConfig config.BotConfig, embConfig config.EmbeddingConfig, db database.Database) *Service {

	dao := repository.NewDAO(db)

	return &Service{
		bots:     createBots(botConfig, embConfig, db, dao),
		database: db,
		//repository: repository.NewRepository(database),
	}
}

func (s *Service) RunBots() error {

	for _, bot := range s.bots {
		if err := bot.Run(); err != nil {
			// log.Fatal("running bot failed:", err)
			fmt.Printf("running bot failed: %s", err.Error())
			return err
		}
	}

	return nil
}

func createBots(botConfig config.BotConfig, embConfig config.EmbeddingConfig, database database.Database, dao repository.DAO) map[string]bot.Bot {
	// Initialize bots
	lineBot, err := bot.NewLineBot(botConfig, database, dao)
	if err != nil {
		//log.Fatal("Failed to initialize LINE bot:", err)
		fmt.Printf("Failed to initialize LINE bot: %s", err.Error())
	}

	tgBot, err := bot.NewTGBot(botConfig, embConfig, database, dao)
	if err != nil {
		//log.Fatal("Failed to initialize Telegram bot:", err)
		fmt.Printf("Failed to initialize Telegram bot: %s", err.Error())
	}

	fbBot, err := bot.NewFBBot(botConfig, database, dao)
	if err != nil {
		log.Fatalf("Failed to create Facebook bot: %v", err)
	}

	igBot, err := bot.NewIGBot(botConfig, database, dao)
	if err != nil {
		log.Fatalf("Failed to create Instagram bot: %v", err)
	}

	generalBot, err := bot.NewGeneralBot(botConfig, embConfig, database, dao)
	if err != nil {
		log.Fatalf("Failed to initialize General bot: %v", err)
	}

	return map[string]bot.Bot{
		"line":      lineBot,
		"telegram":  tgBot,
		"facebook":  fbBot,
		"instagram": igBot,
		"general":   generalBot,
	}
}

func (s *Service) GetBot(tag string) bot.Bot {
	return s.bots[tag]
}

type UserProfile struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ID        string `json:"id"` // Facebook User ID
}

func (s *Service) Init() error {
	// running bots
	for _, bot := range s.bots {
		if err := bot.Run(); err != nil {
			// log.Fatal("running bot failed:", err)
			fmt.Printf("running bot failed: %s", err.Error())
			return err
		}
	}
	return nil
}
