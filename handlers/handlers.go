package handlers

import (
	"crossplatform_chatbot/bot"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleGeneralWebhook handles incoming POST requests from the frontend
func HandlerGeneralBot(c *gin.Context, b bot.GeneralBot) {
	// Parse the incoming request from the frontend and extract the message
	/*var req struct {
		Message   string `json:"message"`
		SessionID string `json:"sessionID"`
	}

	// Bind the incoming request body to the req struct
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("failed to bind request: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind request"})
		return
	}*/

	// Delegate the handling of the message to the generalBot
	b.HandleGeneralMessage(c)

	// Return an OK status
	c.Status(http.StatusOK)
}

// HandleDocumentUpload handles document uploads, processes the document, chunks it, and stores the embeddings
func HandlerDocumentUpload(c *gin.Context, b bot.GeneralBot) {

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

	// Call bot to process the document  TODO service
	err = b.ProcessDocument(sessionID, filePath)
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
