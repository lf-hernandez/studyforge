package repository

import (
	"database/sql"
	"fmt"
	"time"

	"studyforge/internal/models"
)

// DocumentRepository handles document database operations
type DocumentRepository struct {
	db *sql.DB
}

// NewDocumentRepository creates a new document repository
func NewDocumentRepository(db *sql.DB) *DocumentRepository {
	return &DocumentRepository{db: db}
}

// Create creates a new document record
func (r *DocumentRepository) Create(doc *models.Document) error {
	query := `
		INSERT INTO documents (session_id, original_filename, stored_filename, file_path, file_size, page_count, upload_date)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query,
		doc.SessionID,
		doc.OriginalFilename,
		doc.StoredFilename,
		doc.FilePath,
		doc.FileSize,
		doc.PageCount,
		doc.UploadDate,
	)
	if err != nil {
		return fmt.Errorf("failed to create document: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get document ID: %w", err)
	}
	doc.ID = int(id)

	return nil
}

// GetByID retrieves a document by ID
func (r *DocumentRepository) GetByID(id int) (*models.Document, error) {
	query := `
		SELECT id, session_id, original_filename, stored_filename, file_path, file_size, page_count, upload_date, last_accessed, is_deleted
		FROM documents
		WHERE id = ? AND is_deleted = FALSE
	`
	doc := &models.Document{}
	var lastAccessed sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&doc.ID,
		&doc.SessionID,
		&doc.OriginalFilename,
		&doc.StoredFilename,
		&doc.FilePath,
		&doc.FileSize,
		&doc.PageCount,
		&doc.UploadDate,
		&lastAccessed,
		&doc.IsDeleted,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("document not found")
		}
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	if lastAccessed.Valid {
		doc.LastAccessed = lastAccessed.Time
	}

	return doc, nil
}

// GetBySessionID retrieves all documents for a session
func (r *DocumentRepository) GetBySessionID(sessionID string) ([]*models.Document, error) {
	query := `
		SELECT id, session_id, original_filename, stored_filename, file_path, file_size, page_count, upload_date, last_accessed, is_deleted
		FROM documents
		WHERE session_id = ? AND is_deleted = FALSE
		ORDER BY upload_date DESC
	`
	rows, err := r.db.Query(query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get documents: %w", err)
	}
	defer rows.Close()

	var documents []*models.Document
	for rows.Next() {
		doc := &models.Document{}
		var lastAccessed sql.NullTime

		err := rows.Scan(
			&doc.ID,
			&doc.SessionID,
			&doc.OriginalFilename,
			&doc.StoredFilename,
			&doc.FilePath,
			&doc.FileSize,
			&doc.PageCount,
			&doc.UploadDate,
			&lastAccessed,
			&doc.IsDeleted,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan document: %w", err)
		}

		if lastAccessed.Valid {
			doc.LastAccessed = lastAccessed.Time
		}

		documents = append(documents, doc)
	}

	return documents, nil
}

// UpdateLastAccessed updates the last accessed timestamp
func (r *DocumentRepository) UpdateLastAccessed(id int) error {
	query := `UPDATE documents SET last_accessed = ? WHERE id = ?`
	_, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update last accessed: %w", err)
	}
	return nil
}

// Delete marks a document as deleted
func (r *DocumentRepository) Delete(id int) error {
	query := `UPDATE documents SET is_deleted = TRUE WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	return nil
}
