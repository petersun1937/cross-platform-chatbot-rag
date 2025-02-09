package service

import (
	"crossplatform_chatbot/ai_clients"
	"crossplatform_chatbot/ai_clients/huggingface"
	"crossplatform_chatbot/ai_clients/mistral"
	"crossplatform_chatbot/ai_clients/openai"
	"crossplatform_chatbot/ai_clients/togetherai"
	"crossplatform_chatbot/bot"
	config "crossplatform_chatbot/configs"
	"crossplatform_chatbot/database"
	"crossplatform_chatbot/repository"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type Service struct {
	bots        map[string]bot.Bot
	database    database.Database
	repository  repository.DAO
	redisClient *redis.Client
	botConfig   *config.BotConfig
	embConfig   config.EmbeddingConfig
	aiClients   ai_clients.AIClients
}

func NewService(botConfig *config.BotConfig, embConfig *config.EmbeddingConfig, redisConfig config.RedisConfig, db database.Database) *Service {
	// Initialize the DAO, OpenAI and Redis clients
	dao := repository.NewDAO(db)
	redisClient := initRedis(redisConfig)

	// Initialize all AI clients and store them in the unified struct
	aiClients := ai_clients.AIClients{
		OpenAI:      openai.NewClient(),
		HuggingFace: huggingface.NewClient(),
		Mistral:     mistral.NewClient(),
		TogetherAI:  togetherai.NewClient(),
	}

	// Create a temporary Service instance to access methods like getOrInitializeTagEmbeddings
	svc := &Service{
		database:    db,
		repository:  dao,
		aiClients:   aiClients,
		redisClient: redisClient,
		botConfig:   botConfig,
		embConfig:   *embConfig,
	}

	// Now create bots (with the updated embConfig if using emb based tagging)
	svc.bots = createBots(botConfig, *embConfig, aiClients, db, dao)

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

func createBots(botConfig *config.BotConfig, embConfig config.EmbeddingConfig, aiClients ai_clients.AIClients, database database.Database, dao repository.DAO) map[string]bot.Bot {
	// Initialize bots
	lineBot, err := bot.NewLineBot(botConfig, embConfig, aiClients, database, dao)
	if err != nil {
		//log.Fatal("Failed to initialize LINE bot:", err)
		log.Fatalf("Failed to initialize LINE bot: %v", err)
	}

	tgBot, err := bot.NewTGBot(botConfig, embConfig, aiClients, database, dao)
	if err != nil {
		//log.Fatal("Failed to initialize Telegram bot:", err)
		log.Fatalf("Failed to initialize Telegram bot: %v", err)
	}

	fbBot, err := bot.NewFBBot(botConfig, embConfig, aiClients, database, dao)
	if err != nil {
		log.Fatalf("Failed to initialize Facebook bot: %v", err)
	}

	igBot, err := bot.NewIGBot(botConfig, embConfig, aiClients, database, dao)
	if err != nil {
		log.Fatalf("Failed to initialize Instagram bot: %v", err)
	}

	generalBot, err := bot.NewGeneralBot(botConfig, embConfig, aiClients, database, dao)
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
