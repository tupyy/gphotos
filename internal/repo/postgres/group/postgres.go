package group

import (
	"context"
	"fmt"

	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/models"
	pgclient "github.com/tupyy/gophoto/utils/pgclient"
	"gorm.io/gorm"
)

type PostgresGroupRepo struct {
	db     *gorm.DB
	client pgclient.Client
}

func New(client pgclient.Client) (*PostgresGroupRepo, error) {
	config := gorm.Config{
		SkipDefaultTransaction: true, // No need transaction for those use cases.
	}

	gormDB, err := client.Open(config)
	if err != nil {
		return &PostgresGroupRepo{}, err
	}

	return &PostgresGroupRepo{db: gormDB, client: client}, nil
}

func (g *PostgresGroupRepo) FirstOrCreate(ctx context.Context, name string) (entity.Group, bool, error) {
	var m models.Groups
	var emptyGroup entity.Group
	var created bool

	tx := g.db.WithContext(ctx).Where("name = ?", name).First(&m)
	if tx.Error != nil {
		if tx.Error != gorm.ErrRecordNotFound {
			return emptyGroup, false, tx.Error
		}

		// create group
		m = models.Groups{Name: name}
		if err := g.db.WithContext(ctx).Create(&m).Error; err != nil {
			return emptyGroup, false, err
		}

		created = true
	}

	ent := fromModel(m)
	if valErr := ent.Validate(); valErr != nil {
		return emptyGroup, false, fmt.Errorf("%w group validation error: %v", entity.ErrInvalidEntity, valErr)
	}

	return ent, created, nil

}

func fromModel(m models.Groups) entity.Group {
	return entity.Group{
		ID:   &m.ID,
		Name: m.Name,
	}
}
