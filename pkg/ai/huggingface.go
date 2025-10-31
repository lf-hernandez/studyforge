package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HuggingFaceClient handles communication with Hugging Face API
type HuggingFaceClient struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// NewHuggingFaceClient creates a new Hugging Face API client
func NewHuggingFaceClient(apiKey, baseURL string) *HuggingFaceClient {
	return &HuggingFaceClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SummaryRequest represents a request to generate a summary
type SummaryRequest struct {
	Inputs     string                 `json:"inputs"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// SummaryResponse represents the API response
type SummaryResponse struct {
	SummaryText string `json:"summary_text"`
}

// GenerateSummary generates a summary using the BART model with chunking
func (c *HuggingFaceClient) GenerateSummary(text string, academicLevel string) (string, error) {
	// BART can handle ~1024 tokens, which is roughly 3000-4000 characters
	// We'll use 3000 as a safe limit per chunk
	maxChunkSize := 3000

	// If text is small enough, summarize directly
	if len(text) <= maxChunkSize {
		return c.summarizeChunk(text, academicLevel)
	}

	// Otherwise, chunk the text and summarize each chunk
	chunks := c.chunkText(text, maxChunkSize)
	fmt.Printf("Text too large (%d chars), splitting into %d chunks\n", len(text), len(chunks))

	var chunkSummaries []string
	for i, chunk := range chunks {
		fmt.Printf("Summarizing chunk %d/%d (%d chars)...\n", i+1, len(chunks), len(chunk))

		summary, err := c.summarizeChunk(chunk, academicLevel)
		if err != nil {
			return "", fmt.Errorf("failed to summarize chunk %d: %w", i+1, err)
		}

		chunkSummaries = append(chunkSummaries, summary)

		// Small delay between requests to avoid rate limiting
		if i < len(chunks)-1 {
			time.Sleep(500 * time.Millisecond)
		}
	}

	// If we have multiple chunk summaries, combine them
	if len(chunkSummaries) > 1 {
		combinedText := ""
		for i, summary := range chunkSummaries {
			combinedText += fmt.Sprintf("Section %d: %s\n\n", i+1, summary)
		}

		// Return combined summaries directly
		// Note: We could re-summarize if too long, but for now just return sections
		fmt.Printf("Combined %d summaries into final result (%d chars)\n", len(chunkSummaries), len(combinedText))
		return combinedText, nil
	}

	return chunkSummaries[0], nil
}

// chunkText splits text into chunks of roughly equal size
func (c *HuggingFaceClient) chunkText(text string, maxSize int) []string {
	var chunks []string
	minChunkSize := 200 // Minimum chunk size to avoid tiny chunks

	// Split by paragraphs first (double newlines)
	paragraphs := splitByParagraphs(text)

	currentChunk := ""
	for _, para := range paragraphs {
		// Skip very short paragraphs (likely page markers)
		if len(para) < 20 {
			continue
		}

		// If adding this paragraph would exceed the limit, save current chunk
		if len(currentChunk)+len(para) > maxSize && len(currentChunk) >= minChunkSize {
			chunks = append(chunks, currentChunk)
			currentChunk = para
		} else {
			if len(currentChunk) > 0 {
				currentChunk += " "
			}
			currentChunk += para
		}
	}

	// Add the last chunk if it's substantial
	if len(currentChunk) >= minChunkSize {
		chunks = append(chunks, currentChunk)
	}

	return chunks
}

// buildEducationalPrompt creates an instructional prompt for educational content summarization
func buildEducationalPrompt(text string, academicLevel string) string {
	// Build instruction based on academic level
	var instruction string

	switch academicLevel {
	case "high_school":
		instruction = "Summarize this history textbook content for high school students. Use clear language and focus on main events, key people, and important dates. Explain cause and effect relationships. Ignore citations and web references. "
	case "undergraduate":
		instruction = "Summarize this history textbook content for undergraduate college students. Focus on key concepts, historical events, significant figures, and their impact. Include important context and connections between events. Ignore citations and web references. "
	case "graduate":
		instruction = "Summarize this history textbook content for graduate students. Emphasize analytical perspectives, historiographical significance, and complex interrelationships between events and themes. Ignore citations and web references. "
	default:
		instruction = "Summarize this educational textbook content. Focus on key concepts, historical events, and important facts. Maintain accuracy and clarity. Ignore citations and web references. "
	}

	// Combine instruction with content
	return instruction + text
}

// summarizeChunk summarizes a single chunk of text
func (c *HuggingFaceClient) summarizeChunk(text string, academicLevel string) (string, error) {
	// Build instructional prompt for educational summarization
	prompt := buildEducationalPrompt(text, academicLevel)

	// Prepare request
	reqBody := SummaryRequest{
		Inputs: prompt,
		Parameters: map[string]interface{}{
			"max_length": 350,
			"min_length": 120,
			"do_sample":  false,
		},
	}

	// Use BART model for summarization
	modelURL := fmt.Sprintf("%s/facebook/bart-large-cnn", c.baseURL)

	responseData, err := c.makeRequest(modelURL, reqBody)
	if err != nil {
		return "", err
	}

	// Parse response
	var responses []SummaryResponse
	if err := json.Unmarshal(responseData, &responses); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(responses) == 0 {
		return "", fmt.Errorf("empty response from API")
	}

	return responses[0].SummaryText, nil
}

// splitByParagraphs splits text into paragraphs
func splitByParagraphs(text string) []string {
	// Split by double newlines or page markers
	var paragraphs []string
	current := ""

	for _, line := range splitLines(text) {
		if line == "" || line == "---" {
			if len(current) > 0 {
				paragraphs = append(paragraphs, current)
				current = ""
			}
		} else {
			if len(current) > 0 {
				current += " "
			}
			current += line
		}
	}

	if len(current) > 0 {
		paragraphs = append(paragraphs, current)
	}

	return paragraphs
}

// splitLines splits text into lines
func splitLines(text string) []string {
	var lines []string
	current := ""

	for _, char := range text {
		if char == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(char)
		}
	}

	if len(current) > 0 {
		lines = append(lines, current)
	}

	return lines
}

// makeRequest makes an HTTP request to the Hugging Face API
func (c *HuggingFaceClient) makeRequest(url string, reqBody interface{}) ([]byte, error) {
	// Marshal request body
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	// Make request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}
