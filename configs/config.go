package config

import (
	"crossplatform_chatbot/utils"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerConfig
	BotConfig
	EmbeddingConfig
	OpenAIConfig
	RedisConfig
	HuggingFaceConfig
	MistralConfig
	TogetherAIConfig
	// DBString            string
	// AppPort             string
	// TelegramBotToken    string
	// LineChannelSecret   string
	// LineChannelToken    string
	// ServerConfig        ServerConfig
	// TelegramAPIURL      string
	// TelegramWebhookURL  string
	// DialogflowProjectID string
	// FacebookAPIURL      string
	// FacebookPageToken   string
	// FacebookVerifyToken string
	//DBUser string
	//DBPwd  string
}

type ServerConfig struct {
	Host     string
	Port     int // generally int
	Timeout  time.Duration
	MaxConn  int
	DBString string
	AppPort  string
}

type BotConfig struct {
	TelegramBotToken          string
	LineChannelSecret         string
	LineChannelToken          string
	TelegramAPIURL            string
	TelegramWebhookURL        string
	DialogflowProjectID       string
	GoogleCredentialsFilePath string
	FacebookAPIURL            string
	FacebookPageToken         string
	FacebookVerifyToken       string
	InstagramVerifyToken      string
	InstagramPageToken        string
	Screaming                 bool
	UseOpenAI                 bool
	UseMistral                bool
	UseMETA                   bool
	UseDialogflow             bool
}

type OpenAIConfig struct {
	OpenaiAPIKey   string
	OpenaiEmbModel string
	OpenaiMsgModel string
	OpenaiTagModel string
	MaxTokens      int
	MaxTagTokens   int
}

type RedisConfig struct {
	RedisEndpoint string
	RedisPassword string
}

type EmbeddingConfig struct {
	//EmbeddingBatchSize int
	ChunkSize      int
	MinChunkSize   int
	OverlapSize    int
	ScoreThreshold float64
	NumTopChunks   int
	TagEmbeddings  map[string][]float64
}

type HuggingFaceConfig struct {
	HuggingFaceAPIKey string
	HuggingFaceModel  string
}

type MistralConfig struct {
	MistralAPIKey string
	MistralModel  string
}

type TogetherAIConfig struct {
	TogetherAIAPIKey string
	TogetherAIModel  string
}

// Singleton instance of Config
var instance *Config
var once sync.Once

/*func init() {
	err := loadConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}
}*/

/*func NewConfig() *Config {
	return &Config{}
}*/

// Returns the singleton instance of Config
func GetConfig() *Config {
	// Ensure the config is initialized only once
	once.Do(func() {
		err := loadConfig()
		if err != nil {
			panic(fmt.Sprintf("Failed to load config: %v", err))
		}
	})
	return instance
}

// Load the configuration into the singleton instance
func loadConfig() error {
	// Load the .env file only if the DATABASE_URL is not already set
	if !isEnvSet("DATABASE_URL") {
		err := godotenv.Load("configs/.env")
		if err != nil {
			fmt.Printf("Warning: .env file not found: %v. Continuing without it...\n", err)
		}
	}

	// Decode Google credentials if set
	encodedCreds := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_JSON")
	googleCredsPath, err := utils.DecodeGoogleCredentials(encodedCreds)
	if err != nil {
		fmt.Println("Warning: Google credentials not set or invalid. Continuing without it...")
	}

	// Initialize the config struct with environment variables
	instance = &Config{
		ServerConfig: ServerConfig{
			Host: os.Getenv("SERVER_HOST"),
			//Port:     getEnvInt("PORT", 8080),
			Port:     getEnvInt("PORT", -1), // Default to -1 to detect if PORT is missing
			Timeout:  getEnvDuration("SERVER_TIMEOUT", 30*time.Second),
			MaxConn:  getEnvInt("SERVER_MAX_CONN", 100),
			DBString: os.Getenv("DATABASE_URL"),
		},
		BotConfig: BotConfig{
			TelegramBotToken:          os.Getenv("TELEGRAM_BOT_TOKEN"),
			LineChannelSecret:         os.Getenv("LINE_CHANNEL_SECRET"),
			LineChannelToken:          os.Getenv("LINE_CHANNEL_TOKEN"),
			TelegramAPIURL:            os.Getenv("TELEGRAM_API_URL"),
			TelegramWebhookURL:        os.Getenv("TELEGRAM_WEBHOOK_URL"),
			DialogflowProjectID:       os.Getenv("DIALOGFLOW_PROJECTID"),
			GoogleCredentialsFilePath: googleCredsPath,
			FacebookAPIURL:            os.Getenv("FACEBOOK_API_URL"),
			FacebookPageToken:         os.Getenv("FACEBOOK_PAGE_TOKEN"),
			FacebookVerifyToken:       os.Getenv("FACEBOOK_VERIFY_TOKEN"),
			InstagramVerifyToken:      os.Getenv("IG_VERIFY_TOKEN"),
			InstagramPageToken:        os.Getenv("IG_PAGE_TOKEN"),
			Screaming:                 false,
			UseOpenAI:                 true,
			UseMistral:                false,
			UseMETA:                   false,
			UseDialogflow:             true,
		},
		OpenAIConfig: OpenAIConfig{
			OpenaiAPIKey:   os.Getenv("OPENAI_API_KEY"),
			OpenaiEmbModel: os.Getenv("OPENAI_EMBED_MODEL"),
			OpenaiMsgModel: os.Getenv("OPENAI_MSG_MODEL"),
			OpenaiTagModel: os.Getenv("OPENAI_TAG_MODEL"),
			MaxTokens:      getEnvInt("OPENAI_MAX_TOKEN_SIZE", 250),
			MaxTagTokens:   getEnvInt("OPENAI_MAX_TAG_TOKEN_SIZE", 4097),
		},
		EmbeddingConfig: EmbeddingConfig{
			//EmbeddingBatchSize: getEnvInt("DOC_EMBEDDING_BATCH_SIZE", 10),
			ChunkSize:      getEnvInt("DOC_CHUNK_SIZE", 500),
			OverlapSize:    getEnvInt("DOC_OVERLAP_CHUNK_SIZE", 100),
			MinChunkSize:   getEnvInt("DOC_MIN_CHUNK_SIZE", 50),
			ScoreThreshold: getEnvFloat("DOC_SCORE_THRESHOLD", 0.65),
			NumTopChunks:   getEnvInt("DOC_NUM_TOP_CHUNKS", 10),
			TagEmbeddings:  make(map[string][]float64),
		},
		RedisConfig: RedisConfig{
			RedisEndpoint: os.Getenv("REDIS_ENDPOINT"),
			RedisPassword: os.Getenv("REDIS_PASSWORD"),
		},
		HuggingFaceConfig: HuggingFaceConfig{
			HuggingFaceAPIKey: os.Getenv("HUGGINGFACE_API_KEY"),
			HuggingFaceModel:  os.Getenv("HUGGINGFACE_MODEL"),
		},
		MistralConfig: MistralConfig{
			MistralAPIKey: os.Getenv("MISTRAL_API_KEY"),
			MistralModel:  os.Getenv("MISTRAL_MODEL"),
		},
		TogetherAIConfig: TogetherAIConfig{
			TogetherAIAPIKey: os.Getenv("TOGETHERAI_API_KEY"),
			TogetherAIModel:  os.Getenv("TOGETHERAI_MODEL"),
		},
	}

	// Validate required config values in a more concise way
	missingVars := []string{}
	if instance.ServerConfig.DBString == "" {
		missingVars = append(missingVars, "DATABASE_URL")
	}
	if instance.BotConfig.TelegramBotToken == "" {
		missingVars = append(missingVars, "TELEGRAM_BOT_TOKEN")
	}
	if instance.BotConfig.TelegramAPIURL == "" {
		missingVars = append(missingVars, "TELEGRAM_API_URL")
	}

	// Return an error if any required environment variables are missing
	if len(missingVars) > 0 {
		return fmt.Errorf("required environment variables missing: %v", missingVars)
	}

	return nil
}

func (c *Config) Init() error {
	// Load environment variables
	return godotenv.Load("configs/.env")
}

// For resetting the singleton instance
func ResetConfig() {
	instance = nil     // Reset the instance for testing purposes
	once = sync.Once{} // Reset the sync.Once to allow re-initialization
}

func isEnvSet(key string) bool {
	_, exists := os.LookupEnv(key)
	return exists
}

// Utility function to get environment variable as an integer
func getEnvInt(name string, defaultVal int) int {
	value, exists := os.LookupEnv(name)
	if !exists {
		if defaultVal == -1 {
			fmt.Printf("Environment variable %s is not set and no default value provided. Exiting...\n", name)
			os.Exit(1) // Exit if a critical environment variable is missing
		}
		return defaultVal
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		fmt.Printf("Environment variable %s has an invalid value: %s. Expected an integer. Exiting...\n", name, value)
		os.Exit(1) // Exit if the value is not a valid integer
	}

	return intValue
}

/*func getEnvInt(name string, defaultVal int) int {
	if value, exists := os.LookupEnv(name); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultVal
}*/

func getEnvFloat(name string, defaultVal float64) float64 {
	if value, exists := os.LookupEnv(name); exists {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultVal
}

// Utility function to get environment variable as a duration
func getEnvDuration(name string, defaultVal time.Duration) time.Duration {
	if value, exists := os.LookupEnv(name); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultVal
}

func (c *ServerConfig) GetHost() string {
	return c.Host
}

func (c *BotConfig) GetTelegramBotToken() string {
	return c.TelegramBotToken
}
