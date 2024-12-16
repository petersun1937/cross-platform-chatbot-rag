package service

import (
	document "crossplatform_chatbot/document_proc"
	"crossplatform_chatbot/models"
	"crossplatform_chatbot/utils"
	"fmt"

	"gorm.io/gorm"
)

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

		// If filename is not already in the map, add it to the list
		if _, exists := uniqueFilenameMap[filename]; !exists {
			uniqueFilenames = append(uniqueFilenames, filename)
			uniqueFilenameMap[filename] = struct{}{} // Store it in the map
		}
	}

	return uniqueFilenames, nil
}

func (s *Service) HandleDocumentUpload(filename, fileID, filePath string) error {
	// step 1: call bot to process documents
	//b := s.GetBot("general").(bot.GeneralBot)

	//documents, tags, err := b.ProcessDocument(filename, fileID, filePath)
	documents, tags, err := s.ProcessDocument(filename, fileID, filePath)
	if err != nil {
		return err
	}

	// dao version
	// return s.repository.CreateDocumentsAndMeta(fileID, documents, tags)

	// service version
	// step 2: make db data
	documentModels := make([]*models.Document, 0)
	documentID := ""
	for _, doc := range documents {
		model := models.Document{
			Filename:  doc.Filename,
			DocID:     doc.DocID,
			ChunkID:   doc.ChunkID,
			DocText:   doc.DocText,
			Embedding: doc.Embedding,
		}
		documentModels = append(documentModels, &model)
		documentID = doc.DocID
	}
	metadata := models.DocumentMetadata{
		DocID: documentID,
		Tags:  tags,
	}

	// step 3: do transaction
	return s.database.GetDB().Transaction(func(tx *gorm.DB) error {
		// batch insert Documents
		if err := tx.Create(documentModels).Error; err != nil {
			return err
		}

		// insert DocumentMetadata
		if err := tx.Create(&metadata).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *Service) ProcessDocument(filename, sessionID, filePath string) ([]models.Document, []string, error) {
	// Extract text from the uploaded file
	docText, err := document.DownloadAndExtractText(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("error processing document: %w", err)
	}

	chunks := document.OverlapChunk(docText, s.embConfig.ChunkSize, s.embConfig.OverlapSize)
	documents := make([]models.Document, 0)
	tagList := []string{}

	for i, chunk := range chunks {
		document, err := s.StoreDocumentChunks(filename, fmt.Sprintf("%s_%s", filename, sessionID), chunk, i)
		if err != nil {
			return nil, nil, err
		}
		documents = append(documents, document)
		// Auto-tagging using OpenAI
		tags, err := s.client.AutoTagWithOpenAI(chunk)
		if err != nil {
			return nil, nil, fmt.Errorf("error auto-tagging document: %w", err)
		}
		tagList = append(tagList, tags...)
	}

	// Remove duplicates from the tag list
	return documents, utils.RemoveDuplicates(tagList), nil
}

func (s *Service) StoreDocumentChunks(filename, docID, chunkText string, chunkID int) (models.Document, error) {
	embedding, err := s.client.EmbedText(chunkText)
	if err != nil {
		return models.Document{}, fmt.Errorf("error embedding chunk: %v", err)
	}

	document := models.Document{
		Filename:  filename,
		DocID:     docID,
		ChunkID:   fmt.Sprintf("%s_chunk_%d", docID, chunkID),
		DocText:   utils.SanitizeText(chunkText),
		Embedding: utils.Float64SliceToPostgresArray(embedding),
	}

	/*err = s.repository.CreateDocumentEmbedding(document)
	if err != nil {
		return models.Document{}, fmt.Errorf("error storing document chunk: %v", err)
	}*/

	return document, nil
}
