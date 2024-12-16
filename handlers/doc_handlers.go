package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// HandlerGetDocuments handles the retrieval of uploaded documents
func (h *Handler) HandlerGetDocuments(c *gin.Context) {
	//fmt.Println("GET request received at /api/document/list")

	// Set CORS headers explicitly for the GET request
	c.Header("Access-Control-Allow-Origin", "*") // Allow all origins
	c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

	filenames, err := h.Service.GetUploadedDocuments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve documents"})
		return
	}

	// Return only the filenames as a JSON array
	c.JSON(http.StatusOK, gin.H{"filenames": filenames})
}

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

	// Generate a unique document ID
	fileID := uuid.New().String()

	if err := h.Service.HandleDocumentUpload(file.Filename, fileID, filePath); err != nil {
		fmt.Printf("Error processing document: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
