package handler

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yyhuni/lunafox/server/internal/dto"
	"github.com/yyhuni/lunafox/server/internal/service"
)

// WordlistHandler handles wordlist API requests
type WordlistHandler struct {
	svc *service.WordlistService
}

// NewWordlistHandler creates a new wordlist handler
func NewWordlistHandler(svc *service.WordlistService) *WordlistHandler {
	return &WordlistHandler{
		svc: svc,
	}
}

// List returns all wordlists
// GET /api/wordlists/
func (h *WordlistHandler) List(c *gin.Context) {
	wordlists, err := h.svc.ListAll()
	if err != nil {
		dto.InternalError(c, "Failed to list wordlists")
		return
	}

	// Initialize empty slice to return [] instead of null
	resp := make([]dto.WordlistResponse, 0, len(wordlists))
	for _, w := range wordlists {
		resp = append(resp, dto.WordlistResponse{
			ID:          w.ID,
			Name:        w.Name,
			Description: w.Description,
			FilePath:    w.FilePath,
			FileSize:    w.FileSize,
			LineCount:   w.LineCount,
			FileHash:    w.FileHash,
			CreatedAt:   w.CreatedAt,
			UpdatedAt:   w.UpdatedAt,
		})
	}

	dto.Success(c, resp)
}

// Get returns a wordlist by ID
// GET /api/wordlists/:id
func (h *WordlistHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid wordlist ID")
		return
	}

	wordlist, err := h.svc.GetByID(id)
	if err != nil {
		if errors.Is(err, service.ErrWordlistNotFound) {
			dto.NotFound(c, "Wordlist not found")
			return
		}
		dto.InternalError(c, "Failed to get wordlist")
		return
	}

	dto.Success(c, dto.WordlistResponse{
		ID:          wordlist.ID,
		Name:        wordlist.Name,
		Description: wordlist.Description,
		FilePath:    wordlist.FilePath,
		FileSize:    wordlist.FileSize,
		LineCount:   wordlist.LineCount,
		FileHash:    wordlist.FileHash,
		CreatedAt:   wordlist.CreatedAt,
		UpdatedAt:   wordlist.UpdatedAt,
	})
}

// GetByName returns a wordlist by name (for worker API)
// GET /api/worker/wordlists/:name
func (h *WordlistHandler) GetByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		dto.BadRequest(c, "Name is required")
		return
	}

	wordlist, err := h.svc.GetByName(name)
	if err != nil {
		if errors.Is(err, service.ErrWordlistNotFound) {
			dto.NotFound(c, "Wordlist not found")
			return
		}
		dto.InternalError(c, "Failed to get wordlist")
		return
	}

	dto.Success(c, dto.WordlistResponse{
		ID:          wordlist.ID,
		Name:        wordlist.Name,
		Description: wordlist.Description,
		FilePath:    wordlist.FilePath,
		FileSize:    wordlist.FileSize,
		LineCount:   wordlist.LineCount,
		FileHash:    wordlist.FileHash,
		CreatedAt:   wordlist.CreatedAt,
		UpdatedAt:   wordlist.UpdatedAt,
	})
}

// DownloadByID serves the wordlist file by ID (RESTful)
// GET /api/wordlists/:id/download
func (h *WordlistHandler) DownloadByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid wordlist ID")
		return
	}

	wordlist, err := h.svc.GetByID(id)
	if err != nil {
		if errors.Is(err, service.ErrWordlistNotFound) {
			dto.NotFound(c, "Wordlist not found")
			return
		}
		dto.InternalError(c, "Failed to get wordlist")
		return
	}

	h.serveWordlistFile(c, wordlist.Name)
}

// DownloadByName serves the wordlist file (path parameter style - RESTful)
// GET /api/worker/wordlists/:name/download
func (h *WordlistHandler) DownloadByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		dto.BadRequest(c, "Name is required")
		return
	}

	h.serveWordlistFile(c, name)
}

// serveWordlistFile is a helper to serve wordlist file by name
func (h *WordlistHandler) serveWordlistFile(c *gin.Context, name string) {
	filePath, err := h.svc.GetFilePath(name)
	if err != nil {
		if errors.Is(err, service.ErrWordlistNotFound) {
			dto.NotFound(c, "Wordlist not found")
			return
		}
		if errors.Is(err, service.ErrFileNotFound) {
			dto.NotFound(c, "Wordlist file not found on server")
			return
		}
		dto.InternalError(c, "Failed to get wordlist")
		return
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		dto.NotFound(c, "Wordlist file not found on server")
		return
	}

	// Serve file
	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(filePath))
	c.Header("Content-Type", "application/octet-stream")
	c.File(filePath)
}

// Create creates a new wordlist with file upload
// POST /api/wordlists/
// Content-Type: multipart/form-data
// Fields: name (required), description (optional), file (required)
func (h *WordlistHandler) Create(c *gin.Context) {
	// Get form fields
	name := c.PostForm("name")
	if name == "" {
		dto.BadRequest(c, "Name is required")
		return
	}

	description := c.PostForm("description")

	// Get uploaded file (streaming, not loaded into memory)
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		dto.BadRequest(c, "File is required")
		return
	}
	defer func() { _ = file.Close() }()

	// Create wordlist with streamed file content
	wordlist, err := h.svc.Create(name, description, header.Filename, file)
	if err != nil {
		if errors.Is(err, service.ErrWordlistExists) {
			dto.BadRequest(c, "Wordlist name already exists")
			return
		}
		if errors.Is(err, service.ErrEmptyName) {
			dto.BadRequest(c, "Wordlist name cannot be empty")
			return
		}
		if errors.Is(err, service.ErrNameTooLong) {
			dto.BadRequest(c, "Wordlist name too long (max 200 characters)")
			return
		}
		if errors.Is(err, service.ErrInvalidName) {
			dto.BadRequest(c, "Wordlist name contains invalid characters (newlines, tabs, etc.)")
			return
		}
		if errors.Is(err, service.ErrInvalidFileType) {
			dto.BadRequest(c, "File appears to be binary, only text files are allowed")
			return
		}
		dto.InternalError(c, "Failed to create wordlist")
		return
	}

	dto.Created(c, dto.WordlistResponse{
		ID:          wordlist.ID,
		Name:        wordlist.Name,
		Description: wordlist.Description,
		FilePath:    wordlist.FilePath,
		FileSize:    wordlist.FileSize,
		LineCount:   wordlist.LineCount,
		FileHash:    wordlist.FileHash,
		CreatedAt:   wordlist.CreatedAt,
		UpdatedAt:   wordlist.UpdatedAt,
	})
}

// Delete deletes a wordlist
// DELETE /api/wordlists/:id
func (h *WordlistHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid wordlist ID")
		return
	}

	if err := h.svc.Delete(id); err != nil {
		if errors.Is(err, service.ErrWordlistNotFound) {
			dto.NotFound(c, "Wordlist not found")
			return
		}
		dto.InternalError(c, "Failed to delete wordlist")
		return
	}

	dto.NoContent(c)
}

// maxEditableSize is the maximum file size allowed for online editing (10MB)
const maxEditableSize = 10 * 1024 * 1024

// GetContent returns the content of a wordlist file
// GET /api/wordlists/:id/content
func (h *WordlistHandler) GetContent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid wordlist ID")
		return
	}

	// Check file size first
	wordlist, err := h.svc.GetByID(id)
	if err != nil {
		if errors.Is(err, service.ErrWordlistNotFound) {
			dto.NotFound(c, "Wordlist not found")
			return
		}
		dto.InternalError(c, "Failed to get wordlist")
		return
	}

	if wordlist.FileSize > maxEditableSize {
		dto.BadRequest(c, "File too large for online editing (max 10MB), please download and edit locally")
		return
	}

	content, err := h.svc.GetContent(id)
	if err != nil {
		if errors.Is(err, service.ErrWordlistNotFound) {
			dto.NotFound(c, "Wordlist not found")
			return
		}
		if errors.Is(err, service.ErrFileNotFound) {
			dto.NotFound(c, "Wordlist file not found")
			return
		}
		dto.InternalError(c, "Failed to get wordlist content")
		return
	}

	dto.Success(c, dto.WordlistContentResponse{Content: content})
}

// UpdateContent updates the content of a wordlist file
// PUT /api/wordlists/:id/content
func (h *WordlistHandler) UpdateContent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid wordlist ID")
		return
	}

	// Check file size first
	wordlist, err := h.svc.GetByID(id)
	if err != nil {
		if errors.Is(err, service.ErrWordlistNotFound) {
			dto.NotFound(c, "Wordlist not found")
			return
		}
		dto.InternalError(c, "Failed to get wordlist")
		return
	}

	if wordlist.FileSize > maxEditableSize {
		dto.BadRequest(c, "File too large for online editing (max 10MB), please re-upload the file")
		return
	}

	var req dto.UpdateWordlistContentRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	// Also check the new content size
	if int64(len(req.Content)) > maxEditableSize {
		dto.BadRequest(c, "Content too large (max 10MB)")
		return
	}

	wordlist, err = h.svc.UpdateContent(id, req.Content)
	if err != nil {
		if errors.Is(err, service.ErrWordlistNotFound) {
			dto.NotFound(c, "Wordlist not found")
			return
		}
		if errors.Is(err, service.ErrFileNotFound) {
			dto.NotFound(c, "Wordlist file not found")
			return
		}
		dto.InternalError(c, "Failed to update wordlist content")
		return
	}

	dto.Success(c, dto.WordlistResponse{
		ID:          wordlist.ID,
		Name:        wordlist.Name,
		Description: wordlist.Description,
		FilePath:    wordlist.FilePath,
		FileSize:    wordlist.FileSize,
		LineCount:   wordlist.LineCount,
		FileHash:    wordlist.FileHash,
		CreatedAt:   wordlist.CreatedAt,
		UpdatedAt:   wordlist.UpdatedAt,
	})
}
