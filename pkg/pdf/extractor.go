package pdf

import (
	"fmt"
	"strings"

	"studyforge/pkg/utils"

	"github.com/ledongthuc/pdf"
)

// Extractor handles PDF text extraction
type Extractor struct{}

// NewExtractor creates a new PDF extractor
func NewExtractor() *Extractor {
	return &Extractor{}
}

// GetPageCount returns the number of pages in a PDF
func (e *Extractor) GetPageCount(filePath string) (int, error) {
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to open PDF: %w", err)
	}
	defer f.Close()

	return r.NumPage(), nil
}

// ExtractText extracts text from specified page range
// Pages are 1-indexed (first page is 1)
func (e *Extractor) ExtractText(filePath string, startPage, endPage int) (string, error) {
	// Validate page range
	if startPage < 1 || endPage < startPage {
		return "", fmt.Errorf("invalid page range: %d-%d", startPage, endPage)
	}

	// Open PDF
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF: %w", err)
	}
	defer f.Close()

	// Get page count
	numPages := r.NumPage()

	// Validate page range against document
	if endPage > numPages {
		return "", fmt.Errorf("end page %d exceeds document pages %d", endPage, numPages)
	}

	// Extract text from each page
	var extractedText strings.Builder

	for pageNum := startPage; pageNum <= endPage; pageNum++ {
		page := r.Page(pageNum)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			return "", fmt.Errorf("failed to extract page %d: %w", pageNum, err)
		}

		// Clean PDF extraction artifacts
		cleanedText := utils.CleanPDFText(text)

		extractedText.WriteString(fmt.Sprintf("--- Page %d ---\n", pageNum))
		extractedText.WriteString(cleanedText)
		extractedText.WriteString("\n\n")
	}

	result := extractedText.String()
	if strings.TrimSpace(result) == "" {
		return "", fmt.Errorf("no text could be extracted (PDF may be image-based or scanned)")
	}

	return result, nil
}

// ValidatePageRange validates a page range against a document
func (e *Extractor) ValidatePageRange(filePath string, startPage, endPage int) error {
	if startPage < 1 {
		return fmt.Errorf("start page must be at least 1")
	}
	if endPage < startPage {
		return fmt.Errorf("end page must be greater than or equal to start page")
	}

	f, r, err := pdf.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open PDF: %w", err)
	}
	defer f.Close()

	pageCount := r.NumPage()
	if endPage > pageCount {
		return fmt.Errorf("end page %d exceeds document pages %d", endPage, pageCount)
	}

	return nil
}
