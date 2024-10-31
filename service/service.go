package service

import (
	"crossplatform_chatbot/bot"
	config "crossplatform_chatbot/configs"
	"crossplatform_chatbot/database"
	"crossplatform_chatbot/repository"
	"fmt"
	"log"
)

type Service struct {
	bots       map[string]bot.Bot
	database   database.Database // TODO
	repository repository.DAO
}

func NewService(botConfig config.BotConfig, embConfig config.EmbeddingConfig, db database.Database) *Service {

	dao := repository.NewDAO(db)

	return &Service{
		bots:       createBots(botConfig, embConfig, db, dao),
		database:   db,
		repository: dao,
	}
}

func (s *Service) RunBots() error {

	for _, bot := range s.bots {
		if err := bot.Run(); err != nil {
			// log.Fatal("running bot failed:", err)
			fmt.Printf("running bot failed: %s", err.Error())
			return err
		}
	}

	return nil
}

func createBots(botConfig config.BotConfig, embConfig config.EmbeddingConfig, database database.Database, dao repository.DAO) map[string]bot.Bot {
	// Initialize bots
	lineBot, err := bot.NewLineBot(botConfig, database, dao)
	if err != nil {
		//log.Fatal("Failed to initialize LINE bot:", err)
		fmt.Printf("Failed to initialize LINE bot: %s", err.Error())
	}

	tgBot, err := bot.NewTGBot(botConfig, embConfig, database, dao)
	if err != nil {
		//log.Fatal("Failed to initialize Telegram bot:", err)
		fmt.Printf("Failed to initialize Telegram bot: %s", err.Error())
	}

	fbBot, err := bot.NewFBBot(botConfig, database, dao)
	if err != nil {
		log.Fatalf("Failed to create Facebook bot: %v", err)
	}

	igBot, err := bot.NewIGBot(botConfig, database, dao)
	if err != nil {
		log.Fatalf("Failed to create Instagram bot: %v", err)
	}

	generalBot, err := bot.NewGeneralBot(botConfig, embConfig, database, dao)
	if err != nil {
		log.Fatalf("Failed to initialize General bot: %v", err)
	}

	return map[string]bot.Bot{
		"line":      lineBot,
		"telegram":  tgBot,
		"facebook":  fbBot,
		"instagram": igBot,
		"general":   generalBot,
	}
}

func (s *Service) GetBot(tag string) bot.Bot {
	return s.bots[tag]
}

type UserProfile struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ID        string `json:"id"` // Facebook User ID
}

func (s *Service) Init() error {
	// running bots
	for _, bot := range s.bots {
		if err := bot.Run(); err != nil {
			// log.Fatal("running bot failed:", err)
			fmt.Printf("running bot failed: %s", err.Error())
			return err
		}
	}
	return nil
}

func (s *Service) GetUploadedDocuments() ([]string, error) {
	documents, err := s.repository.GetAllDocuments() // Fetch all documents from the repository
	if err != nil {
		return nil, err
	}

	// Use a map to store unique filenames
	uniqueFilenameMap := make(map[string]struct{})
	//uniqueFilenames := []string{}
	uniqueFilenames := make([]string, 0, len(documents))

	// Loop through documents and collect unique filenames
	for _, doc := range documents {
		filename := doc.Filename

		// Check if the filename is a duplicate
		/*if count, exists := filenameCount[filename]; exists {
			// Increment the counter and append it to the filename to make it unique
			newFilename := fmt.Sprintf("%s(%d)", filename, count+1)
			uniqueFilenames = append(uniqueFilenames, newFilename)
			filenameCount[filename] = count + 1
		} else {
			// If it's a unique filename, store it directly
			uniqueFilenames = append(uniqueFilenames, filename)
			filenameCount[filename] = 0
		}*/

		// If filename is not already in the map, add it to the list
		if _, exists := uniqueFilenameMap[filename]; !exists {
			uniqueFilenames = append(uniqueFilenames, filename)
			uniqueFilenameMap[filename] = struct{}{} // Store it in the map
		}
	}

	return uniqueFilenames, nil
}

// // GetDB returns the gorm.DB instance from the service's database
// func (s *Service) GetDB() *gorm.DB {
// 	return s.database.GetDB()
// }

// func (s *Service) ValidateUser(userID string, req models.User) (*string, error) {
// 	_, err := s.repository.GetUser(userID)
// 	if err != nil {
// 		if err.Error() == "record not found" {
// 			if err := s.repository.CreateUser(&req); err != nil {
// 				return nil, fmt.Errorf("error creating user: %v", err)
// 			}

// 			token, err := token.GenerateToken(userID, "user")
// 			if err != nil {
// 				return nil, fmt.Errorf("error generating token: %v", err)
// 			}

// 			return &token, nil
// 		}
// 		return nil, err
// 	}
// 	return nil, nil
// }

// func (s *Service) ValidateUser(userIDStr string, req ValidateUserReq) (*string, error) {
// 	repo := NewRepository(s.database)

// 	// Check if the user exists
// 	_, err := repo.GetUser(userIDStr)
// 	if err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			// If user not found, create a new user
// 			err = repo.CreateUser(userIDStr, req)
// 			if err != nil {
// 				return nil, err
// 			}

// 			// Generate a JWT token for the new user
// 			token, err := token.GenerateToken(userIDStr, "user") // Ensure GenerateToken accepts string
// 			if err != nil {
// 				fmt.Printf("Error generating JWT: %s", err.Error())
// 				return nil, err
// 			}

// 			return &token, nil
// 		}

// 		// Other errors when fetching user
// 		return nil, err
// 	}

// 	return nil, err
// }

// // StoreDocumentEmbedding stores the document and its embedding into the database
// func (s *Service) StoreDocumentEmbedding(docID, docText string, embedding []float64) error {
// 	docText = sanitizeText(docText)
// 	embeddingStr := utils.Float64SliceToPostgresArray(embedding)

// 	docEmbedding := models.DocumentEmbedding{
// 		DocID:     docID,
// 		DocText:   docText,
// 		Embedding: embeddingStr,
// 		CreatedAt: time.Now(),
// 		UpdatedAt: time.Now(),
// 	}

// 	return s.repository.CreateDocumentEmbedding(&docEmbedding)
// }

// func sanitizeText(input string) string {
// 	validRunes := []rune{}
// 	for _, r := range input {
// 		if r == utf8.RuneError {
// 			continue // Skip invalid characters
// 		}
// 		validRunes = append(validRunes, r)
// 	}
// 	return string(validRunes)
// }

// // GetAllDocumentEmbeddings retrieves document embeddings from the repository.
// func (s *Service) GetAllDocumentEmbeddings() (map[string][]float64, map[string]string, error) {
// 	return s.repository.FetchEmbeddings()
// }

// func (s *Service) StoreDocumentEmbedding(docID string, docText string, embedding []float64) error {
// 	repo := NewRepository(s.database)

// 	// Sanitize the document text
// 	docText = sanitizeText(docText)

// 	// Convert embedding to PostgreSQL-compatible array string
// 	embeddingStr := utils.Float64SliceToPostgresArray(embedding)

// 	docEmbedding := models.DocumentEmbedding{
// 		DocID:     docID,
// 		DocText:   docText,
// 		Embedding: embeddingStr,
// 		CreatedAt: time.Now(),
// 		UpdatedAt: time.Now(),
// 	}

// 	// Store the sanitized document embedding
// 	err := repo.CreateDocumentEmbedding(&docEmbedding)
// 	if err != nil {
// 		return fmt.Errorf("error storing document embedding: %v", err)
// 	}

// 	return nil
// }

// func (s *Service) GetAllDocumentEmbeddings() (map[string][]float64, map[string]string, error) {
// 	return repository.FetchEmbeddings(s.GetDB())
// }

// func (s *Service) ValidateUser(userIDStr string, req ValidateUserReq) (*string, error) {
// 	// Check if the user exists in the database
// 	var dbUser models.User
// 	// err := s.dao.CreatePlayer()
// 	err := s.database.GetDB().Where("user_id = ? AND deleted_at IS NULL", userIDStr).First(&dbUser).Error

// 	// If the user does not exist, create a new user record
// 	if err != nil {

// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			// User does not exist, create a new user record
// 			dbUser = models.User{
// 				Model: gorm.Model{
// 					ID: 1,
// 				},
// 				UserID:       userIDStr,
// 				FirstName:    req.FirstName,
// 				LastName:     req.LastName,
// 				UserName:     req.UserName,
// 				LanguageCode: req.LanguageCode,
// 			}

// 			//err = database.DB.Create(&dbUser).Error
// 			err = s.database.GetDB().Create(&dbUser).Error

// 			if err != nil {
// 				fmt.Printf("Error creating user: %s", err.Error())
// 				return nil, err
// 			}

// 			// Generate a JWT token for the new user
// 			token, err := token.GenerateToken(userIDStr, "user") // Ensure GenerateToken accepts string
// 			if err != nil {
// 				fmt.Printf("Error generating JWT: %s", err.Error())
// 				return nil, err
// 			}

// 			return &token, nil

// 			// // Send the token to the user
// 			// msg := tgbotapi.NewMessage(message.Chat.ID, "Welcome! Your access token is: "+token)
// 			// utils.TgBot.Send(msg)
// 		} else {
// 			// Handle other types of errors
// 			fmt.Printf("Error retrieving user: %s", err.Error())
// 		}

// 	}

// 	return nil, err
// }
