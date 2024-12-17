package utils

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Global variables to hold the bot instances for TG and LINE
//var TgBot *tgbotapi.BotAPI
//var LineBot *linebot.Client

func SanitizeText(input string) string {
	validRunes := []rune{}
	for _, r := range input {
		if r == utf8.RuneError {
			continue // Skip invalid characters
		}
		validRunes = append(validRunes, r)
	}
	return string(validRunes)
}

// Helper function to remove duplicates from a list of tags
func RemoveDuplicates(tags []string) []string {
	uniqueTags := make(map[string]bool)
	var result []string

	for _, tag := range tags {
		if !uniqueTags[tag] { // Check if the tag is unique
			uniqueTags[tag] = true
			result = append(result, tag)
		}
	}

	return result
}

// Convert float64 slice to PostgreSQL float8[] string format
func Float64SliceToPostgresArray(embedding []float64) string {
	var result strings.Builder
	result.WriteString("{")
	for i, value := range embedding {
		if i > 0 {
			result.WriteString(",")
		}
		result.WriteString(fmt.Sprintf("%f", value))
	}
	result.WriteString("}")
	return result.String()
}

// Convert data type to store embeddings in Postgres
func PostgresArrayToFloat64Slice(embeddingStr string) ([]float64, error) {
	// Remove curly braces from the string
	embeddingStr = strings.Trim(embeddingStr, "{}")

	// Split the string by commas
	stringValues := strings.Split(embeddingStr, ",")

	// Convert the string values back to float64
	floatValues := make([]float64, len(stringValues))
	for i, strVal := range stringValues {
		val, err := strconv.ParseFloat(strings.TrimSpace(strVal), 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing embedding value: %v", err)
		}
		floatValues[i] = val
	}

	return floatValues, nil
}

func ParseEmbeddingString(embeddingStr string) ([]float64, error) {
	// Remove curly braces and split by commas
	trimmed := strings.Trim(embeddingStr, "{}")
	parts := strings.Split(trimmed, ",")

	embedding := make([]float64, len(parts))
	for i, part := range parts {
		value, err := strconv.ParseFloat(strings.TrimSpace(part), 64)
		if err != nil {
			return nil, err
		}
		embedding[i] = value
	}

	return embedding, nil
}

func AverageEmbeddings(embeddings [][]float64) ([]float64, error) {
	if len(embeddings) == 0 {
		return nil, fmt.Errorf("no embeddings provided")
	}

	length := len(embeddings[0])
	combinedEmbedding := make([]float64, length)

	for _, embedding := range embeddings {
		for i := range embedding {
			combinedEmbedding[i] += embedding[i]
		}
	}

	// Divide each element by the number of embeddings to get the average
	for i := range combinedEmbedding {
		combinedEmbedding[i] /= float64(len(embeddings))
	}

	return combinedEmbedding, nil
}

// Helper function to decode GOOGLE_APPLICATION_CREDENTIALS_JSON and create a temp file
func DecodeGoogleCredentials(encodedCreds string) (string, error) {

	if encodedCreds == "" {
		return "", fmt.Errorf("environment variable GOOGLE_APPLICATION_CREDENTIALS_JSON is not set")
	}

	// Decode the base64 string
	credsBytes, err := base64.StdEncoding.DecodeString(encodedCreds)
	if err != nil {
		return "", fmt.Errorf("failed to decode GOOGLE_APPLICATION_CREDENTIALS_JSON: %v", err)
	}

	// Write to a temporary file using os.CreateTemp
	tmpFile, err := os.CreateTemp("", "google-creds-*.json")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file for Google credentials: %v", err)
	}
	defer tmpFile.Close()

	// Write the decoded credentials to the temporary file
	if _, err := tmpFile.Write(credsBytes); err != nil {
		return "", fmt.Errorf("failed to write credentials to temp file: %v", err)
	}

	return tmpFile.Name(), nil
}
