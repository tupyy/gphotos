package group

import (
	"context"
	"errors"

	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/repo"
	"github.com/tupyy/gophoto/models"
	pgclient "github.com/tupyy/gophoto/utils/pgclient"
	"gorm.io/gorm"
)

type groupRepo struct {
	db     *gorm.DB
	client pgclient.Client
}

func NewPostgresRepo(client pgclient.Client) (*groupRepo, error) {
	config := gorm.Config{
		SkipDefaultTransaction: true, // No need transaction for those use cases.
	}

	gormDB, err := client.Open(config)
	if err != nil {
		return &groupRepo{}, err
	}

	return &groupRepo{db: gormDB, client: client}, nil
}

// Create creates the group.
// TODO should the group contain the users and create also the users_groups entries?
func (g *groupRepo) Create(ctx context.Context, group entity.Group) (int32, error) {
	var m models.Groups

	// create group
	m = models.Groups{Name: group.Name}
	if err := g.db.WithContext(ctx).Create(&m).Error; err != nil {
		return -1, err
	}

	return m.ID, nil
}

func (g *groupRepo) Update(ctx context.Context, group entity.Group) error {
	var m models.Groups

	// create group
	m = models.Groups{ID: *group.ID, Name: group.Name}
	if err := g.db.WithContext(ctx).Save(&m).Error; err != nil {
		return err
	}

	return nil
}

func (g *groupRepo) Delete(ctx context.Context, groupID int32) error {
	return repo.ErrNotImplementated
}

func (g *groupRepo) Get(ctx context.Context) ([]entity.Group, error) {
	return []entity.Group{}, repo.ErrNotImplementated
}

func (g *groupRepo) GetByID(ctx context.Context, id int32) (entity.Group, error) {
	return entity.Group{}, repo.ErrNotImplementated
}

func (g *groupRepo) GetByName(ctx context.Context, name string) (entity.Group, error) {
	var m models.Groups

	tx := g.db.WithContext(ctx).Where("name = ?", name).First(&m)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return entity.Group{}, repo.ErrGroupNotFound
		}

		return entity.Group{}, repo.ErrInternalError
	}

	return entity.Group{ID: &m.ID, Name: m.Name}, nil
}

func (g *groupRepo) GetByUserID(ctx context.Context, userID string) ([]entity.Group, error) {
	return []entity.Group{}, repo.ErrNotImplementated
}

func fromModel(m models.Groups) entity.Group {
	return entity.Group{
		ID:   &m.ID,
		Name: m.Name,
	}
}
