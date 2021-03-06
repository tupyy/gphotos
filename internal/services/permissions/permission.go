package permissions

import (
	"github.com/tupyy/gophoto/internal/entity"
)

type StrategyType int

const (
	// UnanimousStrategy means that all policies must be true.
	UnanimousStrategy StrategyType = iota
	// AtLeastOneStrategy means that at least one policies must be true.
	AtLeastOneStrategy
)

type Policy interface {
	Resolve(entity.Album, entity.User) bool
}

// albumPermissionService resolve a set of conditions set on an album against an user.
// For example in case of editing an album at least one of three conditions must met:
//    - user is the owner of the album
//	  - user has edit permission set directly to him by the owner
//    - the user's group has edit permission set by the owner
// To resolve this case the Album
type albumPermissionService struct {
	policies []Policy
	strategy StrategyType
}

// Create a new albumPermissionResolver with AtLeastOneStrategy by default.
func NewAlbumPermissionService() *albumPermissionService {
	return &albumPermissionService{
		policies: make([]Policy, 0, 3), // often we have 3 policies
		strategy: AtLeastOneStrategy,
	}
}

func (apr *albumPermissionService) Policy(p Policy) *albumPermissionService {
	apr.policies = append(apr.policies, p)

	return apr
}

func (apr *albumPermissionService) Strategy(s StrategyType) *albumPermissionService {
	apr.strategy = s

	return apr
}

func (apr *albumPermissionService) Resolve(album entity.Album, user entity.User) bool {
	switch apr.strategy {
	case AtLeastOneStrategy:
		for _, policy := range apr.policies {
			if policy.Resolve(album, user) {
				return true
			}
		}

		return false
	case UnanimousStrategy:
		result, first := false, true

		for _, policy := range apr.policies {
			resultPolicy := policy.Resolve(album, user)
			if first {
				result = resultPolicy
				first = false
			} else {
				result = result && resultPolicy
			}
		}

		return result
	default:
		return false
	}
}

// OwnerPolicy checks if the user is the owner of the album.
type OwnerPolicy struct{}

func (i OwnerPolicy) Resolve(a entity.Album, u entity.User) bool {
	return a.OwnerID == u.ID
}

// RolePolicy checks if the user has a certain role.
type RolePolicy struct {
	Role entity.Role
}

func (r RolePolicy) Resolve(a entity.Album, u entity.User) bool {
	return r.Role == u.Role
}

// UserPermissionPolicy checks if the album gives the user the permission.
type UserPermissionPolicy struct {
	Permission entity.Permission
}

func (up UserPermissionPolicy) Resolve(a entity.Album, u entity.User) bool {
	return entity.HasUserPermission(a, u.ID, up.Permission)
}

type AnyUserPermissionPolicty struct{}

func (ap AnyUserPermissionPolicty) Resolve(a entity.Album, u entity.User) bool {
	return entity.HasUserPermissions(a, u.ID)
}

// GroupPermissionPolicy checks if one of the users's group has the permission set.
type GroupPermissionPolicy struct {
	Permission entity.Permission
}

func (gp GroupPermissionPolicy) Resolve(a entity.Album, u entity.User) bool {
	for _, group := range u.Groups {
		if entity.HasGroupPermission(a, group.Name, gp.Permission) {
			return true
		}
	}

	return false
}

type AnyGroupPermissionPolicy struct{}

func (ap AnyGroupPermissionPolicy) Resolve(a entity.Album, u entity.User) bool {
	for _, g := range u.Groups {
		if entity.HasGroupPermissions(a, g.Name) {
			return true
		}
	}

	return false
}
