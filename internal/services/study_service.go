package services

import (
	"encoding/json"
	"fmt"
	"time"

	"studyforge/internal/models"
	"studyforge/internal/repository"
	"studyforge/pkg/ai"
)

// StudyService handles study material generation
type StudyService struct {
	aiClient    *ai.HuggingFaceClient
	pdfService  *PDFService
	contentRepo *repository.ContentRepository
	docRepo     *repository.DocumentRepository
}

// NewStudyService creates a new study service
func NewStudyService(
	aiClient *ai.HuggingFaceClient,
	pdfService *PDFService,
	contentRepo *repository.ContentRepository,
	docRepo *repository.DocumentRepository,
) *StudyService {
	return &StudyService{
		aiClient:    aiClient,
		pdfService:  pdfService,
		contentRepo: contentRepo,
		docRepo:     docRepo,
	}
}

// GenerateSummaryRequest contains parameters for summary generation
type GenerateSummaryRequest struct {
	SessionID     string
	DocumentID    int
	PageStart     int
	PageEnd       int
	AcademicLevel string
}

// GenerateSummaryResponse contains the generated summary
type GenerateSummaryResponse struct {
	ContentID      int    `json:"content_id"`
	Summary        string `json:"summary"`
	GenerationTime int    `json:"generation_time"`
	ModelUsed      string `json:"model_used"`
}

// GenerateSummary generates a summary from specified pages
func (s *StudyService) GenerateSummary(req *GenerateSummaryRequest) (*GenerateSummaryResponse, error) {
	startTime := time.Now()

	// Get document
	doc, err := s.docRepo.GetByID(req.DocumentID)
	if err != nil {
		return nil, fmt.Errorf("document not found: %w", err)
	}

	// Verify session matches
	if doc.SessionID != req.SessionID {
		return nil, fmt.Errorf("unauthorized access to document")
	}

	// Validate page range
	if err := s.pdfService.ValidatePageRange(doc.FilePath, req.PageStart, req.PageEnd); err != nil {
		return nil, err
	}

	// Extract text from PDF
	text, _, err := s.pdfService.ExtractText(req.DocumentID, doc.FilePath, req.PageStart, req.PageEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to extract text: %w", err)
	}

	// Generate summary using AI
	summary, err := s.aiClient.GenerateSummary(text, req.AcademicLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to generate summary: %w", err)
	}

	generationTime := int(time.Since(startTime).Milliseconds())

	// Create output content structure
	outputData := map[string]interface{}{
		"summary":        summary,
		"pages":          fmt.Sprintf("%d-%d", req.PageStart, req.PageEnd),
		"academic_level": req.AcademicLevel,
	}

	outputJSON, err := json.Marshal(outputData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal output: %w", err)
	}

	// Save to database
	generatedContent := &models.GeneratedContent{
		SessionID:      req.SessionID,
		DocumentID:     req.DocumentID,
		ContentType:    "summary",
		AcademicLevel:  req.AcademicLevel,
		InputPages:     fmt.Sprintf("%d-%d", req.PageStart, req.PageEnd),
		OutputContent:  string(outputJSON),
		AIModel:        "facebook/bart-large-cnn",
		GenerationTime: generationTime,
		CreatedAt:      time.Now(),
	}

	if err := s.contentRepo.CreateGenerated(generatedContent); err != nil {
		return nil, fmt.Errorf("failed to save content: %w", err)
	}

	return &GenerateSummaryResponse{
		ContentID:      generatedContent.ID,
		Summary:        summary,
		GenerationTime: generationTime,
		ModelUsed:      "facebook/bart-large-cnn",
	}, nil
}

// GetGeneratedContent retrieves previously generated content
func (s *StudyService) GetGeneratedContent(contentID int, sessionID string) (*models.GeneratedContent, error) {
	content, err := s.contentRepo.GetGeneratedByID(contentID)
	if err != nil {
		return nil, err
	}

	// Verify session
	if content.SessionID != sessionID {
		return nil, fmt.Errorf("unauthorized access to content")
	}

	return content, nil
}
