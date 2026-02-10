package application

import (
	"context"

	catalogdomain "github.com/yyhuni/lunafox/server/internal/modules/catalog/domain"
)

var (
	ErrEngineNotFound = catalogdomain.ErrEngineNotFound
	ErrEngineExists   = catalogdomain.ErrEngineExists
	ErrInvalidEngine  = catalogdomain.ErrInvalidEngine
)

type EngineCommandStore interface {
	GetByID(id int) (*catalogdomain.ScanEngine, error)
	ExistsByName(name string, excludeID ...int) (bool, error)
	Create(engine *catalogdomain.ScanEngine) error
	Update(engine *catalogdomain.ScanEngine) error
	Delete(id int) error
}

type EngineCommandService struct {
	store EngineCommandStore
}

func NewEngineCommandService(store EngineCommandStore) *EngineCommandService {
	return &EngineCommandService{store: store}
}

func (service *EngineCommandService) CreateEngine(ctx context.Context, name, configuration string) (*catalogdomain.ScanEngine, error) {
	_ = ctx

	engine, err := catalogdomain.NewScanEngine(name, configuration)
	if err != nil {
		return nil, err
	}

	exists, err := service.store.ExistsByName(engine.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEngineExists
	}

	if err := service.store.Create(engine); err != nil {
		return nil, err
	}

	return engine, nil
}

func (service *EngineCommandService) UpdateEngine(ctx context.Context, id int, name, configuration string) (*catalogdomain.ScanEngine, error) {
	_ = ctx

	engine, err := service.store.GetByID(id)
	if err != nil {
		return nil, err
	}

	previousName := engine.Name
	if err := engine.Rename(name); err != nil {
		return nil, err
	}
	if engine.Name != previousName {
		exists, err := service.store.ExistsByName(engine.Name, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrEngineExists
		}
	}

	engine.Reconfigure(configuration)

	if err := service.store.Update(engine); err != nil {
		return nil, err
	}

	return engine, nil
}

func (service *EngineCommandService) PatchEngine(ctx context.Context, id int, name, configuration *string) (*catalogdomain.ScanEngine, error) {
	_ = ctx

	engine, err := service.store.GetByID(id)
	if err != nil {
		return nil, err
	}

	if name != nil {
		previousName := engine.Name
		if err := engine.Rename(*name); err != nil {
			return nil, err
		}

		if engine.Name != previousName {
			exists, err := service.store.ExistsByName(engine.Name, id)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, ErrEngineExists
			}
		}
	}

	if configuration != nil {
		engine.Reconfigure(*configuration)
	}

	if err := service.store.Update(engine); err != nil {
		return nil, err
	}

	return engine, nil
}

func (service *EngineCommandService) DeleteEngine(ctx context.Context, id int) error {
	_ = ctx

	if _, err := service.store.GetByID(id); err != nil {
		return err
	}
	return service.store.Delete(id)
}
