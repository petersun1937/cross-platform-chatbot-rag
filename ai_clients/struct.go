package ai_clients

import (
	"crossplatform_chatbot/ai_clients/huggingface"
	"crossplatform_chatbot/ai_clients/mistral"
	"crossplatform_chatbot/ai_clients/openai"
	"crossplatform_chatbot/ai_clients/togetherai"
)

// AIClient defines a common interface for all AI services
type AIClients struct {
	OpenAI      *openai.Client
	HuggingFace *huggingface.Client
	Mistral     *mistral.Client
	TogetherAI  *togetherai.Client
}
