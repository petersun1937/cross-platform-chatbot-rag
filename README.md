
# CrossPlatform-TechSupport-Chatbot
A multi-platform chatbot that provides intelligent customer and tech support using OpenAI, Dialogflow, and a Retrieval-Augmented Generation (RAG) system.

## Table of Contents
- [Features](#features)
- [Demo](#demo)
- [How It Works](#how-it-works)
- [Installation](#installation)
- [Tech Stack](#tech-stack)
- [License](#license)
- [Contact](#contact)

## Features
- **Multi-platform Support**: Messenger, Telegram, LINE, Instagram, or custom web page.
- **Retrieval-Augmented Generation (RAG)**: Intelligent, document-based responses using semantic search.
- **Dynamic Intent Handling**: Dialogflow for intent matching and tagging documents for improved efficiency.
- **Context-aware Responses**: Combines RAG and OpenAI models for enhanced conversational AI.
- **Document Processing**: Upload, chunk, and store documents with embeddings for semantic search.

<!---   Handles FAQs, troubleshooting, and customer inquiries -->

## Demo
<!---   - **Video Demo**: [Coming Soon] -->
- **Presentation slides**: [Support chatbot with RAG](https://docs.google.com/presentation/d/10M90QfSjpdLvMHcqu3oT6lyrRwAfknDK/view)
- **Live Demo**: [Chatbot demo with custom frontend](https://petersun1937.github.io/Custom_Frontend_Chatbot)

## How It Works
1. **Intent Handling**:
   - Dialogflow matches user inputs with known intents and tags uploaded documents.
   - Only documents matching specific tags are searched to improve efficiency.
2. **RAG Process**:
   - Uploaded documents are chunked with overlapping sections.
   - Embeddings are generated and stored for semantic search.
   - Relevant chunks are retrieved using a weighted combination of cosine similarity and fuzzy matching scores.
   - Retrieved context is added to prompts for response generation using GPT models.
   Persistent Conversation Context:
3. **Persistent Conversation Context**:
   - Redis stores user conversation history in key-value pairs, allowing personalized, context-aware responses.
   - History is fetched and included in prompts for OpenAI and Dialogflow, ensuring continuity across interactions.
4. **Cross-platform Integration**:
   - APIs for Messenger, LINE, Telegram, and Instagram.
   - Custom web frontend built with React.


## Tech Stack
- **Frontend**: React
- **Backend**: Go (Gin framework)
- **Database**: PostgreSQL
- **In-memory Store**: Redis (for conversation history and context storage)
- **APIs**: OpenAI, Dialogflow, Telegram, LINE, META
- **Cloud Deployment**: Google Cloud (backend), GitHub Pages (frontend)
- **Tools**: PDF processing libraries for text extraction



## Installation

- **Clone the Repository**:
   ```bash
   git clone https://github.com/petersun1937/cross-platform-chatbot-rag.git
   cd cross-platform-chatbot-rag
   ```

- **Backend Setup**:
1. **Prerequisites**:
   - Install Go (v1.16 or higher) and Python (v3.6 or higher).
   - Verify installation:
     ```bash
     go version
     python --version
     ```
2. **Install Python Packages**:
   - The backend requires several Python packages for processing PDFs. To install them, run:
   ```bash
   pip install pdfplumber pytesseract pdf2image PyPDF2
   ```
   - These packages handle PDF text extraction and Optical Character Recognition (OCR) for images within PDFs.

3. **Set Up Environment Variables**:
   - Create a `.env` file in the `configs/` directory to store variables such as API keys and database configurations (refer to `sample.env`).

4. **Run the Backend Server**:
   - Navigate to the project directory and start the server:
     ```bash
     go run main.go
     ```
      
### Frontend Setup (Optional)
   - Refer to [my custom frontend repo](https://github.com/petersun1937/Custom_Frontend_Chatbot)
<!--- 
   - **Tesseract OCR Installation** (Optional for OCR capabilities):
      - **Linux**: Install Tesseract via the package manager:
        ```bash
        sudo apt-get install tesseract-ocr
        ```
      - **Windows**: Download and install [Tesseract](https://github.com/tesseract-ocr/tesseract/wiki).
      - Ensure Tesseract is accessible through your system's PATH.
-->

<!---        
- Example `.env` file structure:
        ```bash
        DATABASE_URL=your_database_url
        API_KEY=your_api_key
        ```
      - Replace `your_database_url` and `your_api_key` with actual values.

   - **Run the backend server**:
      - Start the Go server by navigating to the project directory and running:
        ```bash
        go run main.go 
        ```
-->
        
<!--- 
- **Frontend Setup**:
   - Navigate to `React_custom_frontend/`:
     ```bash
     cd frontend
     npm install
     npm start
     ```
   - The frontend will run at [http://localhost:3000](http://localhost:3000).
   -->

## Potential Use Cases
- Streamlining FAQs and troubleshooting for customer support.
- Automating knowledge retrieval for internal teams.
- Enhancing collaborative workflows with intelligent document processing.
- Acting as a meeting assistant for document organization and context provision.
- Extending to incorporate custom-trained language models.


## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contact
If you have any questions or suggestions, feel free to reach out to cxs1937@psu.edu.