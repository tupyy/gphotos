package keycloak

import (
	"context"
	"errors"

	"github.com/Nerzal/gocloak/v8"
	"github.com/tupyy/gophoto/internal/conf"
)

const (
	AccessToken int = iota
	RefreshToken
)

var (
	ErrRevokeToken = errors.New("failed to revoke token")
)

type KeycloakRepo struct {
	client gocloak.GoCloak
	conf   conf.KeycloakConfig
}

func New(keycloakConf conf.KeycloakConfig) *KeycloakRepo {
	client := gocloak.NewClient(keycloakConf.BaseURL)

	return &KeycloakRepo{
		client: client,
		conf:   keycloakConf,
	}
}

func (k *KeycloakRepo) Logout(accessToken, refreshToken, clientID string) error {
	//token, err := k.connect()
	//if err != nil {
	//	return err
	//}

	ctx, cancel := context.WithTimeout(context.Background(), conf.GetHttpRequestTimeout())
	defer cancel()

	return k.client.LogoutPublicClient(ctx, clientID, k.conf.Realm, "YWRtaW46YWRtaW4=", refreshToken)
}

func (k *KeycloakRepo) connect() (*gocloak.JWT, error) {
	token, err := k.client.LoginClient(context.Background(), "server", "id", k.conf.Realm)
	if err != nil {
		return nil, err
	}

	return token, nil
}
