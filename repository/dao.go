package repository

import (
	"crossplatform_chatbot/database"
	"crossplatform_chatbot/models"
	"crossplatform_chatbot/utils"
	"fmt"
	"time"
)

// DAO interface defines all necessary methods for different entities.
type DAO interface {
	CreateUser(userIDStr string, req ValidateUserReq) error
	GetUser(userIDStr string) (*models.User, error)
	CreateDocumentEmbedding(filename, docID, docText string, embedding []float64) error
	FetchEmbeddings() (map[string][]float64, map[string]string, error)
	GetAllDocuments() ([]models.Document, error)
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
func (d *dao) CreateDocumentEmbedding(filename, docID, docText string, embedding []float64) error {
	docText = utils.SanitizeText(docText)
	embeddingStr := utils.Float64SliceToPostgresArray(embedding)

	docEmbedding := models.Document{
		Filename:  filename,
		DocID:     docID,
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
			return nil, nil, fmt.Errorf("error parsing embedding for docID %s: %v", embedding.DocID, err)
		}
		documentEmbeddings[embedding.DocID] = floatSlice
		docText[embedding.DocID] = embedding.DocText
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

/*// database access object
type GormDB struct {
	DB *gorm.DB
}

func (g *GormDB) Create(value interface{}) error {
	return g.DB.Create(value).Error
}

func (g *GormDB) Where(query interface{}, args ...interface{}) Database {
	return &GormDB{DB: g.DB.Where(query, args...)}
}

func (g *GormDB) First(out interface{}, where ...interface{}) error {
	return g.DB.First(out, where...).Error
}

func (g *GormDB) Save(value interface{}) error {
	return g.DB.Save(value).Error
}

func (g *GormDB) Model(value interface{}) Database {
	g.DB = g.DB.Model(value)
	return g
}

func (g *GormDB) Take(out interface{}, where ...interface{}) error {
	return g.DB.Take(out, where...).Error
}

func (g *GormDB) Delete(value interface{}, where ...interface{}) error {
	return g.DB.Delete(value, where...).Error
}

func (g *GormDB) Find(out interface{}, where ...interface{}) error {
	return g.DB.Find(out, where...).Error
}

func (g *GormDB) Updates(values interface{}) error {
	return g.DB.Updates(values).Error
}*/

////////
// func (d *dao) CreateUser() error {
// 	// save into postgres
// 	d.database.GetDB().Create(model)
// 	// d.database.GetPostgresDB().Create(model)
// }

// func (d *dao) CreatePlayer() error {
// 	// save into mongodb
// 	d.database.GetMongoDB().Create(model)
// }

// func (d *dao) CreateTask() error {
// 	// save mysql
// 	d.database.GetMongoMySQLDB().Create(model)
// }

// func (d *dao) GetTask(id string) Task {
// 	// remote api server
// 	d.api.GetTak(id)
// }
