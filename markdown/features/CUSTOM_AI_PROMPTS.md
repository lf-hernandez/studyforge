# Custom AI Prompts & Interactive Study Tools Feature

## Overview
Enable users to create custom AI prompts for interactive Q&A, study material generation, and content analysis from uploaded PDFs. This feature transforms StudyForge from a simple summarizer into an interactive AI-powered study assistant.

## Core Components

### 1. Custom Prompt Interface
- **Text Area**: Large input field for custom prompts/questions
- **Academic Level Dropdown**: Select target academic level
  - High School
  - Undergraduate
  - Graduate
  - Professional
  - Custom/General
- **Query Type Selection**: Pre-built prompt templates
  - Question Answering
  - Study Notes Generation
  - Quiz Creation
  - Concept Explanation
  - Timeline Creation
  - Custom Query

### 2. Interactive Q&A System
Users can ask specific questions about the content:
- "Where did Columbus go in 1492?"
- "What were the main causes of the Protestant Reformation?"
- "Explain the Treaty of Tordesillas and its impact"
- "Who were the key figures in Spanish exploration?"
- "What economic factors drove European exploration?"

### 3. Study Material Generator
Generate different types of study materials:
- **Study Notes**: Organized, bullet-point summaries
- **Quizzes**: Multiple choice, true/false, short answer
- **Flashcards**: Key terms and definitions
- **Timelines**: Chronological event sequences
- **Concept Maps**: Visual relationship diagrams
- **Essay Questions**: Open-ended analytical prompts

### 4. PDF Export Feature
- Download generated content as PDF
- Include original question/prompt
- Formatted for printing and offline study
- Optional watermark/branding

## Technical Implementation

### Backend Changes

#### 1. New API Endpoint: `/api/study/custom-prompt`
```go
type CustomPromptRequest struct {
    SessionID      string   `json:"session_id"`
    DocumentID     int      `json:"document_id"`
    PageStart      int      `json:"page_start"`
    PageEnd        int      `json:"page_end"`
    Prompt         string   `json:"prompt"`
    PromptType     string   `json:"prompt_type"`
    AcademicLevel  string   `json:"academic_level"`
    OutputFormat   string   `json:"output_format"`
}

type CustomPromptResponse struct {
    ContentID      int      `json:"content_id"`
    Response       string   `json:"response"`
    PromptType     string   `json:"prompt_type"`
    GenerationTime int      `json:"generation_time"`
    ModelUsed      string   `json:"model_used"`
    PDFUrl         string   `json:"pdf_url,omitempty"`
}
```

#### 2. Enhanced AI Service
```go
// pkg/ai/prompt_builder.go
type PromptBuilder struct {
    academicLevel string
    promptType    string
}

func (pb *PromptBuilder) BuildCustomPrompt(userPrompt, content string) string {
    // Build context-aware prompt based on type
}

func (pb *PromptBuilder) BuildQuizPrompt(content string, numQuestions int) string {
    // Generate quiz questions
}

func (pb *PromptBuilder) BuildStudyNotesPrompt(content string) string {
    // Create structured study notes
}
```

#### 3. PDF Generation Service
```go
// pkg/pdf/generator.go
type PDFGenerator struct {
    title    string
    content  string
    metadata map[string]string
}

func (g *PDFGenerator) GenerateStudyPDF() ([]byte, error) {
    // Create formatted PDF with content
}
```

### Frontend Changes

#### 1. New UI Component: Custom Prompt Panel
```html
<div class="custom-prompt-panel">
    <!-- Prompt Type Selector -->
    <div class="prompt-type-selector">
        <button class="type-btn active" data-type="question">Ask Question</button>
        <button class="type-btn" data-type="notes">Study Notes</button>
        <button class="type-btn" data-type="quiz">Generate Quiz</button>
        <button class="type-btn" data-type="flashcards">Flashcards</button>
        <button class="type-btn" data-type="custom">Custom</button>
    </div>

    <!-- Academic Level Dropdown -->
    <select id="academic-level-select" class="form-select">
        <option value="high_school">High School</option>
        <option value="undergraduate" selected>Undergraduate</option>
        <option value="graduate">Graduate</option>
        <option value="professional">Professional</option>
        <option value="general">General</option>
    </select>

    <!-- Custom Prompt Input -->
    <textarea id="custom-prompt-input"
              class="prompt-textarea"
              placeholder="Ask a question or describe what you want to generate...">
    </textarea>

    <!-- Template Suggestions -->
    <div class="template-suggestions">
        <span class="suggestion-chip">Where did...</span>
        <span class="suggestion-chip">Explain the significance of...</span>
        <span class="suggestion-chip">Compare and contrast...</span>
        <span class="suggestion-chip">What were the causes of...</span>
    </div>

    <!-- Action Buttons -->
    <div class="action-buttons">
        <button id="generate-btn" class="btn-primary">Generate Response</button>
        <button id="download-pdf-btn" class="btn-secondary" disabled>Download as PDF</button>
    </div>
</div>

<!-- Response Display Area -->
<div class="response-panel">
    <div class="response-header">
        <h3>AI Response</h3>
        <span class="response-type-badge">Question Answer</span>
        <button class="copy-btn">Copy</button>
    </div>
    <div id="response-content" class="response-content">
        <!-- AI-generated content appears here -->
    </div>
</div>
```

#### 2. JavaScript Functionality
```javascript
// web/js/custom_prompts.js
class CustomPromptManager {
    constructor() {
        this.currentPromptType = 'question';
        this.academicLevel = 'undergraduate';
        this.lastResponse = null;
    }

    async sendCustomPrompt(prompt) {
        const response = await api.post('/api/study/custom-prompt', {
            document_id: this.currentDocumentId,
            page_start: this.pageStart,
            page_end: this.pageEnd,
            prompt: prompt,
            prompt_type: this.currentPromptType,
            academic_level: this.academicLevel
        });

        this.lastResponse = response;
        this.displayResponse(response);
        return response;
    }

    displayResponse(response) {
        // Format and display based on response type
        switch(response.prompt_type) {
            case 'quiz':
                this.displayQuiz(response.response);
                break;
            case 'notes':
                this.displayStudyNotes(response.response);
                break;
            case 'flashcards':
                this.displayFlashcards(response.response);
                break;
            default:
                this.displayText(response.response);
        }
    }

    async downloadAsPDF() {
        if (!this.lastResponse) return;

        const pdfBlob = await api.downloadPDF(this.lastResponse.content_id);
        const url = URL.createObjectURL(pdfBlob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `study_material_${Date.now()}.pdf`;
        a.click();
    }
}
```

## Prompt Templates

### Question Answering Template
```
Based on the provided textbook content, answer the following question in detail:
[USER_QUESTION]

Provide specific examples, dates, and references from the text. Format your answer for a [ACADEMIC_LEVEL] student.
```

### Study Notes Template
```
Create comprehensive study notes from this textbook content for a [ACADEMIC_LEVEL] student. Include:
1. Main concepts and key terms
2. Important dates and events
3. Significant figures and their contributions
4. Cause and effect relationships
5. Summary points for each major topic
```

### Quiz Generation Template
```
Generate a [QUIZ_TYPE] quiz with [NUM_QUESTIONS] questions based on this textbook content. For a [ACADEMIC_LEVEL] level, include:
- Multiple choice questions (4 options each)
- True/False questions
- Short answer questions
- Answer key with explanations
```

### Flashcard Template
```
Create flashcards from this textbook content for a [ACADEMIC_LEVEL] student. Format:
Front: [Term/Question/Concept]
Back: [Definition/Answer/Explanation]

Focus on key vocabulary, important dates, significant figures, and core concepts.
```

## User Workflows

### Workflow 1: Asking Specific Questions
1. User uploads PDF and selects page range
2. Clicks "Ask Question" tab
3. Types: "Where did Columbus go in 1492?"
4. Selects academic level: "Undergraduate"
5. Clicks "Generate Response"
6. Receives detailed answer with context
7. Downloads response as PDF for studying

### Workflow 2: Creating Study Materials
1. User selects content pages
2. Chooses "Study Notes" option
3. AI generates organized notes with:
   - Topic headings
   - Bullet points
   - Key terms highlighted
   - Important dates
4. User reviews and downloads PDF

### Workflow 3: Quiz Generation
1. User selects chapter/pages
2. Chooses "Generate Quiz"
3. Specifies number of questions (e.g., 20)
4. AI creates varied question types
5. User can:
   - Take quiz interactively
   - Download with answer key
   - Share with classmates

## Database Schema Updates

```sql
-- Add new table for custom prompts
CREATE TABLE custom_prompts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT NOT NULL,
    document_id INTEGER NOT NULL,
    prompt TEXT NOT NULL,
    prompt_type TEXT NOT NULL,
    academic_level TEXT,
    response TEXT NOT NULL,
    generation_time INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES sessions(id),
    FOREIGN KEY (document_id) REFERENCES documents(id)
);

-- Add PDF storage table
CREATE TABLE generated_pdfs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content_id INTEGER NOT NULL,
    pdf_data BLOB NOT NULL,
    file_name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (content_id) REFERENCES generated_content(id)
);
```

## UI/UX Considerations

### Interactive Elements
- **Auto-suggestions**: Show common questions as user types
- **History**: Recent prompts for quick re-use
- **Templates**: Pre-built prompts for common tasks
- **Examples**: Show example questions for each category

### Response Formatting
- **Markdown Support**: Rich text formatting
- **Code Highlighting**: For technical content
- **Tables**: For structured data
- **Lists**: Bulleted and numbered
- **Citations**: Link back to source pages

### Mobile Responsiveness
- Collapsible panels for space efficiency
- Touch-friendly buttons and inputs
- Swipe gestures for navigation
- Responsive PDF viewer

## Testing Strategy

### Unit Tests
- Prompt builder functions
- PDF generation
- Response parsing
- Template rendering

### Integration Tests
- End-to-end prompt submission
- PDF download functionality
- Academic level handling
- Multi-format responses

### User Acceptance Tests
- Question accuracy
- Response relevance
- PDF quality
- Performance under load

## Performance Considerations

### Caching Strategy
- Cache common questions and responses
- Store generated PDFs temporarily
- Reuse extracted text across queries
- Implement response pagination

### Optimization
- Lazy load PDF generator
- Stream large responses
- Compress stored PDFs
- Rate limit API calls

## Security Considerations

### Input Validation
- Sanitize user prompts
- Limit prompt length (e.g., 1000 chars)
- Validate academic level values
- Check page range boundaries

### Rate Limiting
- Max 10 custom prompts per minute
- Max 5 PDF generations per session
- Implement exponential backoff
- Monitor for abuse patterns

## Future Enhancements

### Phase 2 Features
- **Collaborative Study**: Share materials with peers
- **Progress Tracking**: Monitor learning progress
- **Spaced Repetition**: Scheduled review reminders
- **Multi-language Support**: Translate content
- **Voice Input**: Ask questions verbally

### Phase 3 Features
- **AI Tutor Mode**: Interactive teaching sessions
- **Concept Linking**: Connect related topics
- **Practice Exams**: Full-length test simulations
- **Study Groups**: Collaborative learning spaces
- **Mobile App**: Native iOS/Android apps

## Success Metrics

### User Engagement
- Average prompts per session
- PDF download rate
- Return user percentage
- Feature adoption rate

### Quality Metrics
- Response accuracy rate
- User satisfaction scores
- Response time (< 3 seconds)
- PDF generation success rate

### Business Metrics
- User retention
- Feature usage growth
- Support ticket reduction
- User testimonials

## Implementation Timeline

### Week 1-2: Backend Development
- Custom prompt API endpoint
- Prompt builder service
- Enhanced AI integration
- Database schema updates

### Week 3-4: Frontend Development
- Custom prompt UI
- Response display components
- PDF download functionality
- Template system

### Week 5: PDF Generation
- PDF generator service
- Formatting engine
- Download API
- Storage optimization

### Week 6: Testing & Polish
- Comprehensive testing
- UI/UX refinements
- Performance optimization
- Documentation

## Conclusion

This feature transforms StudyForge into a comprehensive AI-powered study assistant, enabling students to interact with their textbook content in dynamic ways. By combining custom prompts, multiple output formats, and PDF export capabilities, we create a versatile learning tool that adapts to individual study needs and preferences.