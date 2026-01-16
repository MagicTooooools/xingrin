package repository

import (
	"github.com/orbit/server/internal/model"
	"github.com/orbit/server/internal/pkg/scope"
	"gorm.io/gorm"
)

// WordlistRepository handles wordlist database operations
type WordlistRepository struct {
	db *gorm.DB
}

// NewWordlistRepository creates a new wordlist repository
func NewWordlistRepository(db *gorm.DB) *WordlistRepository {
	return &WordlistRepository{db: db}
}

// Create creates a new wordlist
func (r *WordlistRepository) Create(wordlist *model.Wordlist) error {
	return r.db.Create(wordlist).Error
}

// FindByID finds a wordlist by ID
func (r *WordlistRepository) FindByID(id int) (*model.Wordlist, error) {
	var wordlist model.Wordlist
	err := r.db.First(&wordlist, id).Error
	if err != nil {
		return nil, err
	}
	return &wordlist, nil
}

// FindByName finds a wordlist by name
func (r *WordlistRepository) FindByName(name string) (*model.Wordlist, error) {
	var wordlist model.Wordlist
	err := r.db.Where("name = ?", name).First(&wordlist).Error
	if err != nil {
		return nil, err
	}
	return &wordlist, nil
}

// FindAll finds all wordlists with pagination
func (r *WordlistRepository) FindAll(page, pageSize int) ([]model.Wordlist, int64, error) {
	var wordlists []model.Wordlist
	var total int64

	if err := r.db.Model(&model.Wordlist{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Scopes(
		scope.WithPagination(page, pageSize),
		scope.OrderByCreatedAtDesc(),
	).Find(&wordlists).Error

	return wordlists, total, err
}

// Update updates a wordlist
func (r *WordlistRepository) Update(wordlist *model.Wordlist) error {
	return r.db.Save(wordlist).Error
}

// Delete deletes a wordlist
func (r *WordlistRepository) Delete(id int) error {
	return r.db.Delete(&model.Wordlist{}, id).Error
}

// ExistsByName checks if wordlist name exists
func (r *WordlistRepository) ExistsByName(name string, excludeID ...int) (bool, error) {
	var count int64
	query := r.db.Model(&model.Wordlist{}).Where("name = ?", name)
	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}
	err := query.Count(&count).Error
	return count > 0, err
}
