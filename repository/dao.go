package repository

import (
	"crossplatform_chatbot/database"
	"crossplatform_chatbot/models"
	"crossplatform_chatbot/utils"
	"fmt"
	"time"

	"github.com/lib/pq"
)

// DAO interface defines all necessary methods for different entities.
type DAO interface {
	CreateUser(userIDStr string, req ValidateUserReq) error
	GetUser(userIDStr string) (*models.User, error)
	CreateDocumentEmbedding(filename, docID, chunkID, docText string, embedding []float64) error
	FetchEmbeddings() (map[string][]float64, map[string]string, error)
	GetAllDocuments() ([]models.Document, error)
	//SaveDocumentMetadata(docID string, tags []string) error
	GetChunkEmbeddings(docID string) ([][]float64, error)
	RetrieveTagEmbeddings() (map[string][]float64, error)
	StoreTagEmbeddings(tagDescriptions map[string]string, embedFunc func(string) ([]float64, error)) error
	GetDocumentChunksByTags(tags []string) ([]models.Document, error)
}

// dao struct implements the DAO interface.
type dao struct {
	db database.Database
}

// type dao struct {
// 	database database.Database
// }

// NewDAO initializes and returns a new DAO instance.
func NewDAO(database database.Database) DAO {
	return &dao{
		db: database,
	}
}

type ValidateUserReq struct {
	FirstName    string
	LastName     string
	UserName     string
	LanguageCode string
}

// CreateUser inserts a new user into the database.
func (d *dao) CreateUser(userIDStr string, req ValidateUserReq) error {
	user := models.User{
		UserID:       userIDStr,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		UserName:     req.UserName,
		LanguageCode: req.LanguageCode,
	}
	return d.db.GetDB().Create(&user).Error
}

// Retrieves a user from the database.
func (d *dao) GetUser(userIDStr string) (*models.User, error) {
	var user models.User
	err := d.db.GetDB().Where("user_id = ? AND deleted_at IS NULL", userIDStr).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateDocumentEmbedding stores the document and its embedding into the database.
func (d *dao) CreateDocumentEmbedding(filename, docID, chunkID, docText string, embedding []float64) error {
	docText = utils.SanitizeText(docText)
	embeddingStr := utils.Float64SliceToPostgresArray(embedding)

	docEmbedding := models.Document{
		Filename:  filename,
		DocID:     docID,
		ChunkID:   chunkID,
		DocText:   docText,
		Embedding: embeddingStr,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return d.db.GetDB().Create(&docEmbedding).Error
}

// FetchEmbeddings retrieves and converts all document embeddings from the database.
func (d *dao) FetchEmbeddings() (map[string][]float64, map[string]string, error) {
	var embeddings []models.Document

	if err := d.db.GetDB().Find(&embeddings).Error; err != nil {
		return nil, nil, fmt.Errorf("error retrieving embeddings: %v", err)
	}

	documentEmbeddings := make(map[string][]float64)
	docText := make(map[string]string)

	for _, embedding := range embeddings {
		floatSlice, err := utils.PostgresArrayToFloat64Slice(embedding.Embedding)
		if err != nil {
			return nil, nil, fmt.Errorf("error parsing embedding for docID %s: %v", embedding.ChunkID, err)
		}
		documentEmbeddings[embedding.ChunkID] = floatSlice
		docText[embedding.ChunkID] = embedding.DocText
	}

	return documentEmbeddings, docText, nil
}

// GetAllDocuments retrieves all uploaded documents from the database.
func (d *dao) GetAllDocuments() ([]models.Document, error) {
	var documents []models.Document

	// Fetch all documents from the database.
	if err := d.db.GetDB().Find(&documents).Error; err != nil {
		return nil, fmt.Errorf("error retrieving documents: %v", err)
	}

	return documents, nil
}

func (d *dao) GetChunkEmbeddings(docID string) ([][]float64, error) {
	var chunks []models.Document
	if err := d.db.GetDB().Where("doc_id = ?", docID).Find(&chunks).Error; err != nil {
		return nil, fmt.Errorf("error retrieving chunks for docID %s: %w", docID, err)
	}

	var embeddings [][]float64
	for _, chunk := range chunks {
		embedding, err := utils.ParseEmbeddingString(chunk.Embedding)
		if err != nil {
			return nil, fmt.Errorf("error parsing embedding for chunk %s: %w", chunk.ChunkID, err)
		}
		embeddings = append(embeddings, embedding)
	}

	return embeddings, nil
}

// RetrieveTagEmbeddings gets embeddings for tags and stores from the database
func (d *dao) RetrieveTagEmbeddings() (map[string][]float64, error) {
	var tagEmbeddings []models.TagEmbedding

	// Use GORM's Raw query to fetch tag embeddings
	if err := d.db.GetDB().Raw("SELECT tag_name, embedding FROM tag_embeddings").Scan(&tagEmbeddings).Error; err != nil {
		return nil, fmt.Errorf("error retrieving tag embeddings: %v", err)
	}

	// Convert to map format
	embeddingsMap := make(map[string][]float64)
	for _, tagEmbedding := range tagEmbeddings {
		embeddingsMap[tagEmbedding.TagName] = tagEmbedding.Embedding
	}

	return embeddingsMap, nil
}

// StoreTagEmbeddings generates embeddings for tags and stores them in the database
func (d *dao) StoreTagEmbeddings(tagDescriptions map[string]string, embedFunc func(string) ([]float64, error)) error {
	for tag, description := range tagDescriptions {
		// Generate the embedding for each tag using the provided embed function
		embedding, err := embedFunc(description)
		if err != nil {
			return fmt.Errorf("error generating embedding for tag %s: %v", tag, err)
		}

		// Insert the tag and embedding into the database
		query := `INSERT INTO tag_embeddings (tag_name, embedding) VALUES ($1, $2) ON CONFLICT (tag_name) DO NOTHING`
		if err := d.db.GetDB().Exec(query, tag, embedding).Error; err != nil {
			return fmt.Errorf("error inserting tag embedding for %s: %v", tag, err)
		}
	}
	return nil
}

// GetDocumentChunksByTags retrieves document chunks matching the specified tags and decodes their embeddings.
func (d *dao) GetDocumentChunksByTags(tags []string) ([]models.Document, error) {
	var docIDs []string

	// Query the document_metadata table to get doc_ids where any of the tags match
	err := d.db.GetDB().Table("document_metadata").
		Where("tags && ?::text[]", pq.Array(tags)). // Use pq.Array and cast to text[]
		Pluck("doc_id", &docIDs).Error
	if err != nil {
		return nil, fmt.Errorf("error retrieving doc_ids by tags: %w", err)
	}

	if len(docIDs) == 0 {
		return nil, nil // No matching documents found for the tags
	}

	// Query the documents table to get the document chunks for the retrieved doc_ids
	var documents []models.Document
	err = d.db.GetDB().Where("doc_id IN ?", docIDs).Find(&documents).Error
	if err != nil {
		return nil, fmt.Errorf("error retrieving document chunks: %w", err)
	}

	// Decode embeddings and replace the string with the decoded embeddings
	for i, doc := range documents {
		embedding, err := utils.PostgresArrayToFloat64Slice(doc.Embedding)
		if err != nil {
			return nil, fmt.Errorf("error decoding embedding for document: %w", err)
		}
		documents[i].Embedding = utils.Float64SliceToPostgresArray(embedding) // Encode back if needed, or just use decoded embedding downstream
	}

	return documents, nil
}
