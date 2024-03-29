package user

import (
	"context"
	"fmt"
	"strings"

	pgclient "github.com/tupyy/gophoto/internal/clients/pg"
	"github.com/tupyy/gophoto/internal/entity"
	domain "github.com/tupyy/gophoto/internal/repos"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserPostgresRepo struct {
	db     *gorm.DB
	client pgclient.Client
}

func NewPostgresRepo(client pgclient.Client) (*UserPostgresRepo, error) {
	config := gorm.Config{
		SkipDefaultTransaction: true, // No need transaction for those use cases.
	}

	gormDB, err := client.Open(config)
	if err != nil {
		return &UserPostgresRepo{}, err
	}

	return &UserPostgresRepo{gormDB, client}, nil
}

// GetByGroups returns a list of albums with at least one persmission for at least on group in the list.
// It does not sort or filter the album here. The sorting and filter is done at cache level.
func (a *UserPostgresRepo) GetRelatedUsers(ctx context.Context, user entity.User) (ids []string, err error) {
	var results []result

	where := fmt.Sprintf("album_user_permissions.user_id = '%s'", user.ID)

	if len(user.Groups) > 0 {
		var groups strings.Builder
		for idx, g := range user.Groups {
			groups.WriteString(fmt.Sprintf("'%s'", g.Name))

			if idx < len(user.Groups)-1 {
				groups.WriteString(",")
			}
		}
		where = fmt.Sprintf("%s OR album_group_permissions.group_name = ANY(ARRAY[%s])", where, groups.String())
	}

	tx := a.db.WithContext(ctx).Table("album").
		Select("album.id, album.owner_id").
		Joins("LEFT JOIN album_user_permissions ON (album.id = album_user_permissions.album_id)").
		Joins("LEFT JOIN album_group_permissions ON (album.id = album_group_permissions.album_id)").
		Where(where).
		Find(&results)
	if tx.Error != nil {
		return []string{}, fmt.Errorf("%w internal error: %v", domain.ErrInternalError, tx.Error)
	}

	if len(results) == 0 {
		zap.S().Warnw("user has no relationships to other users", "user_id", user.Username)
		return []string{}, nil
	}

	return mapper(results), nil
}

type result struct {
	ID      int32  `gorm:"column_name:id;type:INT4"`
	OwnerID string `gorm:"column:owner_id;type:INT4;"`
}

func mapper(results []result) []string {
	ids := make([]string, 0, len(results))

	for _, r := range results {
		ids = append(ids, r.OwnerID)
	}

	return ids
}
