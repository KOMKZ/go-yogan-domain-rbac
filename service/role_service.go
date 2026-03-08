package service

import (
	"context"

	domainerrors "github.com/KOMKZ/go-yogan-domain-rbac/errors"
	"github.com/KOMKZ/go-yogan-domain-rbac/model"
	"github.com/KOMKZ/go-yogan-domain-rbac/repository"
	"github.com/KOMKZ/go-yogan-framework/logger"
	"go.uber.org/zap"
)

type RoleService struct {
	roleRepo repository.RoleRepository
	permRepo repository.PermissionRepository
	logger   *logger.CtxZapLogger
}

func NewRoleService(roleRepo repository.RoleRepository, permRepo repository.PermissionRepository, log *logger.CtxZapLogger) *RoleService {
	return &RoleService{roleRepo: roleRepo, permRepo: permRepo, logger: log}
}

type CreateRoleInput struct {
	RoleCode    string
	RoleName    string
	Description string
}

type UpdateRoleInput struct {
	RoleName    *string
	Description *string
	Status      *int
}

type PageResult struct {
	Records []model.Role `json:"records"`
	Total   int64        `json:"total"`
	Size    int          `json:"size"`
	Current int          `json:"current"`
	Pages   int          `json:"pages"`
}

func (s *RoleService) Create(ctx context.Context, input CreateRoleInput) (*model.Role, error) {
	existing, err := s.roleRepo.FindByCode(ctx, input.RoleCode)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, domainerrors.ErrRoleCodeExists
	}

	role := &model.Role{
		RoleCode:    input.RoleCode,
		RoleName:    input.RoleName,
		Description: input.Description,
		Status:      1,
	}

	if err := s.roleRepo.Create(ctx, role); err != nil {
		s.logger.ErrorCtx(ctx, "create role failed", zap.String("code", input.RoleCode), zap.Error(err))
		return nil, err
	}
	s.logger.InfoCtx(ctx, "role created", zap.Uint("role_id", role.ID), zap.String("code", role.RoleCode))
	return role, nil
}

func (s *RoleService) Update(ctx context.Context, id uint, input UpdateRoleInput) (*model.Role, error) {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, domainerrors.ErrRoleNotFound
	}

	if input.RoleName != nil {
		role.RoleName = *input.RoleName
	}
	if input.Description != nil {
		role.Description = *input.Description
	}
	if input.Status != nil {
		role.Status = *input.Status
	}

	if err := s.roleRepo.Update(ctx, role); err != nil {
		return nil, err
	}
	return role, nil
}

func (s *RoleService) Delete(ctx context.Context, id uint) error {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if role == nil {
		return domainerrors.ErrRoleNotFound
	}
	if err := s.roleRepo.Delete(ctx, id); err != nil {
		s.logger.ErrorCtx(ctx, "delete role failed", zap.Uint("role_id", id), zap.Error(err))
		return err
	}
	s.logger.InfoCtx(ctx, "role deleted", zap.Uint("role_id", id))
	return nil
}

func (s *RoleService) GetByID(ctx context.Context, id uint) (*model.Role, error) {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, domainerrors.ErrRoleNotFound
	}
	return role, nil
}

func (s *RoleService) Paginate(ctx context.Context, page, pageSize int, keyword string) (*PageResult, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	roles, total, err := s.roleRepo.Paginate(ctx, page, pageSize, keyword)
	if err != nil {
		return nil, err
	}

	pages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		pages++
	}

	return &PageResult{
		Records: roles,
		Total:   total,
		Size:    pageSize,
		Current: page,
		Pages:   pages,
	}, nil
}

// PermissionGroup represents a group of permissions by GroupCode.
type PermissionGroup struct {
	GroupCode   string             `json:"group_code"`
	GroupName   string             `json:"group_name"`
	Permissions []PermissionItemVO `json:"permissions"`
}

type PermissionItemVO struct {
	ID           uint   `json:"id"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	ResourceCode string `json:"resource_code"`
	Type         string `json:"type"`
	Selected     bool   `json:"selected"`
}

type RolePermissionsVO struct {
	Groups          []PermissionGroup `json:"groups"`
	FlatPermissions []string          `json:"flat_permissions"`
}

var groupNames = map[string]string{
	"SYSTEM": "系统管理",
}

func (s *RoleService) GetAllPermissions(ctx context.Context) ([]PermissionGroup, error) {
	permissions, err := s.permRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return groupPermissions(permissions, nil), nil
}

func (s *RoleService) GetRolePermissions(ctx context.Context, roleID uint) (*RolePermissionsVO, error) {
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, domainerrors.ErrRoleNotFound
	}

	allPerms, err := s.permRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	rolePerms, err := s.permRepo.FindByRoleID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	selectedIDs := make(map[uint]bool, len(rolePerms))
	flatPermissions := make([]string, 0, len(rolePerms))
	for _, p := range rolePerms {
		selectedIDs[p.ID] = true
		flatPermissions = append(flatPermissions, p.PermissionCode)
	}

	return &RolePermissionsVO{
		Groups:          groupPermissions(allPerms, selectedIDs),
		FlatPermissions: flatPermissions,
	}, nil
}

func (s *RoleService) AssignPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return err
	}
	if role == nil {
		return domainerrors.ErrRoleNotFound
	}
	return s.permRepo.AssignToRole(ctx, roleID, permissionIDs)
}

func groupPermissions(permissions []model.Permission, selectedIDs map[uint]bool) []PermissionGroup {
	groupMap := make(map[string]*PermissionGroup)
	groupOrder := make([]string, 0)

	for _, p := range permissions {
		group, exists := groupMap[p.GroupCode]
		if !exists {
			name := groupNames[p.GroupCode]
			if name == "" {
				name = p.GroupCode
			}
			group = &PermissionGroup{
				GroupCode:   p.GroupCode,
				GroupName:   name,
				Permissions: make([]PermissionItemVO, 0),
			}
			groupMap[p.GroupCode] = group
			groupOrder = append(groupOrder, p.GroupCode)
		}

		selected := false
		if selectedIDs != nil {
			selected = selectedIDs[p.ID]
		}

		group.Permissions = append(group.Permissions, PermissionItemVO{
			ID:           p.ID,
			Code:         p.PermissionCode,
			Name:         p.PermissionName,
			ResourceCode: p.ResourceCode,
			Type:         p.PermissionType,
			Selected:     selected,
		})
	}

	result := make([]PermissionGroup, 0, len(groupOrder))
	for _, code := range groupOrder {
		result = append(result, *groupMap[code])
	}
	return result
}
