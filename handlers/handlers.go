package handlers

import (
	"crossplatform_chatbot/bot"
	"crossplatform_chatbot/models"
	"crossplatform_chatbot/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	Service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		Service: service,
	}
}

// HandlerGeneralBot handles incoming POST requests from the frontend
func (h *Handler) HandlerGeneralBot(c *gin.Context) {
	var req models.GeneralRequest

	// Parse the incoming request from the frontend and bind to the req struct
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("failed to bind request: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// // Store the context (to use later for sending the response) TODO to bot?
	b := h.Service.GetBot("general").(bot.GeneralBot)

	// Store the context using the sessionID
	b.StoreContext(req.SessionID, c)

	// Delegate the request to the service layer.
	if err := h.Service.HandleGeneral(req); err != nil {
		fmt.Printf("Error handling general request: %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to handle request"})
		return
	}

	// Return an OK status after successfully processing the message
	c.Status(http.StatusOK)
}

// HandleDocumentUpload handles document uploads, processes the document, chunks it, and stores the embeddings
func (h *Handler) HandlerDocumentUpload(c *gin.Context) {

	// Parse the file from the form-data
	file, err := c.FormFile("document")
	if err != nil {
		fmt.Printf("Error receiving file: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file upload"})
		return
	}

	// Retrieve session ID from the form-data
	sessionID := c.PostForm("sessionID")
	if sessionID == "" {
		fmt.Println("Error: Missing sessionID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	// Save the uploaded file to a temporary location
	filePath := fmt.Sprintf("/tmp/%s", file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		fmt.Printf("Error saving file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving the file"})
		return
	}

	b := h.Service.GetBot("general").(bot.GeneralBot)

	// Generate a unique document ID
	uniqueDocID := uuid.New().String()

	// Call bot to process the document
	//h.Service.handleDocumentUpload(filePath) //In service
	err = b.ProcessDocument(file.Filename, uniqueDocID, filePath)
	if err != nil {
		fmt.Printf("Error processing document: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send success response
	fmt.Println("Document processed successfully")
	c.JSON(http.StatusOK, gin.H{
		"response": "Document processed successfully",
	})
	//c.JSON(http.StatusOK, gin.H{"message": "Document processed successfully"})
}

// HandlerGetDocuments handles the retrieval of uploaded documents
func (h *Handler) HandlerGetDocuments(c *gin.Context) {
	//fmt.Println("GET request received at /api/document/list")

	// Set CORS headers explicitly for the GET request
	/*c.Header("Access-Control-Allow-Origin", "*") // Allow all origins
	c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")*/

	filenames, err := h.Service.GetUploadedDocuments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve documents"})
		return
	}

	// Return only the filenames as a JSON array
	c.JSON(http.StatusOK, gin.H{"filenames": filenames})
}
