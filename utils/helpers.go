package utils

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	dialogflowpb "cloud.google.com/go/dialogflow/apiv2/dialogflowpb"
)

// Global variables to hold the bot instances for TG and LINE
//var TgBot *tgbotapi.BotAPI
//var LineBot *linebot.Client

// Send a text query to Dialogflow and returns the response
func DetectIntentText(projectID, sessionID, text, languageCode string) (*dialogflowpb.DetectIntentResponse, error) {
	// Create a background context for the API call
	ctx := context.Background()

	// Create a new Dialogflow Sessions client
	client, err := dialogflow.NewSessionsClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close() // Ensure the client is closed when done

	// Construct the session path for the Dialogflow API
	sessionPath := fmt.Sprintf("projects/%s/agent/sessions/%s", projectID, sessionID)

	// Create the DetectIntentRequest with the session path and query input
	req := &dialogflowpb.DetectIntentRequest{
		Session: sessionPath,
		QueryInput: &dialogflowpb.QueryInput{
			Input: &dialogflowpb.QueryInput_Text{
				Text: &dialogflowpb.TextInput{
					Text:         text,
					LanguageCode: languageCode,
				},
			},
		},
	}

	// Send the request and return the response or error
	return client.DetectIntent(ctx, req)
}

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
