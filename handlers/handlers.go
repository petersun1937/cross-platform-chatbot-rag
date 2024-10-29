package handlers

import (
	"crossplatform_chatbot/bot"
	"crossplatform_chatbot/models"
	"crossplatform_chatbot/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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

	// // Store the context (to use later for sending the response)
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

// // HandleGeneralWebhook handles incoming POST requests from the frontend
// func HandlerGeneralBot(c *gin.Context, b bot.GeneralBot) {
// 	// Parse the incoming request from the frontend and extract the message
// 	/*var req struct {
// 		Message   string `json:"message"`
// 		SessionID string `json:"sessionID"`
// 	}

// 	// Bind the incoming request body to the req struct
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		fmt.Printf("failed to bind request: %s\n", err.Error())
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind request"})
// 		return
// 	}*/

// 	// Delegate the handling of the message to the generalBot
// 	b.HandleGeneralMessage(c)

// 	// Return an OK status
// 	c.Status(http.StatusOK)
// }

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

	// Call bot to process the document
	err = b.ProcessDocument(file.Filename, sessionID, filePath)
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
