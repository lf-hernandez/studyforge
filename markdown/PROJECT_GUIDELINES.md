# StudyForge Project Guidelines

## Project Overview
StudyForge is a college textbook PDF study assistant that helps students generate study materials from their textbooks using AI. The application processes PDF files, extracts text (with OCR support), and generates summaries, quizzes, study guides, and other educational content.

## Development Philosophy
- **No Frameworks Policy**: Use standard libraries wherever possible
- **Simplicity First**: Prioritize simple, maintainable solutions
- **Security by Design**: Consider security implications in every feature
- **Progressive Enhancement**: Start with basic functionality, enhance gradually

## Technology Stack

### Backend
- **Language**: Go (Golang)
- **HTTP Server**: Standard library `net/http` (NO frameworks like Gin, Echo, etc.)
- **Database**: SQLite with `database/sql` standard library (NO ORMs)
- **PDF Processing**: pdfcpu library (minimal external dependency)
- **OCR**: Tesseract via gosseract wrapper
- **AI Integration**: Hugging Face API (free tier)

### Frontend
- **HTML**: Semantic HTML5
- **CSS**: Pure CSS3 (NO Bootstrap, Tailwind, etc.)
- **JavaScript**: Vanilla ES6+ (NO React, Vue, Angular, jQuery, etc.)
- **PDF Viewing**: PDF.js library (necessary for in-browser PDF display)

## Coding Standards

### Go Code Standards

#### File Organization
```go
// Each file should have a clear, single responsibility
// File naming: use snake_case (e.g., pdf_handler.go, user_service.go)
```

#### Package Structure
- `cmd/` - Application entry points
- `internal/` - Private application code
- `pkg/` - Public libraries that could be imported by other projects

#### Error Handling
```go
// Always check and handle errors explicitly
result, err := someFunction()
if err != nil {
    // Log the error with context
    log.Printf("failed to process PDF: %v", err)
    // Return wrapped error with context
    return fmt.Errorf("processing PDF %s: %w", filename, err)
}
```

#### Database Queries
```go
// Use prepared statements to prevent SQL injection
stmt, err := db.Prepare("SELECT * FROM documents WHERE id = ?")
if err != nil {
    return nil, err
}
defer stmt.Close()

// Use transactions for multiple related operations
tx, err := db.Begin()
if err != nil {
    return err
}
defer tx.Rollback() // Will be no-op if tx.Commit() is called
```

#### HTTP Handlers
```go
// Standard handler signature
func handlePDFUpload(w http.ResponseWriter, r *http.Request) {
    // Set appropriate headers
    w.Header().Set("Content-Type", "application/json")

    // Validate method
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Process request...
}
```

### JavaScript Standards

#### File Organization
```javascript
// Use IIFE to avoid global scope pollution
(function() {
    'use strict';

    // Module code here
})();

// Or use ES6 modules where appropriate
export function processDocument() {
    // Implementation
}
```

#### API Communication
```javascript
// Use async/await for cleaner async code
async function uploadPDF(file) {
    try {
        const formData = new FormData();
        formData.append('file', file);

        const response = await fetch('/api/upload', {
            method: 'POST',
            body: formData
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        return await response.json();
    } catch (error) {
        console.error('Upload failed:', error);
        throw error;
    }
}
```

#### DOM Manipulation
```javascript
// Cache DOM queries
const uploadButton = document.getElementById('upload-btn');
const progressBar = document.querySelector('.progress-bar');

// Use event delegation for dynamic content
document.addEventListener('click', function(event) {
    if (event.target.classList.contains('delete-btn')) {
        handleDelete(event.target.dataset.id);
    }
});
```

### CSS Standards

#### Organization
```css
/* Use logical grouping and comments */

/* === Base Styles === */
* {
    box-sizing: border-box;
}

/* === Layout === */
.container {
    max-width: 1200px;
    margin: 0 auto;
}

/* === Components === */
.button {
    /* Component styles */
}

/* === Utilities === */
.text-center {
    text-align: center;
}
```

#### Naming Convention
```css
/* Use BEM-like naming for components */
.card { }
.card__header { }
.card__body { }
.card--featured { }

/* Use semantic class names */
.pdf-viewer { }
.study-generator { }
.quiz-section { }
```

## Security Guidelines

### File Upload Security
1. **Validate file type**: Only accept PDF files
2. **Check file size**: Implement reasonable limits (e.g., 50MB)
3. **Scan for malware**: Use ClamAV or similar if possible
4. **Store outside web root**: Never serve uploaded files directly
5. **Use random filenames**: Prevent directory traversal attacks

### Input Validation
```go
// Validate all user inputs
func validatePageRange(start, end int, maxPages int) error {
    if start < 1 || end < 1 {
        return errors.New("page numbers must be positive")
    }
    if start > end {
        return errors.New("start page must be less than end page")
    }
    if end > maxPages {
        return fmt.Errorf("end page %d exceeds document pages %d", end, maxPages)
    }
    return nil
}
```

### Session Management
- Use secure, random session tokens (UUID v4)
- Implement session expiration (24-48 hours)
- Use secure cookies with HttpOnly and SameSite flags
- Implement rate limiting per session

### SQL Injection Prevention
```go
// NEVER concatenate user input into SQL queries
// BAD:
query := fmt.Sprintf("SELECT * FROM users WHERE id = %s", userID)

// GOOD:
query := "SELECT * FROM users WHERE id = ?"
rows, err := db.Query(query, userID)
```

## API Design Patterns

### RESTful Endpoints
```
POST   /api/documents/upload     - Upload PDF
GET    /api/documents/{id}        - Get document info
DELETE /api/documents/{id}        - Delete document
POST   /api/documents/{id}/extract - Extract text from pages
POST   /api/study/generate        - Generate study materials
GET    /api/study/{id}           - Retrieve generated content
GET    /api/health                - Health check endpoint
```

### Request/Response Format
```json
// Request
{
    "document_id": "uuid",
    "page_start": 1,
    "page_end": 10,
    "material_type": "summary",
    "academic_level": "undergraduate"
}

// Success Response
{
    "success": true,
    "data": {
        "content": "Generated content...",
        "metadata": {
            "generated_at": "2024-01-01T00:00:00Z",
            "model_used": "facebook/bart-large-cnn"
        }
    }
}

// Error Response
{
    "success": false,
    "error": {
        "code": "INVALID_INPUT",
        "message": "Page range exceeds document length"
    }
}
```

## Testing Requirements

### Unit Tests
- Test all service functions
- Test input validation
- Test error handling
- Aim for >70% code coverage

### Integration Tests
- Test complete API endpoints
- Test database operations
- Test file upload/processing flow

### Test File Naming
```
pdf_service.go      → pdf_service_test.go
api_handlers.go     → api_handlers_test.go
```

## Documentation Standards

### Code Comments
```go
// Package pdf provides utilities for PDF processing and text extraction
package pdf

// ExtractText extracts text content from specified page range.
// Returns extracted text and any OCR errors encountered.
// Pages are 1-indexed (first page is 1, not 0).
func ExtractText(filepath string, startPage, endPage int) (string, error) {
    // Implementation
}
```

### README Files
Each major component should have a README explaining:
- Purpose and functionality
- Setup instructions
- API documentation
- Examples of usage

## Version Control

### Git Workflow
1. **Branch naming**: `feature/description`, `bugfix/description`, `hotfix/description`
2. **Commit messages**: Use conventional commits
   - `feat: add PDF upload functionality`
   - `fix: correct page extraction logic`
   - `docs: update API documentation`
   - `refactor: simplify session management`
   - `test: add unit tests for OCR service`

### Code Review Checklist
- [ ] Code follows project standards
- [ ] Security considerations addressed
- [ ] Tests written and passing
- [ ] Documentation updated
- [ ] No sensitive data in code
- [ ] Error handling implemented

## Development Environment

### Required Tools
- Go 1.20+
- SQLite 3
- Tesseract OCR
- Node.js (for PDF.js only)
- Git

### Environment Variables
```bash
# .env file (never commit this)
SERVER_PORT=8080
DATABASE_PATH=./data/studyforge.db
UPLOAD_DIR=./uploads
MAX_FILE_SIZE=52428800  # 50MB in bytes
SESSION_TIMEOUT=86400   # 24 hours in seconds
HUGGINGFACE_API_KEY=your_api_key_here
```

## Performance Guidelines

### Database Optimization
- Index frequently queried columns
- Use connection pooling
- Implement query result caching
- Regular VACUUM on SQLite

### File Processing
- Stream large files instead of loading into memory
- Process PDFs in chunks
- Cache extracted text in database
- Implement background job processing for OCR

### API Response Times
- Target: <200ms for simple queries
- Target: <2s for PDF text extraction
- Target: <5s for AI content generation
- Implement progress indicators for long operations

## Monitoring and Logging

### Logging Levels
```go
log.Debug("Detailed debug information")
log.Info("General information")
log.Warn("Warning messages")
log.Error("Error messages")
log.Fatal("Critical errors that stop execution")
```

### Metrics to Track
- API response times
- File upload sizes and processing times
- AI API usage and costs
- Error rates by endpoint
- Active sessions count

## Deployment Considerations

### Production Checklist
- [ ] Environment variables configured
- [ ] Database migrations run
- [ ] File permissions set correctly
- [ ] SSL/TLS configured
- [ ] Rate limiting enabled
- [ ] Error logging configured
- [ ] Backup strategy implemented
- [ ] Monitoring alerts set up

### Scaling Considerations
- Horizontal scaling via load balancer
- Database replication for read scaling
- CDN for static assets
- Queue system for heavy processing
- Caching layer (Redis) if needed

## License and Legal

### Open Source Dependencies
- Document all third-party libraries
- Ensure license compatibility
- Attribute properly in documentation

### User Privacy
- Implement data retention policies
- Clear privacy policy
- GDPR compliance considerations
- User data deletion capability

## Contact and Support

### Team Responsibilities
- Backend Development: [Owner]
- Frontend Development: [Owner]
- AI Integration: [Owner]
- Testing & QA: [Owner]

### Getting Help
- Check documentation first
- Search existing issues
- Create detailed bug reports
- Include error messages and logs

---

Last Updated: [Current Date]
Version: 1.0.0