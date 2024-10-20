package service

import (
	"crossplatform_chatbot/database"
	"crossplatform_chatbot/models"
	"crossplatform_chatbot/utils"
	"fmt"

	"gorm.io/gorm"
)

type Repository struct {
	database database.Database
}

func NewRepository(database database.Database) *Repository {
	return &Repository{
		database: database,
	}
}

// CreateUser creates a new user in the database
func (r *Repository) CreateUser(userIDStr string, req ValidateUserReq) error {
	//err = database.DB.Create(&dbUser).Error
	return r.database.GetDB().Create(&models.User{
		Model:        gorm.Model{},
		UserID:       userIDStr,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		UserName:     req.UserName,
		LanguageCode: req.LanguageCode,
	}).Error
}

func (r *Repository) GetUser(userIDStr string) (*models.User, error) {
	var dbUser models.User
	err := r.database.GetDB().Where("user_id = ? AND deleted_at IS NULL", userIDStr).First(&dbUser).Error
	if err != nil {
		return nil, err
	}
	return &dbUser, nil
}

// GetAllDocumentEmbeddings retrieves all document embeddings from the database
func (s *Service) GetAllDocumentEmbeddings() (map[string][]float64, map[string]string, error) {
	var embeddings []models.DocumentEmbedding
	if err := s.GetDB().Find(&embeddings).Error; err != nil {
		return nil, nil, fmt.Errorf("error retrieving document embeddings from the database: %v", err)
	}

	// Convert the retrieved embeddings into a map
	documentEmbeddings := make(map[string][]float64)
	docText := make(map[string]string) // Map to hold docID to chunk text
	for _, embedding := range embeddings {
		floatSlice, err := utils.PostgresArrayToFloat64Slice(embedding.Embedding)
		if err != nil {
			return nil, nil, fmt.Errorf("error parsing embedding for docID %s: %v", embedding.DocID, err)
		}
		documentEmbeddings[embedding.DocID] = floatSlice
		docText[embedding.DocID] = embedding.DocText // Store the chunk's text
	}

	return documentEmbeddings, docText, nil
}
