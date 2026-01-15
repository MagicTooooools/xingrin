package service

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/xingrin/server/internal/dto"
	"github.com/xingrin/server/internal/model"
	"github.com/xingrin/server/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrWordlistNotFound = errors.New("wordlist not found")
	ErrWordlistExists   = errors.New("wordlist name already exists")
	ErrEmptyName        = errors.New("wordlist name cannot be empty")
	ErrFileNotFound     = errors.New("wordlist file not found")
)

// WordlistService handles wordlist business logic
type WordlistService struct {
	repo    *repository.WordlistRepository
	basePath string
}

// NewWordlistService creates a new wordlist service
func NewWordlistService(repo *repository.WordlistRepository, basePath string) *WordlistService {
	return &WordlistService{
		repo:     repo,
		basePath: basePath,
	}
}

// Create creates a new wordlist with file upload
func (s *WordlistService) Create(name, description, filename string, fileContent io.Reader) (*model.Wordlist, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrEmptyName
	}

	exists, err := s.repo.ExistsByName(name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrWordlistExists
	}

	// Ensure storage directory exists
	if err := os.MkdirAll(s.basePath, 0755); err != nil {
		return nil, err
	}

	// Sanitize filename
	safeFilename := sanitizeFilename(filename)
	fullPath := filepath.Join(s.basePath, safeFilename)

	// Write file and calculate hash simultaneously
	file, err := os.Create(fullPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	hasher := sha256.New()
	writer := io.MultiWriter(file, hasher)

	written, err := io.Copy(writer, fileContent)
	if err != nil {
		os.Remove(fullPath) // Cleanup on error
		return nil, err
	}

	fileHash := hex.EncodeToString(hasher.Sum(nil))

	// Count lines
	lineCount, err := countLines(fullPath)
	if err != nil {
		lineCount = 0 // Non-fatal error
	}

	wordlist := &model.Wordlist{
		Name:        name,
		Description: description,
		FilePath:    fullPath,
		FileSize:    written,
		LineCount:   lineCount,
		FileHash:    fileHash,
	}

	if err := s.repo.Create(wordlist); err != nil {
		os.Remove(fullPath) // Cleanup on error
		return nil, err
	}

	return wordlist, nil
}

// List returns paginated wordlists
func (s *WordlistService) List(query *dto.PaginationQuery) ([]model.Wordlist, int64, error) {
	return s.repo.FindAll(query.GetPage(), query.GetPageSize())
}

// GetByID returns a wordlist by ID
func (s *WordlistService) GetByID(id int) (*model.Wordlist, error) {
	wordlist, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWordlistNotFound
		}
		return nil, err
	}
	return wordlist, nil
}

// GetByName returns a wordlist by name
func (s *WordlistService) GetByName(name string) (*model.Wordlist, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrWordlistNotFound
	}

	wordlist, err := s.repo.FindByName(name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWordlistNotFound
		}
		return nil, err
	}
	return wordlist, nil
}

// Delete deletes a wordlist and its file
func (s *WordlistService) Delete(id int) error {
	wordlist, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrWordlistNotFound
		}
		return err
	}

	// Delete file (best effort)
	if wordlist.FilePath != "" {
		os.Remove(wordlist.FilePath)
	}

	return s.repo.Delete(id)
}

// GetContent returns the content of a wordlist file
func (s *WordlistService) GetContent(id int) (string, error) {
	wordlist, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrWordlistNotFound
		}
		return "", err
	}

	if wordlist.FilePath == "" {
		return "", ErrFileNotFound
	}

	content, err := os.ReadFile(wordlist.FilePath)
	if err != nil {
		return "", ErrFileNotFound
	}

	return string(content), nil
}

// UpdateContent updates the content of a wordlist file
func (s *WordlistService) UpdateContent(id int, content string) (*model.Wordlist, error) {
	wordlist, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWordlistNotFound
		}
		return nil, err
	}

	if wordlist.FilePath == "" {
		return nil, ErrFileNotFound
	}

	// Write new content
	if err := os.WriteFile(wordlist.FilePath, []byte(content), 0644); err != nil {
		return nil, err
	}

	// Recalculate stats
	fileInfo, err := os.Stat(wordlist.FilePath)
	if err != nil {
		return nil, err
	}

	// Calculate hash
	hasher := sha256.New()
	hasher.Write([]byte(content))
	fileHash := hex.EncodeToString(hasher.Sum(nil))

	// Count lines
	lineCount := strings.Count(content, "\n")
	if len(content) > 0 && !strings.HasSuffix(content, "\n") {
		lineCount++
	}

	// Update record
	wordlist.FileSize = fileInfo.Size()
	wordlist.LineCount = lineCount
	wordlist.FileHash = fileHash

	if err := s.repo.Update(wordlist); err != nil {
		return nil, err
	}

	return wordlist, nil
}

// GetFilePath returns the file path of a wordlist by name (for download)
func (s *WordlistService) GetFilePath(name string) (string, error) {
	wordlist, err := s.GetByName(name)
	if err != nil {
		return "", err
	}

	if wordlist.FilePath == "" || !fileExists(wordlist.FilePath) {
		return "", ErrFileNotFound
	}

	return wordlist.FilePath, nil
}

// Helper functions

func sanitizeFilename(filename string) string {
	// Remove path separators to prevent directory traversal
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, "\\", "_")

	if filename == "" {
		filename = "wordlist.txt"
	}

	// Add .txt extension if missing
	if filepath.Ext(filename) == "" {
		filename += ".txt"
	}

	return filename
}

func countLines(filepath string) (int, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		count++
	}

	return count, scanner.Err()
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
