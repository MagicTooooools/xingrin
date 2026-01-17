package dto

import "time"

// UpdateWordlistContentRequest represents update wordlist content request
type UpdateWordlistContentRequest struct {
	Content string `json:"content" binding:"required"`
}

// WordlistResponse represents wordlist response
type WordlistResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	FilePath    string    `json:"filePath"`
	FileSize    int64     `json:"fileSize"`
	LineCount   int       `json:"lineCount"`
	FileHash    string    `json:"fileHash"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// WordlistContentResponse represents wordlist content response
type WordlistContentResponse struct {
	Content string `json:"content"`
}
