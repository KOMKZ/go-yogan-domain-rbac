package repository

import (
	"context"

	"github.com/KOMKZ/go-yogan-domain-rbac/model"
)

type UserRoleRepository interface {
	FindRolesByUserID(ctx context.Context, userID uint) ([]model.Role, error)
	FindRoleIDsByUserID(ctx context.Context, userID uint) ([]uint, error)
	AssignRolesToUser(ctx context.Context, userID uint, roleIDs []uint) error
	FindPermissionsByUserID(ctx context.Context, userID uint) ([]model.Permission, error)
	FindPermissionCodesByUserID(ctx context.Context, userID uint) ([]string, error)
	HasRole(ctx context.Context, userID uint, roleCode string) (bool, error)
	HasPermission(ctx context.Context, userID uint, permissionCode string) (bool, error)
}
