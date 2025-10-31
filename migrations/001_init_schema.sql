-- StudyForge Database Schema
-- Migration 001: Initial Schema

-- Sessions table for user tracking
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_accessed TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ip_address TEXT,
    user_agent TEXT,
    is_active BOOLEAN DEFAULT TRUE
);

-- Documents table for uploaded PDFs
CREATE TABLE IF NOT EXISTS documents (
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
    FOREIGN KEY (session_id) REFERENCES sessions(id)
);

-- Extracted content cache
CREATE TABLE IF NOT EXISTS extracted_content (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    document_id INTEGER NOT NULL,
    page_start INTEGER NOT NULL,
    page_end INTEGER NOT NULL,
    content TEXT NOT NULL,
    extraction_time INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (document_id) REFERENCES documents(id),
    UNIQUE(document_id, page_start, page_end)
);

-- Generated study materials (MVP: summaries only)
CREATE TABLE IF NOT EXISTS generated_content (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT NOT NULL,
    document_id INTEGER NOT NULL,
    content_type TEXT NOT NULL,
    academic_level TEXT,
    input_pages TEXT,
    output_content TEXT NOT NULL,
    ai_model TEXT NOT NULL,
    generation_time INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES sessions(id),
    FOREIGN KEY (document_id) REFERENCES documents(id)
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_sessions_active ON sessions(is_active);
CREATE INDEX IF NOT EXISTS idx_documents_session ON documents(session_id);
CREATE INDEX IF NOT EXISTS idx_extracted_content_document ON extracted_content(document_id);
CREATE INDEX IF NOT EXISTS idx_generated_content_session ON generated_content(session_id);
CREATE INDEX IF NOT EXISTS idx_generated_content_document ON generated_content(document_id);
