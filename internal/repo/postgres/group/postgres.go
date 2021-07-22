package group

import (
	"context"
	"errors"
	"fmt"

	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/repo"
	"github.com/tupyy/gophoto/models"
	pgclient "github.com/tupyy/gophoto/utils/pgclient"
	"gorm.io/gorm"
)

type GroupRepo struct {
	db     *gorm.DB
	client pgclient.Client
}

func NewPostgresRepo(client pgclient.Client) (*GroupRepo, error) {
	config := gorm.Config{
		SkipDefaultTransaction: true, // No need transaction for those use cases.
	}

	gormDB, err := client.Open(config)
	if err != nil {
		return &GroupRepo{}, err
	}

	return &GroupRepo{db: gormDB, client: client}, nil
}

// Create creates the group.
// TODO should the group contain the users and create also the users_groups entries?
func (g *GroupRepo) Create(ctx context.Context, group entity.Group) (int32, error) {
	var m models.Groups

	// create group
	m = models.Groups{Name: group.Name}
	if err := g.db.WithContext(ctx).Create(&m).Error; err != nil {
		return -1, err
	}

	return m.ID, nil
}

func (g *GroupRepo) Update(ctx context.Context, group entity.Group) error {
	var m models.Groups

	// create group
	m = models.Groups{ID: *group.ID, Name: group.Name}
	if err := g.db.WithContext(ctx).Save(&m).Error; err != nil {
		return err
	}

	return nil
}

func (g *GroupRepo) Delete(ctx context.Context, groupID int32) error {
	return repo.ErrNotImplementated
}

func (g *GroupRepo) Get(ctx context.Context) ([]entity.Group, error) {
	var groups []models.Groups

	tx := g.db.WithContext(ctx).Find(&groups)
	if tx.Error != nil {
		return []entity.Group{}, fmt.Errorf("%w %v", repo.ErrInternalError, tx.Error)
	}

	entities := make([]entity.Group, 0, len(groups))

	for _, g := range groups {
		entities = append(entities, entity.Group{ID: &g.ID, Name: g.Name})
	}

	return entities, nil
}

func (g *GroupRepo) GetByID(ctx context.Context, id int32) (entity.Group, error) {
	return entity.Group{}, repo.ErrNotImplementated
}

func (g *GroupRepo) GetByName(ctx context.Context, name string) (entity.Group, error) {
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

func (g *GroupRepo) GetByUserID(ctx context.Context, userID string) ([]entity.Group, error) {
	return []entity.Group{}, repo.ErrNotImplementated
}

func fromModel(m models.Groups) entity.Group {
	return entity.Group{
		ID:   &m.ID,
		Name: m.Name,
	}
}
