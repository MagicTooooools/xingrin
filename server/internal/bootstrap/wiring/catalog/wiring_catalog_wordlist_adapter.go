package catalogwiring

import (
	catalogdomain "github.com/yyhuni/lunafox/server/internal/modules/catalog/domain"
	catalogrepo "github.com/yyhuni/lunafox/server/internal/modules/catalog/repository"
)

type catalogWordlistStoreAdapter struct {
	repo *catalogrepo.WordlistRepository
}

func newCatalogWordlistStoreAdapter(repo *catalogrepo.WordlistRepository) *catalogWordlistStoreAdapter {
	return &catalogWordlistStoreAdapter{repo: repo}
}

func (adapter *catalogWordlistStoreAdapter) FindAll(page, pageSize int) ([]catalogdomain.Wordlist, int64, error) {
	wordlists, total, err := adapter.repo.FindAll(page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return catalogModelWordlistListToDomain(wordlists), total, nil
}

func (adapter *catalogWordlistStoreAdapter) List() ([]catalogdomain.Wordlist, error) {
	wordlists, err := adapter.repo.List()
	if err != nil {
		return nil, err
	}
	return catalogModelWordlistListToDomain(wordlists), nil
}

func (adapter *catalogWordlistStoreAdapter) FindByID(id int) (*catalogdomain.Wordlist, error) {
	wordlist, err := adapter.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return catalogModelWordlistToDomain(wordlist), nil
}

func (adapter *catalogWordlistStoreAdapter) FindByName(name string) (*catalogdomain.Wordlist, error) {
	wordlist, err := adapter.repo.FindByName(name)
	if err != nil {
		return nil, err
	}
	return catalogModelWordlistToDomain(wordlist), nil
}

func (adapter *catalogWordlistStoreAdapter) ExistsByName(name string) (bool, error) {
	return adapter.repo.ExistsByName(name)
}

func (adapter *catalogWordlistStoreAdapter) Create(wordlist *catalogdomain.Wordlist) error {
	modelWordlist := catalogDomainWordlistToModel(wordlist)
	if err := adapter.repo.Create(modelWordlist); err != nil {
		return err
	}
	*wordlist = *catalogModelWordlistToDomain(modelWordlist)
	return nil
}

func (adapter *catalogWordlistStoreAdapter) Update(wordlist *catalogdomain.Wordlist) error {
	return adapter.repo.Update(catalogDomainWordlistToModel(wordlist))
}

func (adapter *catalogWordlistStoreAdapter) Delete(id int) error {
	return adapter.repo.Delete(id)
}
