package document_proc

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/nguyenthenguyen/docx"
)

// downloadAndExtractText handles both local files and remote URLs, including .txt, .docx, and .pdf
func DownloadAndExtractText(filePathOrURL string) (string, error) {
	// Check the file type based on the extension
	if strings.HasSuffix(filePathOrURL, ".txt") {
		return readLocalFile(filePathOrURL) // TXT
	} else if strings.HasSuffix(filePathOrURL, ".docx") {
		return extractTextFromDocx(filePathOrURL) // DOCX
	} else if strings.HasSuffix(filePathOrURL, ".pdf") {
		return extractTextFromPDF(filePathOrURL) // PDF
		//return extractTextFromPDFInChunks(filePathOrURL)
	}

	return "", fmt.Errorf("unsupported file type")
}

// readLocalFile reads the content of a .txt file
func readLocalFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening local file: %v", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("error reading local file content: %v", err)
	}

	return string(content), nil
}

// extractTextFromDocx extracts text from a DOCX file
func extractTextFromDocx(filePath string) (string, error) {
	// Open the docx file
	r, err := docx.ReadDocxFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening docx file: %v", err)
	}
	defer r.Close()

	// Extract the raw content from the docx file
	docx1 := r.Editable()
	rawContent := docx1.GetContent()

	// Clean the extracted content by removing XML tags
	cleanContent := cleanDocxText(rawContent)

	return cleanContent, nil
}

// cleanDocxText filters the unnecessary XML tags for cleaner text
func cleanDocxText(rawText string) string {
	// Regular expression to match XML tags
	re := regexp.MustCompile("<[^>]*>")

	// Replace all XML tags with an empty string
	cleanText := re.ReplaceAllString(rawText, "")

	// Trim any excessive white spaces
	cleanText = strings.TrimSpace(cleanText)

	return cleanText
}

// extractTextFromPDF extracts text from a PDF file using pdfcpu (discarded, use python instead)
/*func extractTextFromPDF(filePath string) (string, error) {
	// Initialize a new configuration for pdfcpu
	conf := model.NewDefaultConfiguration()

	// Open the PDF file
	f, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening file: %v", err)
	}
	defer f.Close()

	// Get the current working directory
	outputDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %v", err)
	}

	// Generate dynamic file name using timestamp
	fileName := fmt.Sprintf("extracted_%d", time.Now().Unix())

	// Extract content from the PDF (you can specify page ranges, e.g., "1-3")
	err = api.ExtractContent(f, outputDir, fileName, nil, conf) // Extract for all pages
	if err != nil {
		return "", fmt.Errorf("error extracting content from PDF: %v", err)
	}

	// Assuming the extracted text is saved as files in the output directory, you can read them here
	var extractedText string

	// Read the extracted files and concatenate their content
	for i := 1; ; i++ {
		pageFile := fmt.Sprintf("%s/%s_page_%d.txt", outputDir, fileName, i)
		content, err := os.ReadFile(pageFile)
		if err != nil {
			// Stop when no more pages are found
			if os.IsNotExist(err) {
				break
			}
			return "", fmt.Errorf("error reading extracted text file: %v", err)
		}
		extractedText += string(content) + "\n"
	}

	// Post-process the extracted text to filter out unnecessary lines
	cleanedText := postProcessPDFContent(extractedText)

	return cleanedText, nil
}*/

// Extracts text from a PDF file using python script (pdfplumber)
func extractTextFromPDF(filePath string) (string, error) {
	// Call the Python script (pdfplumber)
	cmd := exec.Command("python", "./python_scripts/extractPDF.py", filePath)
	//cmd := exec.Command("python3", "./python_scripts/extractPDF.py", filePath)

	// Capture the output of the Python script
	var out bytes.Buffer
	cmd.Stdout = &out

	// Run the command and check for errors
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error running Python script: %v", err)
	}

	// Return the extracted text
	return out.String(), nil
}

//var documentEmbeddings = make(map[string][]float64)

/*func ChunkDocument(text string, chunkSize int) []string {
	words := strings.Fields(text) // Split the document into words
	var chunks []string
	for i := 0; i < len(words); i += chunkSize {
		end := i + chunkSize
		if end > len(words) {
			end = len(words)
		}
		chunks = append(chunks, strings.Join(words[i:end], " "))
	}
	fmt.Printf("Document chunked into %d chunks.\n", len(chunks)) // Log the chunk count
	return chunks
}*/
