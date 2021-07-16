package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/repo"
	"github.com/tupyy/gophoto/models"
	"github.com/tupyy/gophoto/utils/logutil"
	pgclient "github.com/tupyy/gophoto/utils/pgclient"
	"gorm.io/gorm"
)

type UserRepo struct {
	db     *gorm.DB
	client pgclient.Client
}

func New(client pgclient.Client) (*UserRepo, error) {
	config := gorm.Config{
		SkipDefaultTransaction: true, // No need transaction for those use cases.
	}

	gormDB, err := client.Open(config)
	if err != nil {
		return &UserRepo{}, err
	}

	return &UserRepo{db: gormDB, client: client}, nil
}

func (u *UserRepo) Create(ctx context.Context, user entity.User) (int, error) {
	tx := u.db.WithContext(ctx).Begin()

	m := fromUserEntity(user)

	tx.Create(&m)
	if tx.Error != nil {
		return 0, tx.Error
	}

	user.ID = &m.ID

	// create group relationships
	if len(user.Groups) > 0 {
		var userGroups = make([]models.UsersGroups, 0, len(user.Groups))

		for _, group := range user.Groups {
			ug := models.UsersGroups{
				UsersID:  *user.ID,
				GroupsID: *group.ID,
			}

			userGroups = append(userGroups, ug)
		}

		tx.CreateInBatches(userGroups, len(userGroups))
		if tx.Error != nil {
			tx.Rollback()
			return 0, tx.Error
		}
	}

	if err := tx.Commit().Error; err != nil {
		return 0, tx.Error
	}

	return int(m.ID), nil
}

// Update updates an user.
func (u *UserRepo) Update(ctx context.Context, user entity.User) (entity.User, error) {
	err := user.Validate()
	if err != nil {
		return user, err
	}

	tx := u.db.WithContext(ctx).Begin()

	// get the groups if any
	var userGroups []models.UsersGroups

	tx.Table("users_groups").
		Joins("INNER JOIN users on users.id = users_groups.users_id").
		Where("users.id = ?", user.ID).
		Find(&userGroups)

	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		return user, tx.Error
	}

	for _, ug := range userGroups {
		found := false
		for _, g := range user.Groups {
			if *g.ID == ug.GroupsID {
				found = true
				break
			}
		}

		if !found {
			tx.Exec("DELETE from users_groups WHERE users_id = ? AND groups_id = ?;", ug.UsersID, ug.GroupsID)
			if tx.Error != nil {
				tx.Rollback()
				return user, tx.Error
			}
		}
	}

	// add the new ones
	for _, g := range user.Groups {
		found := false
		for _, ug := range userGroups {
			if ug.GroupsID == *g.ID {
				found = true
				break
			}
		}

		if !found {
			tx.Create(models.UsersGroups{UsersID: *user.ID, GroupsID: *g.ID})
			if tx.Error != nil {
				tx.Rollback()
				return user, tx.Error
			}
		}
	}

	tx.Save(fromUserEntity(user))
	if tx.Error != nil {
		tx.Rollback()
		return user, tx.Error
	}

	if err := tx.Commit().Error; err != nil {
		return user, err
	}

	return user, nil

}

func (u *UserRepo) Get(ctx context.Context, username string) (entity.User, error) {
	var m models.Users
	var emptyUser = entity.User{}

	tx := u.db.WithContext(ctx).Where("username = ?", username).First(&m)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return emptyUser, repo.ErrUserNotFound
		}
		return emptyUser, tx.Error
	}

	user := toUserEntity(m)

	// get the groups if any
	var userGroups []models.Groups

	tx = u.db.WithContext(ctx).Table("groups").
		Select("groups.id, groups.name").
		Joins("LEFT JOIN users_groups ON groups.id = users_groups.groups_id").
		Joins("LEFT JOIN users on users.id = users_groups.users_id").
		Where("users.id = ?", user.ID).
		Find(&userGroups)

	if tx.Error != nil {
		return emptyUser, tx.Error
	}

	for _, ug := range userGroups {
		user.Groups = append(user.Groups, toGroupEntity(ug))
	}

	if valErr := user.Validate(); valErr != nil {
		return emptyUser, fmt.Errorf("%w validation error: %v", entity.ErrInvalidEntity, valErr)
	}

	return user, nil

}

func (u *UserRepo) GetUsers(ctx context.Context) ([]entity.User, error) {
	var results []customUser
	var users []entity.User

	tx := u.db.WithContext(ctx).Table("users").
		Select("users.*, STRING_AGG(groups.name, ',') as group_names, ARRAY_REMOVE(ARRAY_AGG(groups.id),NULL) as group_ids").
		Joins("INNER JOIN users_groups ON ( users_groups.users_id = users.id )").
		Joins("INNER JOIN groups ON (users_groups.groups_id = groups.id)").
		Group("users.id").
		Find(&results)

	if tx.Error != nil {
		logutil.GetDefaultLogger().WithError(tx.Error).Error("fetching users")

		return users, tx.Error
	}

	if len(results) == 0 {
		logutil.GetDefaultLogger().Warn("no user found")

		return users, repo.ErrUserNotFound
	}

	for _, m := range results {
		user := m.toUser()

		users = append(users, user)
	}

	return users, nil

}
func toUserEntity(m models.Users) entity.User {
	var r entity.Role
	switch m.Role {
	case "admin":
		r = entity.RoleAdmin
	case "editor":
		r = entity.RoleEditor
	case "user":
		r = entity.RoleUser
	default:
		r = entity.RoleUser
	}

	return entity.User{
		ID:       &m.ID,
		Username: m.Username,
		CanShare: *m.CanShare,
		UserID:   m.UserID,
		Role:     r,
	}
}

type customUser struct {
	ID         int32         `gorm:"primary_key;column:id;type:INT4;"`
	Username   string        `gorm:"column:username;type:TEXT;"`
	Role       models.Role   `gorm:"column:role;type:ROLE;"`
	UserID     string        `gorm:"column:user_id;type:TEXT;"`
	CanShare   *bool         `gorm:"column:can_share;type:BOOL;"`
	GroupNames string        `gorm:"column:group_names;type:TEXT;"`
	GroupIDS   pq.Int64Array `gorm:"column:group_ids;type:_INT4;"`
}

func (m customUser) toUser() entity.User {
	var r entity.Role
	switch m.Role {
	case "admin":
		r = entity.RoleAdmin
	case "editor":
		r = entity.RoleEditor
	case "user":
		r = entity.RoleUser
	default:
		r = entity.RoleUser
	}

	user := entity.User{
		ID:       &m.ID,
		Username: m.Username,
		CanShare: *m.CanShare,
		UserID:   m.UserID,
		Role:     r,
	}

	if m.GroupNames == "" {
		return user
	}

	groupNames := strings.Split(m.GroupNames, ",")
	for idx, name := range groupNames {
		id := int32(m.GroupIDS[idx])

		user.Groups = append(user.Groups, entity.Group{ID: &id, Name: name})
	}

	return user

}
func fromUserEntity(u entity.User) models.Users {
	var r models.Role

	switch u.Role {
	case entity.RoleAdmin:
		r = "admin"
	case entity.RoleUser:
		r = "user"
	case entity.RoleEditor:
		r = "editor"
	}

	if u.ID != nil {
		return models.Users{
			ID:       *u.ID,
			Username: u.Username,
			Role:     r,
			CanShare: &u.CanShare,
			UserID:   u.UserID,
		}
	}

	return models.Users{
		Username: u.Username,
		Role:     r,
		CanShare: &u.CanShare,
		UserID:   u.UserID,
	}

}

func toGroupEntity(m models.Groups) entity.Group {
	return entity.Group{
		ID:   &m.ID,
		Name: m.Name,
	}
}
