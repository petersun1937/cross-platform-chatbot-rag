package document

import (
	"crossplatform_chatbot/openai"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// Compute similarity score and retrieve  the top N chunks from database.
func RetrieveTopNChunks(query string, documentEmbeddings map[string][]float64, topN int, docIDToText map[string]string, threshold float64) ([]string, error) {
	fmt.Println("Embedding query for similarity search...")
	client := openai.NewClient()
	queryEmbedding, err := client.EmbedText(query)
	if err != nil {
		return nil, fmt.Errorf("error embedding query: %v", err)
	}

	fmt.Println("Calculating similarity between query and document chunks...")
	type chunkScore struct {
		chunkID string
		score   float64
	}
	var scores []chunkScore

	// Calculate similarity for each document chunk
	for chunkID, embedding := range documentEmbeddings {
		//score := cosineSimilarity(queryEmbedding, embedding)
		cosineScore := cosineSimilarity(queryEmbedding, embedding)
		keywordScore := keywordMatchScore(query, docIDToText[chunkID])
		combinedScore := weightedScore(cosineScore, keywordScore)

		fmt.Printf("Combined score for chunk %s: %f\n", chunkID, combinedScore)

		// Only add the chunk if it meets the threshold
		if combinedScore >= threshold {
			scores = append(scores, chunkScore{chunkID, combinedScore})
		}
		/*if score >= similarityThreshold { // Filter out low similarity scores
			scores = append(scores, chunkScore{chunkID, score})
			fmt.Printf("Similarity score for chunk %s: %f\n", chunkID, score)
		}*/
	}

	// Sort the chunks based on score (highest score first)
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	// Collect the top N chunks' actual text using docIDToText
	var topChunksText []string
	for i := 0; i < topN && i < len(scores); i++ {
		chunkID := scores[i].chunkID
		if text, exists := docIDToText[chunkID]; exists {
			topChunksText = append(topChunksText, text)
			fmt.Printf("Top relevant chunk selected. %s: %f\n", chunkID, scores[i].score)
		} else {
			topChunksText = append(topChunksText, fmt.Sprintf("Text not found for chunk: %s", chunkID))
		}
	}

	//fmt.Println("Top relevant chunks selected. %s: %f\n", chunkID, combinedScore)
	return topChunksText, nil
}

// Chunks the text by full sentences, keeping each chunk under a certain word limit
func ChunkDocumentBySentence(text string, chunkSize int) []string {
	sentences := splitIntoSentences(text) // Split the document into sentences
	var chunks []string
	var currentChunk []string
	currentWordCount := 0

	// Iterate over the sentences and add them to chunks
	for _, sentence := range sentences {
		wordCount := len(strings.Fields(sentence)) // Count the words in the sentence

		// If adding this sentence exceeds the chunk size, start a new chunk
		if currentWordCount+wordCount > chunkSize && len(currentChunk) > 0 {
			chunks = append(chunks, strings.Join(currentChunk, " "))
			currentChunk = []string{}
			currentWordCount = 0
		}

		currentChunk = append(currentChunk, sentence)
		currentWordCount += wordCount
	}

	// Add the last chunk if it's not empty
	if len(currentChunk) > 0 {
		chunks = append(chunks, strings.Join(currentChunk, " "))
	}

	return chunks
}

// Helper function to split text into paragraphs
func splitIntoParagraphs(text string) []string {
	// Split text into paragraphs using two or more newlines as the delimiter
	return strings.Split(text, "\n\n")
}

// Helper function to check if a string is an enumerated point (e.g., "1. ", "2. ", etc.)
func isEnumeratedPoint(line string) bool {
	re := regexp.MustCompile(`^\d+\.\s`)
	return re.MatchString(line)
}

// Helper function to chunk text into paragraphs and points
func ChunkSmartly(text string, maxChunkSize int, minWordsPerChunk int) []string {
	var chunks []string

	// Split the text into paragraphs first
	paragraphs := splitIntoParagraphs(text)
	for _, paragraph := range paragraphs {
		// Split paragraph into lines to detect enumerated points
		lines := strings.Split(paragraph, "\n")
		pointGroup := ""
		for _, line := range lines {
			line = strings.TrimSpace(line) // Ensure that each line is trimmed
			if line == "" {
				continue // Skip empty lines
			}

			// If the line starts with an enumerated point, group it
			if isEnumeratedPoint(line) {
				if pointGroup != "" {
					// Add the previous point group to chunks
					chunks = append(chunks, strings.TrimSpace(pointGroup))
					pointGroup = ""
				}
				pointGroup += line + " " // Group lines with enumerated points
			} else {
				// If a non-enumerated point, check if it's a new paragraph or continuation
				if pointGroup != "" {
					pointGroup += line + " "
				} else {
					// Add standalone sentences or paragraphs directly
					if len(line) <= maxChunkSize {
						chunks = append(chunks, strings.TrimSpace(line))
					} else {
						// Split into smaller sentences if it's still too long
						chunks = append(chunks, splitIntoSentences(line)...)
					}
				}
			}
		}
		if pointGroup != "" {
			chunks = append(chunks, strings.TrimSpace(pointGroup)) // Add the last grouped point
		}
	}

	// Combine chunks that are smaller than the minimum word count
	var combinedChunks []string
	currentChunk := ""
	currentWordCount := 0

	for _, chunk := range chunks {
		wordsInChunk := len(strings.Fields(chunk))

		// If adding this chunk doesn't exceed the minimum word count, combine it with the current chunk
		if currentWordCount+wordsInChunk < minWordsPerChunk {
			currentChunk += " " + chunk
			currentWordCount += wordsInChunk
		} else {
			// If current chunk meets the minimum size, append it to combinedChunks
			if len(strings.TrimSpace(currentChunk)) > 0 {
				combinedChunks = append(combinedChunks, strings.TrimSpace(currentChunk))
			}
			// Start a new chunk
			currentChunk = chunk
			currentWordCount = wordsInChunk
		}
	}

	// Add the last chunk if it hasn't been added yet
	if len(strings.TrimSpace(currentChunk)) > 0 {
		combinedChunks = append(combinedChunks, strings.TrimSpace(currentChunk))
	}

	return combinedChunks
}

// Helper function to split text into sentences, accounting for decimals
func splitIntoSentences(text string) []string {
	var sentences []string
	sentence := ""
	for i, r := range text {
		sentence += string(r)

		// Check for sentence-ending punctuation ('.', '!', '?')
		if r == '.' || r == '!' || r == '?' {
			// Ensure it's not a decimal point by checking if the character before the period is a digit
			if r == '.' && i > 0 && unicode.IsDigit(rune(text[i-1])) {
				continue // Skip splitting here if it's part of a number (e.g., 1.1, 1.2)
			}

			// Add the sentence to the list and reset the sentence accumulator
			sentences = append(sentences, strings.TrimSpace(sentence))
			sentence = ""
		}
	}

	// Add any remaining text as the last sentence
	if len(strings.TrimSpace(sentence)) > 0 {
		sentences = append(sentences, strings.TrimSpace(sentence))
	}

	return sentences
}

// SemanticChunk performs semantic chunking, embedding sentences and chunking them based on similarity
// func SemanticChunk(text string, similarityThreshold float64) ([]string, [][]float64, error) {
// 	// Step 1: Split the document into sentences or smaller units
// 	sentences := strings.Split(text, ". ")

// 	// Step 2: Embed all sentences in one batch
// 	embeddings, err := openai.EmbedSentencesBatch(sentences)
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("error embedding sentences: %v", err)
// 	}

// 	// Step 3: Initialize variables for chunking
// 	var chunks []string
// 	var chunkEmbeddings [][]float64 // To store the embeddings for each chunk
// 	var currentChunk []string
// 	var currentEmbedding []float64
// 	prevEmbedding := embeddings[0]

// 	// Add the first sentence to the current chunk
// 	currentChunk = append(currentChunk, sentences[0])
// 	currentEmbedding = append(currentEmbedding, embeddings[0]...)

// 	// Step 4: Loop through the remaining sentences
// 	for i := 1; i < len(sentences); i++ {
// 		currEmbedding := embeddings[i]

// 		// Calculate similarity between consecutive sentence embeddings
// 		similarity := cosineSimilarity(prevEmbedding, currEmbedding)

// 		// If the similarity drops below the threshold, create a new chunk
// 		if similarity < similarityThreshold {
// 			// Add the current chunk to the list
// 			chunks = append(chunks, strings.Join(currentChunk, ". "))
// 			chunkEmbeddings = append(chunkEmbeddings, currentEmbedding)

// 			// Start a new chunk
// 			currentChunk = []string{sentences[i]}
// 			currentEmbedding = embeddings[i]
// 		} else {
// 			// Continue adding to the current chunk
// 			currentChunk = append(currentChunk, sentences[i])
// 			currentEmbedding = append(currentEmbedding, embeddings[i]...)
// 		}

// 		// Update previous embedding for the next comparison
// 		prevEmbedding = currEmbedding
// 	}

// 	// Step 5: Add the last chunk
// 	if len(currentChunk) > 0 {
// 		chunks = append(chunks, strings.Join(currentChunk, ". "))
// 		chunkEmbeddings = append(chunkEmbeddings, currentEmbedding)
// 	}

// 	// Return the chunks and their corresponding embeddings
// 	return chunks, chunkEmbeddings, nil
// }

// Overlap chunking to preserve context across chunks
func OverlapChunk(text string, chunkSize int, overlap int) []string {
	words := strings.Fields(text) // Split text into words
	var chunks []string

	for i := 0; i < len(words); i += chunkSize - overlap {
		end := i + chunkSize
		if end > len(words) {
			end = len(words)
		}
		chunk := strings.Join(words[i:end], " ")
		chunks = append(chunks, chunk)
	}

	return chunks
}
