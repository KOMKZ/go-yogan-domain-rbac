package repository

import (
	"context"

	"github.com/KOMKZ/go-yogan-domain-rbac/model"
)

type PermissionRepository interface {
	FindAll(ctx context.Context) ([]model.Permission, error)
	FindByIDs(ctx context.Context, ids []uint) ([]model.Permission, error)
	FindByRoleID(ctx context.Context, roleID uint) ([]model.Permission, error)
	AssignToRole(ctx context.Context, roleID uint, permissionIDs []uint) error
}
