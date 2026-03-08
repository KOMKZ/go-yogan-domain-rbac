package repository

import (
	"context"

	"github.com/KOMKZ/go-yogan-domain-rbac/model"
)

type RoleRepository interface {
	Create(ctx context.Context, role *model.Role) error
	Update(ctx context.Context, role *model.Role) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*model.Role, error)
	FindByCode(ctx context.Context, code string) (*model.Role, error)
	FindByIDs(ctx context.Context, ids []uint) ([]model.Role, error)
	Paginate(ctx context.Context, page, pageSize int, keyword string) ([]model.Role, int64, error)
	Count(ctx context.Context) (int64, error)
}
