import PyPDF2
import pdfplumber
from pdf2image import convert_from_path
import pytesseract
import logging
import sys

sys.stdout.reconfigure(encoding='utf-8')

# Set up logging configuration
logging.basicConfig(
    filename='pdf_extraction.log',  # Log file name
    level=logging.INFO,             # Log all info level and above
    format='%(asctime)s - %(levelname)s - %(message)s',  # Log format
)

def decrypt_pdf(file_path):
    try:
        with open(file_path, 'rb') as file:
            reader = PyPDF2.PdfReader(file)
            if reader.is_encrypted:
                logging.info(f"PDF is encrypted, attempting to decrypt: {file_path}")
                reader.decrypt('')
            decrypted_file_path = f"{file_path}_decrypted.pdf"
            with open(decrypted_file_path, 'wb') as decrypted_file:
                writer = PyPDF2.PdfWriter()
                for page in range(len(reader.pages)):
                    writer.add_page(reader.pages[page])
                writer.write(decrypted_file)
            logging.info(f"Decryption completed, saved to {decrypted_file_path}")
            return decrypted_file_path
    except Exception as e:
        logging.error(f"Error decrypting PDF: {e}", exc_info=True)
        raise e
    return file_path

def extract_text_from_pdf(file_path):
    extracted_text = ''
    logging.info(f"Starting text extraction for file: {file_path}")
    
    with pdfplumber.open(file_path) as pdf:
        logging.info(f"PDF opened successfully, total pages: {len(pdf.pages)}")
        
        for i, page in enumerate(pdf.pages):
            try:
                text = page.extract_text()
                if text:
                    extracted_text += text + '\n'
                    logging.info(f"Extracted text from page {i + 1}")
                else:
                    logging.warning(f"No text found on page {i + 1}, processing with OCR...")
                    # If no text is found, run OCR on the page
                    images = convert_from_path(file_path, first_page=i+1, last_page=i+1)
                    for image in images:
                        ocr_text = pytesseract.image_to_string(image)
                        extracted_text += ocr_text + '\n'
                    logging.info(f"OCR text extraction completed for page {i + 1}")
            except Exception as e:
                logging.error(f"Error processing page {i + 1}: {e}", exc_info=True)
                print(f"Error processing page {i + 1}: {e}", file=sys.stderr)
    
    logging.info(f"Text extraction completed for file: {file_path}")
    return extracted_text

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python pdf_extractor.py <path_to_pdf>")
        logging.error("No file path provided. Exiting script.")
        sys.exit(1)
    
    file_path = sys.argv[1]  # Get the file path from the command-line arguments
    
    try:
        logging.info(f"Script started with file: {file_path}")
        
        # Attempt to decrypt the PDF if necessary
        decrypted_file_path = decrypt_pdf(file_path)
        
        text = extract_text_from_pdf(decrypted_file_path)
        print(text)  # Print extracted text to stdout
        logging.info("Text output successfully printed.")
    except Exception as e:
        logging.error(f"Error extracting text: {e}", exc_info=True)
        print(f"Error extracting text: {e}")
        sys.exit(1)




'''def extract_text_from_pdf(file_path):
    extracted_text = ''
    with pdfplumber.open(file_path) as pdf:
        for page in pdf.pages:
            text = page.extract_text()
            if text:
                extracted_text += text + '\n'
    return extracted_text

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python pdf_extractor.py <path_to_pdf>")
        sys.exit(1)
    
    file_path = sys.argv[1]  # Get the file path from the command-line arguments
    try:
        text = extract_text_from_pdf(file_path)
        print(text)  # Print extracted text to stdout
    except Exception as e:
        print(f"Error extracting text: {e}")
        sys.exit(1)'''
