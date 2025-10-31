# StudyForge

An AI-powered PDF textbook study assistant that extracts text from PDF documents and generates educational summaries, study materials, and interactive Q&A using Hugging Face models.

## Features

- PDF text extraction with page range selection
- AI-powered content summarization tailored to academic levels
- Text cleaning to handle PDF extraction artifacts
- Session-based document management
- RESTful API with standard library Go backend
- Vanilla JavaScript frontend (no frameworks)
- SQLite database for data persistence

## Prerequisites

- Go 1.20 or higher
- Hugging Face API key (free tier available)
- Web browser with JavaScript enabled

## Installation

### 1. Clone the repository

```bash
git clone git@github.com:lf-hernandez/studyforge.git
cd studyforge
```

### 2. Install Go dependencies

```bash
go mod download
```

### 3. Set up environment variables

Create a `.env` file in the project root:

```bash
# Server Configuration
SERVER_HOST=localhost
SERVER_PORT=8080

# Database
DATABASE_PATH=./data/studyforge.db

# File Storage
UPLOAD_DIR=./uploads
MAX_FILE_SIZE=52428800

# Session
SESSION_TIMEOUT=3600

# Hugging Face API
HUGGINGFACE_API_KEY=your_api_key_here
HUGGINGFACE_API_URL=https://api-inference.huggingface.co/models

# Logging
LOG_LEVEL=info
```

Replace `your_api_key_here` with your actual Hugging Face API key.

### 4. Create required directories

```bash
mkdir -p data uploads
```

## Running the Application

### Start the server

```bash
# Export environment variables
export HUGGINGFACE_API_KEY="your_api_key_here"
export HUGGINGFACE_API_URL="https://api-inference.huggingface.co/models"

# Run the server
go run cmd/server/main.go
```

The server will start on http://localhost:8080

### Access the web interface

Open your browser and navigate to:
```
http://localhost:8080
```

## Usage

1. **Upload a PDF**: Click "Upload PDF" and select your textbook file (max 50MB)
2. **Select Pages**: Enter the start and end page numbers for extraction
3. **Choose Academic Level**: Select from High School, Undergraduate, or Graduate
4. **Generate Summary**: Click "Generate Summary" to create an AI-powered summary
5. **View Results**: The summary will appear in sections, organized by the extracted content

## API Endpoints

### Health Check
```
GET /api/health
```

### PDF Management
```
POST /api/pdf/upload          - Upload a PDF file
GET  /api/pdf/documents        - List uploaded documents
GET  /api/pdf/extract          - Extract text from pages
```

### Study Material Generation
```
POST /api/study/generate       - Generate summary from pages
GET  /api/study/content/:id    - Retrieve generated content
```

## Project Structure

```
studyforge/
├── cmd/
│   └── server/
│       └── main.go            # Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers/          # HTTP request handlers
│   │   ├── middleware/        # HTTP middleware
│   │   └── router.go          # Route definitions
│   ├── config/
│   │   └── config.go          # Configuration management
│   ├── models/                # Data models
│   ├── repository/            # Database operations
│   └── services/              # Business logic
├── pkg/
│   ├── ai/
│   │   └── huggingface.go    # AI integration
│   ├── pdf/
│   │   └── extractor.go      # PDF text extraction
│   └── utils/                 # Utility functions
├── web/
│   ├── index.html             # Frontend interface
│   ├── css/
│   │   └── main.css          # Styles
│   └── js/
│       ├── api.js            # API client
│       └── app.js            # Application logic
├── migrations/                # Database migrations
├── markdown/                  # Documentation
└── go.mod                    # Go dependencies
```

## Configuration

### Database

SQLite database is automatically created at `./data/studyforge.db` on first run.

### Session Management

Sessions are UUID-based and stored in the database. Default timeout is 1 hour.

### PDF Processing

- Maximum file size: 50MB (configurable)
- Supported format: PDF
- Text extraction: Uses ledongthuc/pdf library
- Text cleaning: Removes PDF artifacts and formatting issues

### AI Integration

- Model: facebook/bart-large-cnn
- Chunking: Automatic text chunking for large documents
- Academic levels: High School, Undergraduate, Graduate
- Response format: Sectioned summaries for multi-chunk content

## Development

### Running tests

```bash
go test ./...
```

### Building for production

```bash
go build -o studyforge cmd/server/main.go
./studyforge
```

### Database migrations

Migrations run automatically on server start. SQL files are in `migrations/` directory.

## Troubleshooting

### Common Issues

1. **"API key not found"**: Ensure HUGGINGFACE_API_KEY is exported in your environment
2. **"Cannot open database"**: Check write permissions for `./data` directory
3. **"PDF extraction failed"**: Ensure PDF is not corrupted or password-protected
4. **"Summary too short"**: The AI chunks large texts; this is normal behavior

### Logs

Server logs are output to stdout. Set LOG_LEVEL=debug for verbose logging.

## Security Considerations

- Sessions are temporary and UUID-based
- File uploads are restricted by size and type
- Database uses prepared statements to prevent SQL injection
- No authentication system (designed for local/personal use)

## Future Enhancements

See `markdown/features/CUSTOM_AI_PROMPTS.md` for planned features including:
- Custom AI prompts and Q&A
- Quiz and flashcard generation
- PDF export functionality
- Multiple study material formats

## Technology Stack

- **Backend**: Go (standard library only, no frameworks)
- **Frontend**: HTML, CSS, Vanilla JavaScript (no frameworks)
- **Database**: SQLite with standard library driver
- **AI**: Hugging Face Inference API
- **PDF**: ledongthuc/pdf library

## License

MIT

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## Support

For issues, questions, or suggestions, please open an issue on GitHub.
