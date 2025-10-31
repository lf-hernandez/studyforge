package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"studyforge/internal/services"
	"studyforge/pkg/utils"
)

// StudyHandler handles study material generation requests
type StudyHandler struct {
	studyService *services.StudyService
}

// NewStudyHandler creates a new study handler
func NewStudyHandler(studyService *services.StudyService) *StudyHandler {
	return &StudyHandler{
		studyService: studyService,
	}
}

// GenerateRequest represents a study material generation request
type GenerateRequest struct {
	DocumentID    int    `json:"document_id"`
	PageStart     int    `json:"page_start"`
	PageEnd       int    `json:"page_end"`
	MaterialType  string `json:"material_type"`   // 'summary' for MVP
	AcademicLevel string `json:"academic_level"`  // 'high_school', 'undergraduate', 'graduate'
}

// HandleGenerate handles study material generation
func (h *StudyHandler) HandleGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	// Get session
	session, err := utils.GetSessionFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "NO_SESSION", "No session found")
		return
	}

	// Parse request body
	var req GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	// Validate request
	if req.DocumentID <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_DOCUMENT_ID", "Invalid document ID")
		return
	}
	if req.PageStart < 1 || req.PageEnd < req.PageStart {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PAGE_RANGE", "Invalid page range")
		return
	}

	// MVP: Only support summaries
	if req.MaterialType != "summary" && req.MaterialType != "" {
		utils.WriteError(w, http.StatusBadRequest, "UNSUPPORTED_TYPE", "Only 'summary' type is supported in MVP")
		return
	}
	req.MaterialType = "summary"

	// Default academic level if not provided
	if req.AcademicLevel == "" {
		req.AcademicLevel = "undergraduate"
	}

	// Generate summary
	serviceReq := &services.GenerateSummaryRequest{
		SessionID:     session.ID,
		DocumentID:    req.DocumentID,
		PageStart:     req.PageStart,
		PageEnd:       req.PageEnd,
		AcademicLevel: req.AcademicLevel,
	}

	log.Printf("Generating summary for document %d, pages %d-%d", req.DocumentID, req.PageStart, req.PageEnd)

	result, err := h.studyService.GenerateSummary(serviceReq)
	if err != nil {
		log.Printf("Failed to generate summary: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "GENERATION_ERROR", err.Error())
		return
	}

	log.Printf("Summary generated successfully (ID: %d)", result.ContentID)

	// Return response
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"content_id":      result.ContentID,
		"material_type":   "summary",
		"summary":         result.Summary,
		"model_used":      result.ModelUsed,
		"generation_time": result.GenerationTime,
	})
}

// HandleGetContent retrieves previously generated content
func (h *StudyHandler) HandleGetContent(w http.ResponseWriter, r *http.Request) {
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

	// Get content ID from URL
	contentIDStr := r.URL.Query().Get("id")
	contentID, err := strconv.Atoi(contentIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid content ID")
		return
	}

	// Get content
	content, err := h.studyService.GetGeneratedContent(contentID, session.ID)
	if err != nil {
		log.Printf("Failed to get content: %v", err)
		utils.WriteError(w, http.StatusNotFound, "NOT_FOUND", "Content not found")
		return
	}

	// Parse output content
	var outputData map[string]interface{}
	if err := json.Unmarshal([]byte(content.OutputContent), &outputData); err != nil {
		log.Printf("Failed to parse output content: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "PARSE_ERROR", "Failed to parse content")
		return
	}

	// Return content
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"content_id":      content.ID,
		"material_type":   content.ContentType,
		"content":         outputData,
		"academic_level":  content.AcademicLevel,
		"pages":           content.InputPages,
		"model_used":      content.AIModel,
		"generation_time": content.GenerationTime,
		"created_at":      content.CreatedAt,
	})
}
