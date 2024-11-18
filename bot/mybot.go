package bot

import (
	config "crossplatform_chatbot/configs"
	"crossplatform_chatbot/database"
	document "crossplatform_chatbot/document_proc"
	openai "crossplatform_chatbot/openai"
	"crossplatform_chatbot/repository"
	"fmt"
	"net/http"
	"strings"

	"cloud.google.com/go/dialogflow/apiv2/dialogflowpb"
	"github.com/gin-gonic/gin"
)

type GeneralBot interface {
	Run() error
	HandleGeneralMessage(sessionID, message string)
	//sendResponse(identifier interface{}, response string) error
	StoreDocumentChunks(Filename, docID, text string, chunkSize, minchunkSize int) error
	ProcessDocument(Filename, sessionID, filePath string) error
	StoreContext(sessionID string, c *gin.Context)
	//SetWebhook(webhookURL string) error
}

type generalBot struct {
	// Add any common fields if necessary, like configuration
	BaseBot
	//ctx context.Context
	// conf         config.BotConfig
	//embConfig    config.EmbeddingConfig
	//openAIclient *openai.Client
	//config map[string]string
}

// func NewGeneralBot(conf *config.Config, service *service.Service) (*generalBot, error) {
// 	baseBot := &BaseBot{
// 		Platform: GENERAL,
// 		Service:  service,
// 	}

// 	return &generalBot{
// 		BaseBot:      baseBot,
// 		conf:         conf.BotConfig,
// 		embConfig:    conf.EmbeddingConfig,
// 		ctx:          context.Background(),
// 		openAIclient: openai.NewClient(),
// 	}, nil
// }

// creates a new GeneralBot instance
func NewGeneralBot(botconf config.BotConfig, embconf config.EmbeddingConfig, database database.Database, dao repository.DAO) (*generalBot, error) {

	return &generalBot{
		BaseBot: BaseBot{
			Platform:     GENERAL,
			conf:         botconf,
			database:     database,
			dao:          dao,
			openAIclient: openai.NewClient(),
			embConfig:    embconf,
		},
	}, nil
}

func (b *generalBot) Run() error {
	// Implement logic for running the bot
	fmt.Println("General bot is running...")
	return nil
}

// func (b *generalBot) HandleGeneralMessage(c *gin.Context) {
func (b *generalBot) HandleGeneralMessage(sessionID, message string) {

	// Process and send the message
	b.ProcessUserMessage(sessionID, message)

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
		documentEmbeddings, chunkText, err := b.BaseBot.dao.FetchEmbeddings()
		//documentEmbeddings, chunkText, err := b.Service.GetAllDocumentEmbeddings()
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
				response, err = b.BaseBot.GetOpenAIResponse(gptPrompt)
				if err != nil {
					response = fmt.Sprintf("OpenAI Error: %v", err)
				}
			} else {
				// If no relevant document found, fallback to OpenAI response
				response, err = b.BaseBot.GetOpenAIResponse(message)
				if err != nil {
					response = fmt.Sprintf("OpenAI Error: %v", err)
				}
			}

		} else {
			//response = fmt.Sprintf("You said: %s", message)
			//HandleDialogflowIntent(message string) (string, error) {
			b.BaseBot.handleMessageDialogflow(GENERAL, sessionID, message, b)
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

func (b *generalBot) sendFrontendMessage(c *gin.Context, message string) error {
	if c == nil {
		return fmt.Errorf("gin context is nil")
	}
	c.JSON(http.StatusOK, gin.H{
		"response": message,
	})
	return nil
}

// TODO
var sessionContextMap = make(map[string]*gin.Context)

// StoreContext stores the context in sessionContextMap using the session ID
func (b *generalBot) StoreContext(sessionID string, c *gin.Context) {
	sessionContextMap[sessionID] = c
}

// Retrieve the context using sessionID when you need to send a response
func getContext(sessionID string) (*gin.Context, error) {
	if context, ok := sessionContextMap[sessionID]; ok {
		return context, nil
	}
	return nil, fmt.Errorf("no context found for session ID %s", sessionID)
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

func (b *generalBot) ProcessDocument(filename, sessionID, filePath string) error {
	// Extract text from the uploaded file (assuming downloadAndExtractText can handle local files)
	docText, err := document.DownloadAndExtractText(filePath)
	if err != nil {
		return fmt.Errorf("error processing document: %w", err)
	}

	// Store document chunks and their embeddings
	//chunkSize := 300
	//minChunkSize := 50
	err = b.StoreDocumentChunks(filename, filename+"_"+sessionID, docText, b.embConfig.ChunkSize, b.embConfig.MinChunkSize)
	if err != nil {
		return fmt.Errorf("error storing document chunks: %w", err)
	}

	// Combine chunk embeddings and auto-tag using embeddings
	// tags, err := document.AutoTagDocumentEmbeddings(sessionID, b.openAIclient, b.BaseBot.dao, b.embConfig.TagEmbeddings)
	// if err != nil {
	// 	return fmt.Errorf("error auto-tagging document: %w", err)
	// }

	//TODO
	// Retrieve all chunk embeddings for the document
	// chunkEmbeddings, err := b.BaseBot.dao.GetChunkEmbeddings(sessionID)
	// if err != nil {
	// 	return fmt.Errorf("error retrieving chunk embeddings: %w", err)
	// }

	// // Combine chunk embeddings into a single document embedding
	// combinedEmbedding, err := utils.AverageEmbeddings(chunkEmbeddings)
	// if err != nil {
	// 	return fmt.Errorf("error combining chunk embeddings: %w", err)
	// }

	// // Retrieve tags based on similarity between document embedding and tag embeddings
	// tags := document.GetRelevantTags(combinedEmbedding, b.embConfig.TagEmbeddings, 0.7)
	// fmt.Println("Auto-tagged with tags:", tags)

	// Auto-tagging using OpenAI
	tags, err := b.BaseBot.openAIclient.AutoTagWithOpenAI(docText)
	if err != nil {
		return fmt.Errorf("error auto-tagging document: %w", err)
	}

	// Save tags in document metadata
	if err := b.BaseBot.dao.SaveDocumentMetadata(sessionID, tags); err != nil {
		return fmt.Errorf("error saving document metadata: %w", err)
	}

	return nil
}

func (b *generalBot) StoreDocumentChunks(filename, docID, text string, chunkSize, overlap int) error {
	// Chunk the document with overlap
	chunks := document.OverlapChunk(text, chunkSize, overlap)

	//client := openai.NewClient()

	for i, chunk := range chunks {
		// Get the embeddings for each chunk
		embedding, err := b.BaseBot.openAIclient.EmbedText(chunk)
		if err != nil {
			return fmt.Errorf("error embedding chunk %d: %v", i, err)
		}

		// Create a unique chunk ID for storage in the database
		chunkID := fmt.Sprintf("%s_chunk_%d", docID, i)
		// Store each chunk and its embedding
		err = b.BaseBot.dao.CreateDocumentEmbedding(filename, docID, chunkID, chunk, embedding) // Store each chunk with its embedding
		if err != nil {
			return fmt.Errorf("error storing chunks: %v", err)
		}
		//b.Service.StoreDocumentEmbedding(chunkID, chunk, embedding)
	}
	fmt.Println("Document embedding complete.")
	return nil
}
