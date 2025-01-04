package service

import (
	"context"
	"fmt"
	"log"
	"strings"

	config "crossplatform_chatbot/configs"

	"github.com/redis/go-redis/v9"
)

func initRedis(redisConfig config.RedisConfig) *redis.Client {
	// Load Redis configuration
	Endpoint := redisConfig.RedisEndpoint
	Password := redisConfig.RedisPassword

	// Create a Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     Endpoint, // Redis server address
		Password: Password, // Redis password
		DB:       0,        // Use default DB
	})

	// Test the connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	fmt.Println("Connected to Redis successfully!")
	return client
}

func (s *Service) saveConversation(chatID, userMessage, botResponse string) error {
	ctx := context.Background()
	key := "conversation:" + chatID // Use chat/session ID as the key
	entry := fmt.Sprintf("User: %s\nBot: %s", userMessage, botResponse)
	return s.redisClient.RPush(ctx, key, entry).Err()
}

func (s *Service) getConversationHistory(chatID string, limit int64) (string, error) {
	ctx := context.Background()
	key := "conversation:" + chatID

	// Fetch the last `limit` entries
	history, err := s.redisClient.LRange(ctx, key, -limit, -1).Result()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve conversation history from Redis: %v", err)
	}

	// Combine the entries into a single string
	return strings.Join(history, "\n"), nil
}
