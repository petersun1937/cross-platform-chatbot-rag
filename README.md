
# CrossPlatform-TechSupport-Chatbot
A multi-platform chatbot that provides intelligent customer and tech support using OpenAI, Dialogflow, and a Retrieval-Augmented Generation (RAG) system.

## Table of Contents
- [Features](#features)
- [Demo](#demo)
- [How It Works](#how-it-works)
- [Installation](#installation)
- [Tech Stack](#tech-stack)
- [License](#license)

## Features
- Multi-platform support (Messenger, Telegram, LINE, or custom web page)
- Document storing and chunking
- OpenAI integration for conversational AI and text embedding/semantic search
- Retrieval-Augmented Generation (RAG) for document-based responses
- Context-aware responses and smart routing

<!---   Handles FAQs, troubleshooting, and customer inquiries -->

## Demo
<!---   - **Video Demo**: [Coming Soon] -->
- **Live Demo**: https://petersun1937.github.io/Custom_Frontend_Chatbot

## How It Works
- Users interact with the chatbot through various platforms (support FB Messenger, Telegram, LINE, or custom platform).
- Dialogflow detects intents for common questions and support requests.
- OpenAI generates conversational responses for unrecognized inputs.
- Document upload, which is then chunked and stored in the database along with text embeddings.
- The Retrieval-Augmented Generation (RAG) system fetches relevant documents for FAQs and troubleshooting.


## Tech Stack
- **Frontend**: React (web interaction)
- **Backend**: Go (Gin framework)
- **Database**: PostgreSQL
- **APIs**: OpenAI API, Dialogflow, META APIs, Telegram API, LINE API
<!---  **Cloud**: AWS (for deployment) -->



## Installation

- **Clone the Repository**:
   ```bash
   git clone https://github.com/petersun1937/CrossPlatform-TechSupport-Chatbot.git
   cd CrossPlatform-TechSupport-Chatbot
   ```

- **Backend Setup**:
   - **Ensure Go and Python are installed**:
      - Make sure that Go (version 1.16 or higher) and Python (version 3.6 or higher) are installed on your system.
      - You can verify the installation by running the following commands:
        ```bash
        go version
        python --version
        ```

  - **Install Python packages**:
      - The backend requires several Python packages for processing PDFs. To install them, run:
        ```bash
        pip install pdfplumber pytesseract pdf2image PyPDF2
        ```
      - These packages handle PDF text extraction and Optical Character Recognition (OCR) for images within PDFs.
<!--- 
   - **Tesseract OCR Installation** (Optional for OCR capabilities):
      - **Linux**: Install Tesseract via the package manager:
        ```bash
        sudo apt-get install tesseract-ocr
        ```
      - **Windows**: Download and install [Tesseract](https://github.com/tesseract-ocr/tesseract/wiki).
      - Ensure Tesseract is accessible through your system's PATH.
-->
   - **Set up environment variables**:
      - Create a `.env` file in the `configs/` directory to store environment variables such as API keys and database configurations (refer to `sample.env`).
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




## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.