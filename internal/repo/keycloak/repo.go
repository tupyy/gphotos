package keycloak

import (
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v8"
	keycloak "github.com/Nerzal/gocloak/v8"
	"github.com/pkg/errors"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/utils/logutil"
)

const (
	masterRealm = "master"
)

type KeycloakRepo struct {
	client keycloak.GoCloak
	token  *keycloak.JWT
	realm  string
}

func New(ctx context.Context, c conf.KeycloakConfig) (*KeycloakRepo, error) {
	client := gocloak.NewClient(c.BaseURL)
	token, err := client.LoginAdmin(ctx, c.AdminUsername, c.AdminPwd, masterRealm)
	if err != nil {
		return nil, err
	}

	logutil.GetDefaultLogger().Info("keycloak client created")

	return &KeycloakRepo{client: client, token: token, realm: c.Realm}, nil
}

func (k *KeycloakRepo) GetUsers(ctx context.Context) ([]entity.User, error) {
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

	return users, nil
}

func (k *KeycloakRepo) GetUserByID(ctx context.Context, id string) (entity.User, error) {
	return entity.User{}, errors.New("no implementatedrlbu")
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

func mapper(u keycloak.User) entity.User {
	user := entity.User{
		Username: *u.Username,
		ID:       *u.ID,
	}

	if u.FirstName != nil {
		user.FirstName = *u.FirstName
	}

	if u.LastName != nil {
		user.LastName = *u.LastName
	}

	if u.Attributes != nil {
		m := *u.Attributes
		if attrs, found := m["can_share"]; found {
			for _, attr := range attrs {
				switch attr {
				case "true":
					user.CanShare = true
				default:
					user.CanShare = false
				}
			}
		}
	}

	return user
}

func ptrBool(b bool) *bool {
	return &b
}
