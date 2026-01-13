package handler

import (
	"errors"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xingrin/go-backend/internal/dto"
	"github.com/xingrin/go-backend/internal/service"
)

// WordlistHandler handles wordlist endpoints
type WordlistHandler struct {
	svc *service.WordlistService
}

// NewWordlistHandler creates a new wordlist handler
func NewWordlistHandler(svc *service.WordlistService) *WordlistHandler {
	return &WordlistHandler{svc: svc}
}

// Create uploads and creates a new wordlist
// POST /api/wordlists
func (h *WordlistHandler) Create(c *gin.Context) {
	name := c.PostForm("name")
	description := c.PostForm("description")

	file, err := c.FormFile("file")
	if err != nil {
		dto.BadRequest(c, "Missing wordlist file")
		return
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		dto.InternalError(c, "Failed to read uploaded file")
		return
	}
	defer src.Close()

	wordlist, err := h.svc.Create(name, description, file.Filename, src)
	if err != nil {
		if errors.Is(err, service.ErrEmptyName) {
			dto.BadRequest(c, "Wordlist name cannot be empty")
			return
		}
		if errors.Is(err, service.ErrWordlistExists) {
			dto.BadRequest(c, "Wordlist name already exists")
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

// List returns paginated wordlists
// GET /api/wordlists
func (h *WordlistHandler) List(c *gin.Context) {
	var query dto.PaginationQuery
	if !dto.BindQuery(c, &query) {
		return
	}

	wordlists, total, err := h.svc.List(&query)
	if err != nil {
		dto.InternalError(c, "Failed to list wordlists")
		return
	}

	var resp []dto.WordlistResponse
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

	dto.Paginated(c, resp, total, query.GetPage(), query.GetPageSize())
}

// Delete deletes a wordlist
// DELETE /api/wordlists/:id
func (h *WordlistHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid wordlist ID")
		return
	}

	err = h.svc.Delete(id)
	if err != nil {
		if errors.Is(err, service.ErrWordlistNotFound) {
			dto.NotFound(c, "Wordlist not found")
			return
		}
		dto.InternalError(c, "Failed to delete wordlist")
		return
	}

	dto.NoContent(c)
}

// Download downloads a wordlist file by name
// GET /api/wordlists/download?wordlist=xxx
func (h *WordlistHandler) Download(c *gin.Context) {
	name := c.Query("wordlist")
	if name == "" {
		dto.BadRequest(c, "Missing parameter: wordlist")
		return
	}

	filePath, err := h.svc.GetFilePath(name)
	if err != nil {
		if errors.Is(err, service.ErrWordlistNotFound) || errors.Is(err, service.ErrFileNotFound) {
			dto.NotFound(c, "Wordlist not found")
			return
		}
		dto.InternalError(c, "Failed to get wordlist")
		return
	}

	c.FileAttachment(filePath, filepath.Base(filePath))
}

// GetContent returns the content of a wordlist file
// GET /api/wordlists/:id/content
func (h *WordlistHandler) GetContent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid wordlist ID")
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

	var req dto.UpdateWordlistContentRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	wordlist, err := h.svc.UpdateContent(id, req.Content)
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
