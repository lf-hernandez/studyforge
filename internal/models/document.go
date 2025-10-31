package models

import "time"

// Document represents an uploaded PDF document
type Document struct {
	ID               int       `json:"id"`
	SessionID        string    `json:"session_id"`
	OriginalFilename string    `json:"original_filename"`
	StoredFilename   string    `json:"stored_filename"`
	FilePath         string    `json:"file_path"`
	FileSize         int64     `json:"file_size"`
	PageCount        int       `json:"page_count"`
	UploadDate       time.Time `json:"upload_date"`
	LastAccessed     time.Time `json:"last_accessed"`
	IsDeleted        bool      `json:"is_deleted"`
}
