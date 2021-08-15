package album_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	ac "github.com/tupyy/gophoto/internal/controllers"
	"github.com/tupyy/gophoto/internal/domain/entity"
)

func TestPermissionResolver(t *testing.T) {
	album := entity.Album{
		OwnerID: "batman",
		UserPermissions: entity.Permissions{
			"batman": {entity.PermissionReadAlbum, entity.PermissionEditAlbum},
		},
		GroupPermissions: entity.Permissions{
			"good_guys": {entity.PermissionDeleteAlbum},
		},
	}

	user := entity.User{
		ID:     "batman",
		Role:   entity.RoleEditor,
		Groups: []entity.Group{{Name: "good_guys"}, {Name: "admin"}},
	}

	apr := ac.NewAlbumPermissionResolver()

	hasPermission := apr.Policy(ac.RolePolicy{entity.RoleEditor}).Strategy(ac.AtLeastOneStrategy).Resolve(album, user)
	assert.True(t, hasPermission)

	apr = ac.NewAlbumPermissionResolver()
	hasPermission = apr.Policy(ac.RolePolicy{entity.RoleAdmin}).Strategy(ac.AtLeastOneStrategy).Resolve(album, user)
	assert.False(t, hasPermission)

	apr = ac.NewAlbumPermissionResolver()
	hasPermission = apr.Policy(ac.RolePolicy{entity.RoleAdmin}).
		Policy(ac.RolePolicy{entity.RoleEditor}).
		Strategy(ac.UnanimousStrategy).Resolve(album, user)
	assert.False(t, hasPermission)

	apr = ac.NewAlbumPermissionResolver()
	hasPermission = apr.Policy(ac.UserPermissionPolicy{entity.PermissionEditAlbum}).
		Policy(ac.UserPermissionPolicy{entity.PermissionDeleteAlbum}).
		Strategy(ac.AtLeastOneStrategy).
		Resolve(album, user)
	assert.True(t, hasPermission)

	apr = ac.NewAlbumPermissionResolver()
	hasPermission = apr.Policy(ac.GroupPermissionPolicy{entity.PermissionEditAlbum}).
		Policy(ac.UserPermissionPolicy{entity.PermissionDeleteAlbum}).
		Strategy(ac.AtLeastOneStrategy).
		Resolve(album, user)
	assert.False(t, hasPermission)
}
