package repository

import (
	"context"
	"errors"

	"github.com/KOMKZ/go-yogan-domain-rbac/model"
	"github.com/KOMKZ/go-yogan-framework/database"
	"gorm.io/gorm"
)

type RoleMySQLRepository struct {
	base *database.BaseRepository[model.Role]
	db   *gorm.DB
}

func NewRoleMySQLRepository(db *gorm.DB) *RoleMySQLRepository {
	return &RoleMySQLRepository{
		base: database.NewBaseRepository[model.Role](db),
		db:   db,
	}
}

func (r *RoleMySQLRepository) Create(ctx context.Context, role *model.Role) error {
	return r.base.Create(ctx, role)
}

func (r *RoleMySQLRepository) Update(ctx context.Context, role *model.Role) error {
	return r.base.Update(ctx, role)
}

func (r *RoleMySQLRepository) Delete(ctx context.Context, id uint) error {
	return r.base.Delete(ctx, id)
}

func (r *RoleMySQLRepository) FindByID(ctx context.Context, id uint) (*model.Role, error) {
	result, err := r.base.FindByID(ctx, id)
	if errors.Is(err, database.ErrRecordNotFound) {
		return nil, nil
	}
	return result, err
}

func (r *RoleMySQLRepository) FindByCode(ctx context.Context, code string) (*model.Role, error) {
	var role model.Role
	err := r.db.WithContext(ctx).Where("role_code = ?", code).First(&role).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleMySQLRepository) FindByIDs(ctx context.Context, ids []uint) ([]model.Role, error) {
	var roles []model.Role
	if len(ids) == 0 {
		return roles, nil
	}
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&roles).Error
	return roles, err
}

func (r *RoleMySQLRepository) Paginate(ctx context.Context, page, pageSize int, keyword string) ([]model.Role, int64, error) {
	var roles []model.Role
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Role{})
	if keyword != "" {
		query = query.Where("role_name LIKE ? OR role_code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

func (r *RoleMySQLRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Role{}).Count(&count).Error
	return count, err
}
