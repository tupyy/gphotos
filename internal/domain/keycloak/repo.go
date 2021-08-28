package keycloak

import (
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v8"
	keycloak "github.com/Nerzal/gocloak/v8"
	"github.com/pkg/errors"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/domain/entity"
	userFilters "github.com/tupyy/gophoto/internal/domain/filters/user"
	"github.com/tupyy/gophoto/utils/logutil"
)

const (
	masterRealm = "master"
)

type KeycloakRepo struct {
	client        keycloak.GoCloak
	token         *keycloak.JWT
	realm         string
	configuration conf.KeycloakConfig
}

func New(ctx context.Context, c conf.KeycloakConfig) (*KeycloakRepo, error) {
	client := gocloak.NewClient(c.BaseURL)
	token, err := client.LoginAdmin(ctx, c.AdminUsername, c.AdminPwd, masterRealm)
	if err != nil {
		return nil, err
	}

	logutil.GetDefaultLogger().Infof("connected to keycloak as %s", c.AdminUsername)

	return &KeycloakRepo{client: client, token: token, realm: c.Realm, configuration: c}, nil
}

func (k *KeycloakRepo) GetUsers(ctx context.Context, filters ...userFilters.Filter) ([]entity.User, error) {
	keycloakUsers, err := k.client.GetUsers(ctx, k.token.AccessToken, k.realm, keycloak.GetUsersParams{Enabled: ptrBool(true)})
	if err != nil {
		logutil.GetDefaultLogger().WithError(err).Error("cannot fetch users from keycloak")

		return []entity.User{}, errors.Wrap(err, "user repo")
	}

	users := make([]entity.User, 0, len(keycloakUsers))
	for _, keycloakUser := range keycloakUsers {
		if keycloakUser == nil {
			continue
		}

		logutil.GetDefaultLogger().WithField("user", fmt.Sprintf("%+v", keycloakUser)).Trace("found user")

		users = append(users, mapper(*keycloakUser))
	}

	// get groups
	for _, user := range users {
		groups, err := k.client.GetUserGroups(ctx, k.token.AccessToken, k.realm, user.ID, keycloak.GetGroupsParams{})
		if err != nil {
			return []entity.User{}, err
		}

		user.Groups = make([]entity.Group, 0, len(groups))
		for _, g := range groups {
			user.Groups = append(user.Groups, entity.Group{Name: *g.Name})
		}
	}

	//filter them
	if len(filters) > 0 {
		users = filter(filters, users)
		logutil.GetDefaultLogger().WithField("count filtered users", len(users)).Debug("filter user")
	}

	return users, nil
}

func (k *KeycloakRepo) GetUserByID(ctx context.Context, id string) (entity.User, error) {
	keycloakUser, err := k.client.GetUserByID(ctx, k.token.AccessToken, k.realm, id)
	if err != nil {
		logutil.GetDefaultLogger().WithError(err).WithField("id", id).Error("cannot fetch user from keycloak")

		return entity.User{}, errors.Wrap(err, "user repo")
	}

	if !*keycloakUser.Enabled {
		logutil.GetDefaultLogger().WithError(err).WithField("id", id).Error("user disabled")

		return entity.User{}, fmt.Errorf("user %s disabled", id)
	}

	return mapper(*keycloakUser), nil
}

func (k *KeycloakRepo) GetGroups(ctx context.Context) ([]entity.Group, error) {
	kgroups, err := k.client.GetGroups(ctx, k.token.AccessToken, k.realm, keycloak.GetGroupsParams{})
	if err != nil {
		return []entity.Group{}, err
	}

	groups := make([]entity.Group, 0, len(kgroups))
	for _, g := range kgroups {
		groups = append(groups, entity.Group{Name: *g.Name})
	}

	return groups, nil
}
