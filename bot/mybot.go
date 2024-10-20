package bot

import (
	"context"
	config "crossplatform_chatbot/configs"
	"crossplatform_chatbot/document"
	openai "crossplatform_chatbot/openai"
	"crossplatform_chatbot/service"
	"fmt"
	"net/http"
	"strings"

	"cloud.google.com/go/dialogflow/apiv2/dialogflowpb"
	"github.com/gin-gonic/gin"
)

type GeneralBot interface {
	Run() error
	HandleGeneralMessage(context *gin.Context)
	StoreDocumentChunks(docID string, text string, chunkSize int, minchunkSize int) error
	ProcessDocument(sessionID string, filePath string) error
	//SetWebhook(webhookURL string) error
}

type generalBot struct {
	// Add any common fields if necessary, like configuration
	*BaseBot
	ctx          context.Context
	conf         config.BotConfig
	embConfig    config.EmbeddingConfig
	openAIclient *openai.Client

	//config map[string]string
}

func NewGeneralBot(conf *config.Config, service *service.Service) (*generalBot, error) {
	baseBot := &BaseBot{
		Platform: GENERAL,
		Service:  service,
	}

	return &generalBot{
		BaseBot:      baseBot,
		conf:         conf.BotConfig,
		embConfig:    conf.EmbeddingConfig,
		ctx:          context.Background(),
		openAIclient: openai.NewClient(),
	}, nil
}

func (b *generalBot) Run() error {
	// Implement logic for running the bot
	fmt.Println("General bot is running...")
	return nil
}

func (b *generalBot) HandleGeneralMessage(c *gin.Context) { // TODO: some to handler
	var req struct {
		SessionID string `json:"sessionID"`
		Message   string `json:"message"`
	}

	// Parse the request body
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("Invalid request: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Store the context (to use later for sending the response)
	storeContext(req.SessionID, c)

	/*user := message.From
	if user == nil {
		return
	}

	userIDStr := strconv.FormatInt(user.ID, 10)
	fmt.Printf("User ID: %s \n", userIDStr)

	token, err := b.validateAndGenerateToken(userIDStr, user)
	if err != nil {
		fmt.Printf("Error validating user: %s", err.Error())
		return
	}

	if token != nil {
		b.sendTelegramMessage(message.Chat.ID, "Welcome! Your access token is: "+*token)
	} else {
		b.processUserMessage(message, user.FirstName, message.Text)
	}*/

	// Process and send the message
	b.ProcessUserMessage(req.SessionID, req.Message)

	// Send the response back to the frontend using sendResponse
	/*err = b.sendFrontendMessage(c, response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while sending the response"})
		return
	}*/

	/*if err := b.sendResponse(req.SessionID, response); err != nil {
		fmt.Printf("An error occurred while sending the response: %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while sending the response"})
		return
	}*/
}

// ProcessUserMessage processes incoming messages
func (b *generalBot) ProcessUserMessage(sessionID string, message string) {
	var response string
	//var err error

	fmt.Printf("Received message %s \n", message)
	fmt.Printf("Chat ID: %s \n", sessionID)

	if strings.HasPrefix(message, "/") {

		response = handleCommand(message)
		/*response, err = handleCommand(sessionID, message, b)
		if err != nil {
			fmt.Printf("An error occurred: %s \n", err.Error())
			response = "An error occurred while processing your command."
		}*/
	} else if screaming && len(message) > 0 {
		response = strings.ToUpper(message)
	} else {
		// Get all document embeddings
		documentEmbeddings, chunkText, err := b.Service.GetAllDocumentEmbeddings()
		if err != nil {
			fmt.Printf("Error retrieving document embeddings: %v", err)
			response = "Error retrieving document embeddings."
		} else if useOpenAI {
			// Perform similarity matching with the user's message
			topChunks, err := document.RetrieveTopNChunks(message, documentEmbeddings, b.embConfig.NumTopChunks, chunkText, b.embConfig.ScoreThreshold) // Retrieve top 3 relevant chunks
			if err != nil {
				fmt.Printf("Error retrieving document chunks: %v", err)
				response = "Error retrieving related document information."
			} else if len(topChunks) > 0 {
				// If there are similar chunks found, provide them as context for GPT
				context := strings.Join(topChunks, "\n")
				gptPrompt := fmt.Sprintf("Context:\n%s\nUser query: %s", context, message)

				// Call GPT with the context and user query
				response, err = GetOpenAIResponse(gptPrompt)
				if err != nil {
					response = fmt.Sprintf("OpenAI Error: %v", err)
				}
			} else {
				// If no relevant document found, fallback to OpenAI response
				response, err = GetOpenAIResponse(message)
				if err != nil {
					response = fmt.Sprintf("OpenAI Error: %v", err)
				}
			}

		} else {
			//response = fmt.Sprintf("You said: %s", message)
			handleMessageDialogflow(GENERAL, sessionID, message, b)
		}
	}

	if response != "" {
		fmt.Printf("Sent message %s \n", response)
		err := b.sendResponse(sessionID, response)
		if err != nil {
			//c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while sending the response"})
			fmt.Printf("An error occurred while sending the response: %s\n", err.Error())
		}
	}

}

func (b *generalBot) sendResponse(identifier interface{}, response string) error {
	// Perform type assertion to convert identifier to string
	if sessionID, ok := identifier.(string); ok {
		// Retrieve context using the sessionID
		c, err := getContext(sessionID)
		if err != nil {
			return fmt.Errorf("failed to retrieve context for sessionID: %s, error: %w", sessionID, err)
		}
		// Call sendFrontendMessage using the retrieved context
		return b.sendFrontendMessage(c, response)
	}
	return fmt.Errorf("invalid identifier type, expected string")
}

var sessionContextMap = make(map[string]*gin.Context)

// Store the context when the session starts
func storeContext(sessionID string, c *gin.Context) {
	sessionContextMap[sessionID] = c
}

// Retrieve the context using sessionID when you need to send a response
func getContext(sessionID string) (*gin.Context, error) {
	if context, ok := sessionContextMap[sessionID]; ok {
		return context, nil
	}
	return nil, fmt.Errorf("no context found for session ID %s", sessionID)
}

func (b *generalBot) sendFrontendMessage(c *gin.Context, message string) error {
	c.JSON(http.StatusOK, gin.H{
		"response": message,
	})
	return nil
}

func (b *generalBot) handleDialogflowResponse(response *dialogflowpb.DetectIntentResponse, identifier interface{}) error {
	// Send the response to the respective platform or frontend
	for _, msg := range response.QueryResult.FulfillmentMessages {
		if text := msg.GetText(); text != nil {
			return b.sendResponse(identifier, text.Text[0])
		}
	}
	return fmt.Errorf("invalid identifier for frontend or platform")
}

/*func (b *generalBot) sendMenu(identifier interface{}) error {
	if sessionID, ok := identifier.(string); ok {
		// Logic to send menu to the frontend user TODO
		// For example, return a message to the frontend via the API response
		fmt.Printf("Sending menu to frontend user with session ID: %s\n", sessionID)
		return nil
	}
	return fmt.Errorf("invalid identifier type for frontend platform")
}*/

/*func (b *generalBot) StoreDocumentChunks(docID string, text string, chunkSize int, overlap int) error {
	// Chunk the document using the semantic chunking logic
	chunks, chunkEmbeddings, err := utils.SemanticChunk(text, 0.3) // Now returns both chunks and their embeddings
	if err != nil {
		return fmt.Errorf("error splitting chunks: %v", err)
	}

	// Store each chunk and its embedding in the database
	for i, chunk := range chunks {
		chunkID := fmt.Sprintf("%s_chunk_%d", docID, i)

		// Store the chunk and its embedding
		err := b.Service.StoreDocumentEmbedding(chunkID, chunk, chunkEmbeddings[i])
		if err != nil {
			return fmt.Errorf("error storing chunk %d: %v", i, err)
		}
	}

	fmt.Println("Document embedding and storage complete.")
	return nil
}*/

func (b *generalBot) StoreDocumentChunks(docID string, text string, chunkSize int, overlap int) error {
	// Chunk the document with overlap
	chunks := document.OverlapChunk(text, chunkSize, overlap)

	//client := openai.NewClient()

	for i, chunk := range chunks {
		// Get the embeddings for each chunk
		embedding, err := b.openAIclient.EmbedText(chunk)
		if err != nil {
			return fmt.Errorf("error embedding chunk %d: %v", i, err)
		}

		// Create a unique chunk ID for storage in the database
		chunkID := fmt.Sprintf("%s_chunk_%d", docID, i)
		// Store each chunk and its embedding
		b.Service.StoreDocumentEmbedding(chunkID, chunk, embedding)
	}
	fmt.Println("Document embedding complete.")
	return nil
}

/*func (b *generalBot) StoreDocumentChunks(docID string, text string, chunkSize int, minChunkSize int) error {
	// Chunk the document using the chunking logic
	chunks, err := utils.SemanticChunk(text, 0.5)
	//chunks := utils.ChunkSmartly(text, chunkSize, minChunkSize)
	if err != nil {
		return fmt.Errorf("error splitting chunks: %v", err)
	}

	for i, chunk := range chunks {
		// Get the embeddings for each chunk
		embedding, err := openai.EmbedText(chunk)
		if err != nil {
			return fmt.Errorf("error embedding chunk %d: %v", i, err)
		}

		// Create a unique chunk ID for storage in the database
		chunkID := fmt.Sprintf("%s_chunk_%d", docID, i)
		// Store each chunk and its embedding
		b.Service.StoreDocumentEmbedding(chunkID, chunk, embedding)
	}
	fmt.Println("Document embedding complete.")
	return nil
}*/

func (b *generalBot) ProcessDocument(sessionID string, filePath string) error {
	// Extract text from the uploaded file (assuming downloadAndExtractText can handle local files)
	docText, err := document.DownloadAndExtractText(filePath)
	if err != nil {
		return fmt.Errorf("error processing document: %w", err)
	}

	// Store document chunks and their embeddings
	//chunkSize := 300
	//minChunkSize := 50
	err = b.StoreDocumentChunks(sessionID, docText, b.embConfig.ChunkSize, b.embConfig.MinChunkSize)
	if err != nil {
		return fmt.Errorf("error storing document chunks: %w", err)
	}

	return nil
}
