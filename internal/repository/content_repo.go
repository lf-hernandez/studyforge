package repository

import (
	"database/sql"
	"fmt"

	"studyforge/internal/models"
)

// ContentRepository handles generated content database operations
type ContentRepository struct {
	db *sql.DB
}

// NewContentRepository creates a new content repository
func NewContentRepository(db *sql.DB) *ContentRepository {
	return &ContentRepository{db: db}
}

// CreateGenerated creates a new generated content record
func (r *ContentRepository) CreateGenerated(content *models.GeneratedContent) error {
	query := `
		INSERT INTO generated_content (session_id, document_id, content_type, academic_level, input_pages, output_content, ai_model, generation_time, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query,
		content.SessionID,
		content.DocumentID,
		content.ContentType,
		content.AcademicLevel,
		content.InputPages,
		content.OutputContent,
		content.AIModel,
		content.GenerationTime,
		content.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create generated content: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get content ID: %w", err)
	}
	content.ID = int(id)

	return nil
}

// GetGeneratedByID retrieves generated content by ID
func (r *ContentRepository) GetGeneratedByID(id int) (*models.GeneratedContent, error) {
	query := `
		SELECT id, session_id, document_id, content_type, academic_level, input_pages, output_content, ai_model, generation_time, created_at
		FROM generated_content
		WHERE id = ?
	`
	content := &models.GeneratedContent{}
	err := r.db.QueryRow(query, id).Scan(
		&content.ID,
		&content.SessionID,
		&content.DocumentID,
		&content.ContentType,
		&content.AcademicLevel,
		&content.InputPages,
		&content.OutputContent,
		&content.AIModel,
		&content.GenerationTime,
		&content.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("content not found")
		}
		return nil, fmt.Errorf("failed to get content: %w", err)
	}
	return content, nil
}

// CreateExtracted creates a new extracted content record (cache)
func (r *ContentRepository) CreateExtracted(content *models.ExtractedContent) error {
	query := `
		INSERT INTO extracted_content (document_id, page_start, page_end, content, extraction_time, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(document_id, page_start, page_end) DO UPDATE SET
			content = excluded.content,
			extraction_time = excluded.extraction_time,
			created_at = excluded.created_at
	`
	result, err := r.db.Exec(query,
		content.DocumentID,
		content.PageStart,
		content.PageEnd,
		content.Content,
		content.ExtractionTime,
		content.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create extracted content: %w", err)
	}

	id, err := result.LastInsertId()
	if err == nil {
		content.ID = int(id)
	}

	return nil
}

// GetExtracted retrieves cached extracted content
func (r *ContentRepository) GetExtracted(documentID, pageStart, pageEnd int) (*models.ExtractedContent, error) {
	query := `
		SELECT id, document_id, page_start, page_end, content, extraction_time, created_at
		FROM extracted_content
		WHERE document_id = ? AND page_start = ? AND page_end = ?
	`
	content := &models.ExtractedContent{}
	err := r.db.QueryRow(query, documentID, pageStart, pageEnd).Scan(
		&content.ID,
		&content.DocumentID,
		&content.PageStart,
		&content.PageEnd,
		&content.Content,
		&content.ExtractionTime,
		&content.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No cache found is not an error
		}
		return nil, fmt.Errorf("failed to get extracted content: %w", err)
	}
	return content, nil
}
