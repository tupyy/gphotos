package keycloak

import (
	"context"

	"github.com/Nerzal/gocloak/v8"
	keycloak "github.com/Nerzal/gocloak/v8"
	"github.com/pkg/errors"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/utils/logutil"
)

type KeycloakUserRepo struct {
	client keycloak.GoCloak
	token  *keycloak.JWT
	realm  string
}

func NewUserRepo(ctx context.Context, url, username, pwd, realm string) (*KeycloakUserRepo, error) {
	client := gocloak.NewClient(url)
	token, err := client.LoginAdmin(ctx, username, pwd, realm)
	if err != nil {
		return nil, err
	}

	return &KeycloakUserRepo{client: client, token: token, realm: realm}, nil
}

func (u *KeycloakUserRepo) Get(ctx context.Context) ([]entity.User, error) {
	keycloakUsers, err := u.client.GetUsers(ctx, u.token.AccessToken, u.realm, keycloak.GetUsersParams{})
	if err != nil {
		logutil.GetDefaultLogger().WithError(err).Error("cannot fetch users from keycloak")

		return []entity.User{}, errors.Wrap(err, "user repo")
	}

	users := make([]entity.User, 0, len(keycloakUsers))
	for _, keycloakUser := range keycloakUsers {
		if keycloakUser == nil {
			continue
		}

		user := entity.User{
			Username: *keycloakUser.Username,
			UserID:   *keycloakUser.ID,
		}

		users = append(users, user)
	}

	return users, nil
}
