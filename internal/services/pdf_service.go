package services

import (
	"fmt"
	"time"

	"studyforge/internal/models"
	"studyforge/internal/repository"
	"studyforge/pkg/pdf"
)

// PDFService handles PDF-related business logic
type PDFService struct {
	extractor  *pdf.Extractor
	contentRepo *repository.ContentRepository
}

// NewPDFService creates a new PDF service
func NewPDFService(contentRepo *repository.ContentRepository) *PDFService {
	return &PDFService{
		extractor:   pdf.NewExtractor(),
		contentRepo: contentRepo,
	}
}

// GetPageCount returns the number of pages in a PDF
func (s *PDFService) GetPageCount(filePath string) (int, error) {
	return s.extractor.GetPageCount(filePath)
}

// ExtractText extracts text from specified page range with caching
func (s *PDFService) ExtractText(documentID int, filePath string, startPage, endPage int) (string, int, error) {
	startTime := time.Now()

	// Check cache first
	cached, err := s.contentRepo.GetExtracted(documentID, startPage, endPage)
	if err != nil {
		return "", 0, fmt.Errorf("cache lookup failed: %w", err)
	}

	if cached != nil {
		// Return cached content
		return cached.Content, cached.ExtractionTime, nil
	}

	// Extract text from PDF
	text, err := s.extractor.ExtractText(filePath, startPage, endPage)
	if err != nil {
		return "", 0, err
	}

	extractionTime := int(time.Since(startTime).Milliseconds())

	// Cache the result
	extractedContent := &models.ExtractedContent{
		DocumentID:     documentID,
		PageStart:      startPage,
		PageEnd:        endPage,
		Content:        text,
		ExtractionTime: extractionTime,
		CreatedAt:      time.Now(),
	}

	if err := s.contentRepo.CreateExtracted(extractedContent); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to cache extracted content: %v\n", err)
	}

	return text, extractionTime, nil
}

// ValidatePageRange validates a page range
func (s *PDFService) ValidatePageRange(filePath string, startPage, endPage int) error {
	return s.extractor.ValidatePageRange(filePath, startPage, endPage)
}
