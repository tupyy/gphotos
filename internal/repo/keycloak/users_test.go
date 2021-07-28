package keycloak_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tupyy/gophoto/internal/repo/keycloak"
)

func TestNominal(t *testing.T) {
	keycloakUserRepo, err := keycloak.NewUserRepo(context.Background(), "http://localhost:9000", "admin", "admin", "gophotos")
	assert.Nil(t, err)

	users, err := keycloakUserRepo.Get(context.Background())
	assert.Nil(t, err)
	assert.Len(t, users, 2)
}
