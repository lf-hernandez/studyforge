# StudyForge Project Outline

## Executive Summary
StudyForge is a web-based application that transforms college textbook PDFs into personalized study materials using AI. Students can upload their textbooks, extract specific page ranges, and generate customized content including summaries, quizzes, study guides, and flashcards tailored to their academic level.

## Core Features

### 1. PDF Processing
- **Upload**: Accept PDF files up to 50MB
- **Text Extraction**: Extract text from native PDFs
- **OCR Support**: Process scanned PDFs using Tesseract
- **Page Selection**: Extract specific page ranges (e.g., pages 1-10)
- **Preview**: Display PDF pages in browser using PDF.js

### 2. AI-Powered Content Generation
- **Summaries**: Concise summaries at different academic levels
- **Quizzes**: Multiple choice and short answer questions
- **Study Notes**: Key points and concepts
- **Study Guides**: Comprehensive guides for exam preparation
- **Flashcards**: Question-answer pairs for memorization
- **Custom Prompts**: Allow users to request specific content types

### 3. User Experience
- **Session Management**: No login required, session-based usage
- **Progress Tracking**: Show processing status for long operations
- **Export Options**: Download generated content as PDF or text
- **History**: Access previously generated materials during session
- **Responsive Design**: Works on desktop and mobile devices

## Technical Architecture

### System Architecture
```
┌─────────────────────────────────────────────────────┐
│                   Web Browser                       │
│  ┌──────────────────────────────────────────────┐  │
│  │         Frontend (HTML/CSS/JS)               │  │
│  │  - PDF Upload Interface                      │  │
│  │  - PDF Viewer (PDF.js)                       │  │
│  │  - Study Material Generator                  │  │
│  │  - Content Display & Export                  │  │
│  └──────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────┘
                           │
                    HTTP REST API
                           │
┌─────────────────────────────────────────────────────┐
│                  Go Backend Server                  │
│  ┌──────────────────────────────────────────────┐  │
│  │            API Layer (net/http)              │  │
│  │  - Request routing                           │  │
│  │  - Input validation                          │  │
│  │  - Response formatting                       │  │
│  └──────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────┐  │
│  │             Service Layer                    │  │
│  │  - PDF processing service                    │  │
│  │  - OCR service (Tesseract)                   │  │
│  │  - AI integration service                    │  │
│  │  - Content generation service                │  │
│  └──────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────┐  │
│  │           Data Access Layer                  │  │
│  │  - SQLite database operations                │  │
│  │  - File system operations                    │  │
│  │  - Cache management                          │  │
│  └──────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────┘
                           │
        ┌──────────────────────────────────┐
        │                                    │
┌──────────────┐                    ┌──────────────┐
│    SQLite    │                    │  Hugging Face│
│   Database   │                    │   AI API     │
└──────────────┘                    └──────────────┘
```

### Technology Stack

#### Backend (Go)
- **HTTP Server**: `net/http` standard library
- **Database**: SQLite with `database/sql`
- **PDF Processing**: pdfcpu library
- **OCR**: Tesseract via gosseract
- **JSON**: `encoding/json` standard library
- **File Handling**: `os`, `io`, `path/filepath` packages
- **Concurrency**: Goroutines and channels

#### Frontend (Vanilla)
- **HTML5**: Semantic markup
- **CSS3**: Custom styling, flexbox/grid layouts
- **JavaScript ES6+**: No frameworks
- **PDF.js**: For PDF viewing
- **Fetch API**: For AJAX requests
- **LocalStorage**: For client-side caching

#### External Services
- **Hugging Face API**: Free tier for AI models
  - Text summarization: facebook/bart-large-cnn
  - Question generation: iarfmoose/t5-base-question-generator
  - General tasks: google/flan-t5-base

## Database Schema

### Tables Structure

```sql
-- Sessions table for user tracking
CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_accessed TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ip_address TEXT,
    user_agent TEXT,
    is_active BOOLEAN DEFAULT TRUE
);

-- Documents table for uploaded PDFs
CREATE TABLE documents (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT NOT NULL,
    original_filename TEXT NOT NULL,
    stored_filename TEXT NOT NULL,
    file_path TEXT NOT NULL,
    file_size INTEGER NOT NULL,
    page_count INTEGER,
    upload_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_accessed TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE,
    metadata JSON,
    FOREIGN KEY (session_id) REFERENCES sessions(id)
);

-- Extracted content cache
CREATE TABLE extracted_content (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    document_id INTEGER NOT NULL,
    page_start INTEGER NOT NULL,
    page_end INTEGER NOT NULL,
    content TEXT NOT NULL,
    is_ocr BOOLEAN DEFAULT FALSE,
    extraction_time INTEGER, -- milliseconds
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (document_id) REFERENCES documents(id),
    UNIQUE(document_id, page_start, page_end)
);

-- Generated study materials
CREATE TABLE generated_content (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT NOT NULL,
    document_id INTEGER NOT NULL,
    content_type TEXT NOT NULL, -- 'summary', 'quiz', 'notes', 'guide', 'flashcards'
    academic_level TEXT, -- 'high_school', 'undergraduate', 'graduate'
    input_pages TEXT, -- e.g., "1-10"
    input_text TEXT,
    output_content JSON NOT NULL,
    ai_model TEXT NOT NULL,
    generation_time INTEGER, -- milliseconds
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES sessions(id),
    FOREIGN KEY (document_id) REFERENCES documents(id)
);

-- API usage tracking
CREATE TABLE api_usage (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT NOT NULL,
    api_name TEXT NOT NULL, -- 'huggingface', 'tesseract'
    endpoint TEXT,
    tokens_used INTEGER,
    response_time INTEGER, -- milliseconds
    status_code INTEGER,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES sessions(id)
);

-- Create indexes for performance
CREATE INDEX idx_sessions_active ON sessions(is_active);
CREATE INDEX idx_documents_session ON documents(session_id);
CREATE INDEX idx_extracted_content_document ON extracted_content(document_id);
CREATE INDEX idx_generated_content_session ON generated_content(session_id);
CREATE INDEX idx_generated_content_document ON generated_content(document_id);
CREATE INDEX idx_api_usage_session ON api_usage(session_id);
```

## API Endpoints Specification

### Document Management

#### POST /api/documents/upload
Upload a PDF document
```json
// Request (multipart/form-data)
{
    "file": <PDF file>
}

// Response
{
    "success": true,
    "data": {
        "document_id": "123",
        "filename": "textbook.pdf",
        "page_count": 450,
        "file_size": 25000000
    }
}
```

#### GET /api/documents/:id
Get document information
```json
// Response
{
    "success": true,
    "data": {
        "id": "123",
        "filename": "textbook.pdf",
        "page_count": 450,
        "upload_date": "2024-01-01T12:00:00Z",
        "file_size": 25000000
    }
}
```

#### DELETE /api/documents/:id
Delete a document
```json
// Response
{
    "success": true,
    "message": "Document deleted successfully"
}
```

### Content Extraction

#### POST /api/documents/:id/extract
Extract text from specified pages
```json
// Request
{
    "page_start": 1,
    "page_end": 10,
    "use_ocr": false
}

// Response
{
    "success": true,
    "data": {
        "content": "Extracted text content...",
        "pages_processed": 10,
        "processing_time": 1500
    }
}
```

### Study Material Generation

#### POST /api/study/generate
Generate study materials
```json
// Request
{
    "document_id": "123",
    "page_start": 1,
    "page_end": 10,
    "material_type": "summary", // 'summary', 'quiz', 'notes', 'guide', 'flashcards'
    "academic_level": "undergraduate", // 'high_school', 'undergraduate', 'graduate'
    "additional_instructions": "Focus on key concepts"
}

// Response
{
    "success": true,
    "data": {
        "content_id": "456",
        "material_type": "summary",
        "content": {
            "title": "Chapter 1 Summary",
            "text": "Generated summary content...",
            "key_points": ["Point 1", "Point 2"]
        },
        "model_used": "facebook/bart-large-cnn",
        "generation_time": 3000
    }
}
```

#### GET /api/study/:id
Retrieve previously generated content
```json
// Response
{
    "success": true,
    "data": {
        "content_id": "456",
        "material_type": "summary",
        "content": {
            "title": "Chapter 1 Summary",
            "text": "Generated summary content..."
        },
        "created_at": "2024-01-01T12:30:00Z"
    }
}
```

### Session Management

#### GET /api/session
Get current session information
```json
// Response
{
    "success": true,
    "data": {
        "session_id": "uuid-string",
        "created_at": "2024-01-01T12:00:00Z",
        "documents_count": 3,
        "generated_content_count": 10
    }
}
```

### System

#### GET /api/health
Health check endpoint
```json
// Response
{
    "status": "healthy",
    "version": "1.0.0",
    "uptime": 3600
}
```

## Development Phases

### Phase 1: Foundation (Week 1-2)
- [x] Project structure setup
- [x] Documentation creation
- [ ] Basic Go server with net/http
- [ ] SQLite database setup and migrations
- [ ] Basic HTML/CSS frontend structure
- [ ] File upload endpoint
- [ ] Session management

### Phase 2: PDF Processing (Week 3-4)
- [ ] Integrate pdfcpu library
- [ ] Implement text extraction
- [ ] Add page range selection
- [ ] Integrate Tesseract OCR
- [ ] Create extraction caching
- [ ] Error handling for various PDF types

### Phase 3: AI Integration (Week 5-6)
- [ ] Hugging Face API client
- [ ] Implement summary generation
- [ ] Implement quiz generation
- [ ] Implement study notes generation
- [ ] Add prompt templates
- [ ] Response caching system

### Phase 4: Frontend Development (Week 7-8)
- [ ] PDF upload interface
- [ ] PDF.js viewer integration
- [ ] Study material generation UI
- [ ] Progress indicators
- [ ] Content display components
- [ ] Export functionality

### Phase 5: Enhancement (Week 9-10)
- [ ] Performance optimization
- [ ] Advanced error handling
- [ ] Rate limiting
- [ ] Session cleanup jobs
- [ ] Analytics dashboard
- [ ] Mobile responsiveness

### Phase 6: Testing & Deployment (Week 11-12)
- [ ] Unit tests for all services
- [ ] Integration tests
- [ ] Load testing
- [ ] Security audit
- [ ] Docker containerization
- [ ] Deployment documentation
- [ ] User documentation

## AI Model Integration

### Recommended Models

#### For Summarization
1. **facebook/bart-large-cnn**
   - Best for news-style summaries
   - Good performance on academic text
   - Free tier: 1000 requests/month

2. **google/flan-t5-base**
   - Versatile, handles multiple tasks
   - Good for instructional summaries
   - Free tier: 1000 requests/month

#### For Question Generation
1. **iarfmoose/t5-base-question-generator**
   - Specialized for educational questions
   - Generates relevant questions from text
   - Free tier: 1000 requests/month

2. **valhalla/t5-base-qa-qg-hl**
   - Generates both questions and answers
   - Good for creating study quizzes
   - Free tier: 1000 requests/month

### Prompt Templates

#### Summary Generation
```
Summarize the following academic text for a {academic_level} student.
Focus on key concepts, important definitions, and main ideas.
Keep the summary concise but comprehensive.

Text: {extracted_text}
```

#### Quiz Generation
```
Create a {quiz_type} quiz with {num_questions} questions based on the following academic text.
Include a mix of conceptual and factual questions appropriate for {academic_level} level.
Format: Question, followed by answer.

Text: {extracted_text}
```

#### Study Notes Generation
```
Create detailed study notes from the following academic text for a {academic_level} student.
Organize the notes with:
1. Main topics
2. Key definitions
3. Important formulas or concepts
4. Examples
5. Summary points

Text: {extracted_text}
```

## Performance Targets

### Response Time Goals
- File upload: < 2 seconds for 10MB
- Text extraction (native PDF): < 1 second per page
- OCR processing: < 3 seconds per page
- AI content generation: < 5 seconds
- Page load: < 1 second
- API responses: < 200ms (non-AI)

### Capacity Planning
- Concurrent users: 100
- Max file size: 50MB
- Storage per session: 200MB
- Session timeout: 24 hours
- Database size: 10GB initial

### Optimization Strategies
1. **Caching**
   - Cache extracted text
   - Cache AI responses
   - Browser caching for static assets

2. **Background Processing**
   - Queue OCR tasks
   - Async AI API calls
   - Batch database operations

3. **Resource Management**
   - Connection pooling
   - Goroutine pools
   - Memory-mapped file reading

## Security Considerations

### Input Validation
- File type verification (PDF only)
- File size limits
- Page range validation
- Text length limits for AI
- SQL injection prevention
- XSS prevention

### Data Protection
- Session token security
- Secure file storage
- HTTPS enforcement
- Rate limiting
- CORS configuration
- Content Security Policy

### Privacy
- Auto-delete files after 24 hours
- No personal data collection
- Clear privacy policy
- GDPR compliance ready
- Data export capability

## Error Handling Strategy

### Error Categories
1. **User Errors** (4xx)
   - Invalid file format
   - File too large
   - Invalid page range
   - Rate limit exceeded

2. **System Errors** (5xx)
   - Database connection failure
   - AI API unavailable
   - OCR failure
   - Server overload

### Error Response Format
```json
{
    "success": false,
    "error": {
        "code": "FILE_TOO_LARGE",
        "message": "File size exceeds 50MB limit",
        "details": {
            "max_size": 52428800,
            "file_size": 75000000
        }
    }
}
```

## Monitoring and Analytics

### Key Metrics
- API response times
- File processing success rate
- AI API usage and costs
- Error rates by type
- User session duration
- Popular content types
- Storage usage

### Logging
- Structured JSON logging
- Log levels: DEBUG, INFO, WARN, ERROR
- Centralized log aggregation
- Error alerting thresholds

## Future Enhancements

### Version 2.0 Features
- User accounts and history
- Collaborative study groups
- Multiple language support
- Custom AI model fine-tuning
- Integration with learning management systems
- Mobile applications
- Batch processing for multiple chapters
- Advanced formatting preservation

### Potential Integrations
- Google Drive / Dropbox
- Citation management tools
- Note-taking applications
- Calendar integration for study schedules
- Export to Anki/Quizlet

## Success Metrics

### Technical KPIs
- 99.9% uptime
- < 5% error rate
- < 5 second average processing time
- > 90% cache hit rate

### User KPIs
- Average session duration > 15 minutes
- Return user rate > 40%
- Content generation success rate > 95%
- User satisfaction score > 4/5

## Risk Assessment

### Technical Risks
1. **AI API Rate Limits**
   - Mitigation: Implement caching and queuing

2. **Large File Processing**
   - Mitigation: Streaming and chunked processing

3. **OCR Accuracy**
   - Mitigation: Multiple OCR engines, user feedback

### Business Risks
1. **AI API Cost Overruns**
   - Mitigation: Usage monitoring and limits

2. **Copyright Concerns**
   - Mitigation: Clear terms of service

3. **Data Privacy Issues**
   - Mitigation: Strong privacy policy and auto-deletion

## Resources and References

### Documentation
- [Go net/http documentation](https://pkg.go.dev/net/http)
- [SQLite documentation](https://www.sqlite.org/docs.html)
- [PDF.js documentation](https://mozilla.github.io/pdf.js/)
- [Hugging Face API docs](https://huggingface.co/docs/api-inference/index)

### Libraries
- [pdfcpu](https://github.com/pdfcpu/pdfcpu)
- [gosseract](https://github.com/otiai10/gosseract)
- [PDF.js](https://github.com/mozilla/pdf.js)

### Tools
- [Tesseract OCR](https://github.com/tesseract-ocr/tesseract)
- [SQLite Browser](https://sqlitebrowser.org/)
- [Postman](https://www.postman.com/) for API testing

---

Last Updated: [Current Date]
Version: 1.0.0
Status: Planning Phase Complete