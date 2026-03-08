package service

import (
	"context"

	domainerrors "github.com/KOMKZ/go-yogan-domain-rbac/errors"
	"github.com/KOMKZ/go-yogan-domain-rbac/model"
	"github.com/KOMKZ/go-yogan-domain-rbac/repository"
	"github.com/KOMKZ/go-yogan-framework/logger"
	"go.uber.org/zap"
)

type UserRoleService struct {
	roleRepo     repository.RoleRepository
	userRoleRepo repository.UserRoleRepository
	logger       *logger.CtxZapLogger
}

func NewUserRoleService(roleRepo repository.RoleRepository, userRoleRepo repository.UserRoleRepository, log *logger.CtxZapLogger) *UserRoleService {
	return &UserRoleService{roleRepo: roleRepo, userRoleRepo: userRoleRepo, logger: log}
}

func (s *UserRoleService) GetUserRoles(ctx context.Context, userID uint) ([]model.Role, error) {
	return s.userRoleRepo.FindRolesByUserID(ctx, userID)
}

func (s *UserRoleService) AssignUserRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	if len(roleIDs) == 0 {
		return s.userRoleRepo.AssignRolesToUser(ctx, userID, roleIDs)
	}

	roles, err := s.roleRepo.FindByIDs(ctx, roleIDs)
	if err != nil {
		return err
	}

	found := make(map[uint]bool, len(roles))
	for _, r := range roles {
		found[r.ID] = true
		if !r.IsEnabled() {
			return domainerrors.ErrRoleDisabled
		}
	}
	for _, id := range roleIDs {
		if !found[id] {
			return domainerrors.ErrRoleNotFound
		}
	}

	if err := s.userRoleRepo.AssignRolesToUser(ctx, userID, roleIDs); err != nil {
		s.logger.ErrorCtx(ctx, "assign user roles failed", zap.Uint("user_id", userID), zap.Error(err))
		return err
	}
	s.logger.InfoCtx(ctx, "user roles assigned", zap.Uint("user_id", userID), zap.Any("role_ids", roleIDs))
	return nil
}

func (s *UserRoleService) GetUserPermissions(ctx context.Context, userID uint) ([]model.Permission, error) {
	return s.userRoleRepo.FindPermissionsByUserID(ctx, userID)
}

func (s *UserRoleService) GetUserPermissionCodes(ctx context.Context, userID uint) ([]string, error) {
	return s.userRoleRepo.FindPermissionCodesByUserID(ctx, userID)
}

func (s *UserRoleService) HasRole(ctx context.Context, userID uint, roleCode string) (bool, error) {
	return s.userRoleRepo.HasRole(ctx, userID, roleCode)
}

func (s *UserRoleService) HasPermission(ctx context.Context, userID uint, permissionCode string) (bool, error) {
	return s.userRoleRepo.HasPermission(ctx, userID, permissionCode)
}

func (s *UserRoleService) Can(ctx context.Context, userID uint, permissionCode string) (bool, error) {
	return s.HasPermission(ctx, userID, permissionCode)
}

func (s *UserRoleService) HasAnyRole(ctx context.Context, userID uint, roleCodes ...string) (bool, error) {
	for _, code := range roleCodes {
		has, err := s.userRoleRepo.HasRole(ctx, userID, code)
		if err != nil {
			return false, err
		}
		if has {
			return true, nil
		}
	}
	return false, nil
}

func (s *UserRoleService) HasAllRoles(ctx context.Context, userID uint, roleCodes ...string) (bool, error) {
	for _, code := range roleCodes {
		has, err := s.userRoleRepo.HasRole(ctx, userID, code)
		if err != nil {
			return false, err
		}
		if !has {
			return false, nil
		}
	}
	return true, nil
}

func (s *UserRoleService) HasAnyPermission(ctx context.Context, userID uint, permissionCodes ...string) (bool, error) {
	for _, code := range permissionCodes {
		has, err := s.userRoleRepo.HasPermission(ctx, userID, code)
		if err != nil {
			return false, err
		}
		if has {
			return true, nil
		}
	}
	return false, nil
}

func (s *UserRoleService) HasAllPermissions(ctx context.Context, userID uint, permissionCodes ...string) (bool, error) {
	for _, code := range permissionCodes {
		has, err := s.userRoleRepo.HasPermission(ctx, userID, code)
		if err != nil {
			return false, err
		}
		if !has {
			return false, nil
		}
	}
	return true, nil
}
