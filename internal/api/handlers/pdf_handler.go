package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"studyforge/internal/config"
	"studyforge/internal/models"
	"studyforge/internal/repository"
	"studyforge/internal/services"
	"studyforge/pkg/utils"

	"github.com/google/uuid"
)

// PDFHandler handles PDF-related requests
type PDFHandler struct {
	cfg        *config.Config
	docRepo    *repository.DocumentRepository
	pdfService *services.PDFService
}

// NewPDFHandler creates a new PDF handler
func NewPDFHandler(
	cfg *config.Config,
	docRepo *repository.DocumentRepository,
	pdfService *services.PDFService,
) *PDFHandler {
	return &PDFHandler{
		cfg:        cfg,
		docRepo:    docRepo,
		pdfService: pdfService,
	}
}

// HandleUpload handles PDF file uploads
func (h *PDFHandler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	// Get session from context
	session, err := utils.GetSessionFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "NO_SESSION", "No session found")
		return
	}

	// Parse multipart form (limit to configured max file size)
	if err := r.ParseMultipartForm(h.cfg.MaxFileSize); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "PARSE_ERROR", "Failed to parse form data")
		return
	}

	// Get file from form
	file, header, err := r.FormFile("file")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "NO_FILE", "No file provided")
		return
	}
	defer file.Close()

	// Validate file size
	if header.Size > h.cfg.MaxFileSize {
		utils.WriteError(w, http.StatusBadRequest, "FILE_TOO_LARGE",
			fmt.Sprintf("File size exceeds %d MB limit", h.cfg.MaxFileSize/1024/1024))
		return
	}

	// Validate file type (must be PDF)
	if filepath.Ext(header.Filename) != ".pdf" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_FILE_TYPE", "Only PDF files are allowed")
		return
	}

	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(h.cfg.UploadDir, 0755); err != nil {
		log.Printf("Failed to create upload directory: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "UPLOAD_ERROR", "Failed to prepare upload")
		return
	}

	// Generate unique filename
	storedFilename := fmt.Sprintf("%s_%s", uuid.New().String(), filepath.Base(header.Filename))
	filePath := filepath.Join(h.cfg.UploadDir, storedFilename)

	// Save file
	dst, err := os.Create(filePath)
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "UPLOAD_ERROR", "Failed to save file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		log.Printf("Failed to copy file: %v", err)
		os.Remove(filePath) // Clean up
		utils.WriteError(w, http.StatusInternalServerError, "UPLOAD_ERROR", "Failed to save file")
		return
	}

	// Get page count
	pageCount, err := h.pdfService.GetPageCount(filePath)
	if err != nil {
		log.Printf("Failed to get page count: %v", err)
		os.Remove(filePath) // Clean up
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PDF", "Failed to process PDF file")
		return
	}

	// Create document record
	doc := &models.Document{
		SessionID:        session.ID,
		OriginalFilename: header.Filename,
		StoredFilename:   storedFilename,
		FilePath:         filePath,
		FileSize:         header.Size,
		PageCount:        pageCount,
		UploadDate:       time.Now(),
		IsDeleted:        false,
	}

	if err := h.docRepo.Create(doc); err != nil {
		log.Printf("Failed to save document record: %v", err)
		os.Remove(filePath) // Clean up
		utils.WriteError(w, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to save document")
		return
	}

	log.Printf("Document uploaded successfully: %s (%d pages)", header.Filename, pageCount)

	// Return response
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"document_id": doc.ID,
		"filename":    doc.OriginalFilename,
		"page_count":  doc.PageCount,
		"file_size":   doc.FileSize,
	})
}

// HandleGetDocument retrieves document information
func (h *PDFHandler) HandleGetDocument(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	// Get session
	session, err := utils.GetSessionFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "NO_SESSION", "No session found")
		return
	}

	// Get document ID from URL
	docIDStr := r.URL.Query().Get("id")
	docID, err := strconv.Atoi(docIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid document ID")
		return
	}

	// Get document
	doc, err := h.docRepo.GetByID(docID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "NOT_FOUND", "Document not found")
		return
	}

	// Verify session
	if doc.SessionID != session.ID {
		utils.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized access")
		return
	}

	// Return document info
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id":          doc.ID,
		"filename":    doc.OriginalFilename,
		"page_count":  doc.PageCount,
		"file_size":   doc.FileSize,
		"upload_date": doc.UploadDate,
	})
}
