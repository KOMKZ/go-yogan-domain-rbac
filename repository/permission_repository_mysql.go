package repository

import (
	"context"

	"github.com/KOMKZ/go-yogan-domain-rbac/model"
	"gorm.io/gorm"
)

type PermissionMySQLRepository struct {
	db *gorm.DB
}

func NewPermissionMySQLRepository(db *gorm.DB) *PermissionMySQLRepository {
	return &PermissionMySQLRepository{db: db}
}

func (r *PermissionMySQLRepository) FindAll(ctx context.Context) ([]model.Permission, error) {
	var perms []model.Permission
	err := r.db.WithContext(ctx).Order("group_code, id").Find(&perms).Error
	return perms, err
}

func (r *PermissionMySQLRepository) FindByIDs(ctx context.Context, ids []uint) ([]model.Permission, error) {
	var perms []model.Permission
	if len(ids) == 0 {
		return perms, nil
	}
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&perms).Error
	return perms, err
}

func (r *PermissionMySQLRepository) FindByRoleID(ctx context.Context, roleID uint) ([]model.Permission, error) {
	var perms []model.Permission
	err := r.db.WithContext(ctx).
		Joins("JOIN sys_role_permissions ON sys_role_permissions.permission_id = sys_permissions.id").
		Where("sys_role_permissions.role_id = ?", roleID).
		Find(&perms).Error
	return perms, err
}

func (r *PermissionMySQLRepository) AssignToRole(ctx context.Context, roleID uint, permissionIDs []uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&model.RolePermission{}).Error; err != nil {
			return err
		}
		if len(permissionIDs) == 0 {
			return nil
		}
		records := make([]model.RolePermission, 0, len(permissionIDs))
		for _, pid := range permissionIDs {
			records = append(records, model.RolePermission{RoleID: roleID, PermissionID: pid})
		}
		return tx.Create(&records).Error
	})
}
