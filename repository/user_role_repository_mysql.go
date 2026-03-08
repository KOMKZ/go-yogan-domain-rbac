package repository

import (
	"context"

	"github.com/KOMKZ/go-yogan-domain-rbac/model"
	"gorm.io/gorm"
)

type UserRoleMySQLRepository struct {
	db *gorm.DB
}

func NewUserRoleMySQLRepository(db *gorm.DB) *UserRoleMySQLRepository {
	return &UserRoleMySQLRepository{db: db}
}

func (r *UserRoleMySQLRepository) FindRolesByUserID(ctx context.Context, userID uint) ([]model.Role, error) {
	var roles []model.Role
	err := r.db.WithContext(ctx).
		Joins("JOIN sys_user_roles ON sys_user_roles.role_id = sys_roles.id").
		Where("sys_user_roles.user_id = ? AND sys_roles.status = 1", userID).
		Find(&roles).Error
	return roles, err
}

func (r *UserRoleMySQLRepository) FindRoleIDsByUserID(ctx context.Context, userID uint) ([]uint, error) {
	var ids []uint
	err := r.db.WithContext(ctx).Model(&model.UserRole{}).
		Where("user_id = ?", userID).Pluck("role_id", &ids).Error
	return ids, err
}

func (r *UserRoleMySQLRepository) AssignRolesToUser(ctx context.Context, userID uint, roleIDs []uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserRole{}).Error; err != nil {
			return err
		}
		if len(roleIDs) == 0 {
			return nil
		}
		records := make([]model.UserRole, 0, len(roleIDs))
		for _, rid := range roleIDs {
			records = append(records, model.UserRole{UserID: userID, RoleID: rid})
		}
		return tx.Create(&records).Error
	})
}

func (r *UserRoleMySQLRepository) FindPermissionsByUserID(ctx context.Context, userID uint) ([]model.Permission, error) {
	var perms []model.Permission
	err := r.db.WithContext(ctx).
		Joins("JOIN sys_role_permissions ON sys_role_permissions.permission_id = sys_permissions.id").
		Joins("JOIN sys_user_roles ON sys_user_roles.role_id = sys_role_permissions.role_id").
		Joins("JOIN sys_roles ON sys_roles.id = sys_user_roles.role_id").
		Where("sys_user_roles.user_id = ? AND sys_roles.status = 1", userID).
		Distinct().
		Find(&perms).Error
	return perms, err
}

func (r *UserRoleMySQLRepository) FindPermissionCodesByUserID(ctx context.Context, userID uint) ([]string, error) {
	var codes []string
	err := r.db.WithContext(ctx).
		Model(&model.Permission{}).
		Joins("JOIN sys_role_permissions ON sys_role_permissions.permission_id = sys_permissions.id").
		Joins("JOIN sys_user_roles ON sys_user_roles.role_id = sys_role_permissions.role_id").
		Joins("JOIN sys_roles ON sys_roles.id = sys_user_roles.role_id").
		Where("sys_user_roles.user_id = ? AND sys_roles.status = 1", userID).
		Distinct().
		Pluck("sys_permissions.permission_code", &codes).Error
	return codes, err
}

func (r *UserRoleMySQLRepository) HasRole(ctx context.Context, userID uint, roleCode string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.UserRole{}).
		Joins("JOIN sys_roles ON sys_roles.id = sys_user_roles.role_id").
		Where("sys_user_roles.user_id = ? AND sys_roles.role_code = ? AND sys_roles.status = 1", userID, roleCode).
		Count(&count).Error
	return count > 0, err
}

func (r *UserRoleMySQLRepository) HasPermission(ctx context.Context, userID uint, permissionCode string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.RolePermission{}).
		Joins("JOIN sys_user_roles ON sys_user_roles.role_id = sys_role_permissions.role_id").
		Joins("JOIN sys_roles ON sys_roles.id = sys_user_roles.role_id").
		Joins("JOIN sys_permissions ON sys_permissions.id = sys_role_permissions.permission_id").
		Where("sys_user_roles.user_id = ? AND sys_permissions.permission_code = ? AND sys_roles.status = 1", userID, permissionCode).
		Count(&count).Error
	return count > 0, err
}
