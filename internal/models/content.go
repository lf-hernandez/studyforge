package models

import "time"

// GeneratedContent represents AI-generated study material
type GeneratedContent struct {
	ID             int       `json:"id"`
	SessionID      string    `json:"session_id"`
	DocumentID     int       `json:"document_id"`
	ContentType    string    `json:"content_type"`    // 'summary', 'quiz', 'notes', etc.
	AcademicLevel  string    `json:"academic_level"`  // 'high_school', 'undergraduate', 'graduate'
	InputPages     string    `json:"input_pages"`     // e.g., "1-10"
	OutputContent  string    `json:"output_content"`  // JSON string
	AIModel        string    `json:"ai_model"`
	GenerationTime int       `json:"generation_time"` // milliseconds
	CreatedAt      time.Time `json:"created_at"`
}

// ExtractedContent represents cached extracted text from PDF
type ExtractedContent struct {
	ID             int       `json:"id"`
	DocumentID     int       `json:"document_id"`
	PageStart      int       `json:"page_start"`
	PageEnd        int       `json:"page_end"`
	Content        string    `json:"content"`
	ExtractionTime int       `json:"extraction_time"` // milliseconds
	CreatedAt      time.Time `json:"created_at"`
}
