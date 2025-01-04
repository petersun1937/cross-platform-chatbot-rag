package service

import (
	"crossplatform_chatbot/bot"
	config "crossplatform_chatbot/configs"
	"crossplatform_chatbot/database"
	"crossplatform_chatbot/openai"
	"crossplatform_chatbot/repository"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type Service struct {
	bots         map[string]bot.Bot
	database     database.Database // TODO
	repository   repository.DAO
	openaiClient *openai.Client
	redisClient  *redis.Client
	botConfig    *config.BotConfig
	embConfig    config.EmbeddingConfig
	//TagEmbeddings map[string][]float64
}

func NewService(botConfig *config.BotConfig, embConfig *config.EmbeddingConfig, redisConfig config.RedisConfig, db database.Database) *Service {
	// Initialize the DAO, OpenAI and Redis clients
	dao := repository.NewDAO(db)
	openaiClient := openai.NewClient()
	redisClient := initRedis(redisConfig)

	// Create a temporary Service instance to access methods like getOrInitializeTagEmbeddings
	svc := &Service{
		database:     db,
		repository:   dao,
		openaiClient: openaiClient,
		redisClient:  redisClient,
		botConfig:    botConfig,
		embConfig:    *embConfig, // pointer no pointer?
	}

	// Now create bots (with the updated embConfig if using emb based tagging)
	svc.bots = createBots(botConfig, *embConfig, db, dao)

	return svc
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

func createBots(botConfig *config.BotConfig, embConfig config.EmbeddingConfig, database database.Database, dao repository.DAO) map[string]bot.Bot {
	// Initialize bots
	lineBot, err := bot.NewLineBot(botConfig, database, embConfig, dao)
	if err != nil {
		//log.Fatal("Failed to initialize LINE bot:", err)
		log.Fatalf("Failed to initialize LINE bot: %v", err)
	}

	tgBot, err := bot.NewTGBot(botConfig, embConfig, database, dao)
	if err != nil {
		//log.Fatal("Failed to initialize Telegram bot:", err)
		log.Fatalf("Failed to initialize Telegram bot: %v", err)
	}

	fbBot, err := bot.NewFBBot(botConfig, database, embConfig, dao)
	if err != nil {
		log.Fatalf("Failed to initialize Facebook bot: %v", err)
	}

	igBot, err := bot.NewIGBot(botConfig, database, embConfig, dao)
	if err != nil {
		log.Fatalf("Failed to initialize Instagram bot: %v", err)
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
