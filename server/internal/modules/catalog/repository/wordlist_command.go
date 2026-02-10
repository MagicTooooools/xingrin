package repository

import (
	"github.com/yyhuni/lunafox/server/internal/modules/catalog/repository/persistence"
)

// Create creates a new wordlist
func (r *WordlistRepository) Create(wordlist *model.Wordlist) error {
	return r.db.Create(wordlist).Error
}

// Update updates a wordlist
func (r *WordlistRepository) Update(wordlist *model.Wordlist) error {
	return r.db.Save(wordlist).Error
}

// Delete deletes a wordlist by ID
func (r *WordlistRepository) Delete(id int) error {
	return r.db.Delete(&model.Wordlist{}, id).Error
}
