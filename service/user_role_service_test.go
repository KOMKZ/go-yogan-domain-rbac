package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	domainerrors "github.com/KOMKZ/go-yogan-domain-rbac/errors"
	"github.com/KOMKZ/go-yogan-domain-rbac/model"
	"github.com/KOMKZ/go-yogan-domain-rbac/repository"
	"github.com/KOMKZ/go-yogan-framework/logger"
)

// --- mock user role repo ---

type mockUserRoleRepo struct {
	userRoles       map[uint][]uint // userID -> roleIDs
	hasRoleResult   map[string]bool // "userID:roleCode" -> bool
	hasPermResult   map[string]bool // "userID:permCode" -> bool
	permCodes       map[uint][]string
	permissions     map[uint][]model.Permission
	roles           map[uint][]model.Role
	assignErr       error
	findRolesErr    error
	findRoleIDsErr  error
	findPermsErr    error
	findPermCodesErr error
	hasRoleErr      error
	hasPermErr      error
}

func newMockUserRoleRepo() *mockUserRoleRepo {
	return &mockUserRoleRepo{
		userRoles:     make(map[uint][]uint),
		hasRoleResult: make(map[string]bool),
		hasPermResult: make(map[string]bool),
		permCodes:     make(map[uint][]string),
		permissions:   make(map[uint][]model.Permission),
		roles:         make(map[uint][]model.Role),
	}
}

func (m *mockUserRoleRepo) FindRolesByUserID(_ context.Context, userID uint) ([]model.Role, error) {
	if m.findRolesErr != nil {
		return nil, m.findRolesErr
	}
	return m.roles[userID], nil
}

func (m *mockUserRoleRepo) FindRoleIDsByUserID(_ context.Context, userID uint) ([]uint, error) {
	if m.findRoleIDsErr != nil {
		return nil, m.findRoleIDsErr
	}
	return m.userRoles[userID], nil
}

func (m *mockUserRoleRepo) AssignRolesToUser(_ context.Context, userID uint, roleIDs []uint) error {
	if m.assignErr != nil {
		return m.assignErr
	}
	m.userRoles[userID] = roleIDs
	return nil
}

func (m *mockUserRoleRepo) FindPermissionsByUserID(_ context.Context, userID uint) ([]model.Permission, error) {
	if m.findPermsErr != nil {
		return nil, m.findPermsErr
	}
	return m.permissions[userID], nil
}

func (m *mockUserRoleRepo) FindPermissionCodesByUserID(_ context.Context, userID uint) ([]string, error) {
	if m.findPermCodesErr != nil {
		return nil, m.findPermCodesErr
	}
	return m.permCodes[userID], nil
}

func (m *mockUserRoleRepo) HasRole(_ context.Context, userID uint, roleCode string) (bool, error) {
	if m.hasRoleErr != nil {
		return false, m.hasRoleErr
	}
	key := fmt.Sprintf("%d:%s", userID, roleCode)
	return m.hasRoleResult[key], nil
}

func (m *mockUserRoleRepo) HasPermission(_ context.Context, userID uint, permissionCode string) (bool, error) {
	if m.hasPermErr != nil {
		return false, m.hasPermErr
	}
	key := fmt.Sprintf("%d:%s", userID, permissionCode)
	return m.hasPermResult[key], nil
}

var _ repository.UserRoleRepository = (*mockUserRoleRepo)(nil)

func fmt_key(userID uint, code string) string {
	return fmt.Sprintf("%d:%s", userID, code)
}

func TestUserRoleService_GetUserRoles_Success(t *testing.T) {
	urRepo := newMockUserRoleRepo()
	urRepo.roles[1] = []model.Role{{ID: 10, RoleCode: "admin"}}
	svc := NewUserRoleService(newMockRoleRepo(), urRepo, logger.GetLogger("rbac_test"))
	roles, err := svc.GetUserRoles(context.Background(), 1)
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if len(roles) != 1 || roles[0].ID != 10 {
		t.Errorf("roles = %v", roles)
	}
}

func TestUserRoleService_GetUserRoles_RepoErr(t *testing.T) {
	urRepo := newMockUserRoleRepo()
	urRepo.findRolesErr = errors.New("find roles failed")
	svc := NewUserRoleService(newMockRoleRepo(), urRepo, logger.GetLogger("rbac_test"))
	_, err := svc.GetUserRoles(context.Background(), 1)
	if err == nil || err.Error() != "find roles failed" {
		t.Errorf("err = %v", err)
	}
}

func TestUserRoleService_AssignUserRoles_Success(t *testing.T) {
	roleRepo := newMockRoleRepo()
	ctx := context.Background()
	_ = roleRepo.Create(ctx, &model.Role{RoleCode: "admin", RoleName: "Admin", Status: 1})
	_ = roleRepo.Create(ctx, &model.Role{RoleCode: "editor", RoleName: "Editor", Status: 1})

	urRepo := newMockUserRoleRepo()
	svc := NewUserRoleService(roleRepo, urRepo, logger.GetLogger("rbac_test"))
	err := svc.AssignUserRoles(ctx, 1, []uint{1, 2})
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if len(urRepo.userRoles[1]) != 2 {
		t.Errorf("assigned = %v", urRepo.userRoles[1])
	}
}

func TestUserRoleService_AssignUserRoles_EmptyRoles(t *testing.T) {
	svc := NewUserRoleService(newMockRoleRepo(), newMockUserRoleRepo(), logger.GetLogger("rbac_test"))
	err := svc.AssignUserRoles(context.Background(), 1, []uint{})
	if err != nil {
		t.Fatalf("err = %v", err)
	}
}

func TestUserRoleService_AssignUserRoles_RoleNotFound(t *testing.T) {
	roleRepo := newMockRoleRepo()
	svc := NewUserRoleService(roleRepo, newMockUserRoleRepo(), logger.GetLogger("rbac_test"))
	err := svc.AssignUserRoles(context.Background(), 1, []uint{999})
	if err != domainerrors.ErrRoleNotFound {
		t.Errorf("err = %v, want ErrRoleNotFound", err)
	}
}

func TestUserRoleService_AssignUserRoles_RoleDisabled(t *testing.T) {
	roleRepo := newMockRoleRepo()
	ctx := context.Background()
	_ = roleRepo.Create(ctx, &model.Role{RoleCode: "disabled", RoleName: "D", Status: 0})

	svc := NewUserRoleService(roleRepo, newMockUserRoleRepo(), logger.GetLogger("rbac_test"))
	err := svc.AssignUserRoles(ctx, 1, []uint{1})
	if err != domainerrors.ErrRoleDisabled {
		t.Errorf("err = %v, want ErrRoleDisabled", err)
	}
}

func TestUserRoleService_AssignUserRoles_FindByIDsErr(t *testing.T) {
	roleRepo := newMockRoleRepo()
	roleRepo.findByIDsErr = errors.New("find ids failed")
	svc := NewUserRoleService(roleRepo, newMockUserRoleRepo(), logger.GetLogger("rbac_test"))
	err := svc.AssignUserRoles(context.Background(), 1, []uint{1})
	if err == nil || err.Error() != "find ids failed" {
		t.Errorf("err = %v", err)
	}
}

func TestUserRoleService_AssignUserRoles_AssignErr(t *testing.T) {
	roleRepo := newMockRoleRepo()
	ctx := context.Background()
	_ = roleRepo.Create(ctx, &model.Role{RoleCode: "a", RoleName: "A", Status: 1})

	urRepo := newMockUserRoleRepo()
	urRepo.assignErr = errors.New("assign failed")
	svc := NewUserRoleService(roleRepo, urRepo, logger.GetLogger("rbac_test"))
	err := svc.AssignUserRoles(ctx, 1, []uint{1})
	if err == nil || err.Error() != "assign failed" {
		t.Errorf("err = %v", err)
	}
}

func TestUserRoleService_GetUserPermissions(t *testing.T) {
	urRepo := newMockUserRoleRepo()
	urRepo.permissions[1] = []model.Permission{{ID: 1, PermissionCode: "user:read"}}
	svc := NewUserRoleService(newMockRoleRepo(), urRepo, logger.GetLogger("rbac_test"))
	perms, err := svc.GetUserPermissions(context.Background(), 1)
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if len(perms) != 1 {
		t.Errorf("len = %d", len(perms))
	}
}

func TestUserRoleService_GetUserPermissionCodes(t *testing.T) {
	urRepo := newMockUserRoleRepo()
	urRepo.permCodes[1] = []string{"user:read", "user:write"}
	svc := NewUserRoleService(newMockRoleRepo(), urRepo, logger.GetLogger("rbac_test"))
	codes, err := svc.GetUserPermissionCodes(context.Background(), 1)
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if len(codes) != 2 {
		t.Errorf("len = %d", len(codes))
	}
}

func TestUserRoleService_HasRole(t *testing.T) {
	urRepo := newMockUserRoleRepo()
	urRepo.hasRoleResult[fmt_key(1, "admin")] = true
	svc := NewUserRoleService(newMockRoleRepo(), urRepo, logger.GetLogger("rbac_test"))

	has, err := svc.HasRole(context.Background(), 1, "admin")
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if !has {
		t.Error("HasRole should be true")
	}

	has, err = svc.HasRole(context.Background(), 1, "editor")
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if has {
		t.Error("HasRole should be false")
	}
}

func TestUserRoleService_HasPermission(t *testing.T) {
	urRepo := newMockUserRoleRepo()
	urRepo.hasPermResult[fmt_key(1, "user:read")] = true
	svc := NewUserRoleService(newMockRoleRepo(), urRepo, logger.GetLogger("rbac_test"))

	has, err := svc.HasPermission(context.Background(), 1, "user:read")
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if !has {
		t.Error("HasPermission should be true")
	}
}

func TestUserRoleService_Can(t *testing.T) {
	urRepo := newMockUserRoleRepo()
	urRepo.hasPermResult[fmt_key(1, "user:read")] = true
	svc := NewUserRoleService(newMockRoleRepo(), urRepo, logger.GetLogger("rbac_test"))

	has, err := svc.Can(context.Background(), 1, "user:read")
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if !has {
		t.Error("Can should be true")
	}
}

func TestUserRoleService_HasAnyRole(t *testing.T) {
	urRepo := newMockUserRoleRepo()
	urRepo.hasRoleResult[fmt_key(1, "admin")] = true
	svc := NewUserRoleService(newMockRoleRepo(), urRepo, logger.GetLogger("rbac_test"))

	has, err := svc.HasAnyRole(context.Background(), 1, "editor", "admin")
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if !has {
		t.Error("HasAnyRole should be true")
	}

	has, err = svc.HasAnyRole(context.Background(), 1, "editor", "viewer")
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if has {
		t.Error("HasAnyRole should be false")
	}
}

func TestUserRoleService_HasAllRoles(t *testing.T) {
	urRepo := newMockUserRoleRepo()
	urRepo.hasRoleResult[fmt_key(1, "admin")] = true
	urRepo.hasRoleResult[fmt_key(1, "editor")] = true
	svc := NewUserRoleService(newMockRoleRepo(), urRepo, logger.GetLogger("rbac_test"))

	has, err := svc.HasAllRoles(context.Background(), 1, "admin", "editor")
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if !has {
		t.Error("HasAllRoles should be true")
	}

	has, err = svc.HasAllRoles(context.Background(), 1, "admin", "viewer")
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if has {
		t.Error("HasAllRoles should be false")
	}
}

func TestUserRoleService_HasAnyPermission(t *testing.T) {
	urRepo := newMockUserRoleRepo()
	urRepo.hasPermResult[fmt_key(1, "user:read")] = true
	svc := NewUserRoleService(newMockRoleRepo(), urRepo, logger.GetLogger("rbac_test"))

	has, err := svc.HasAnyPermission(context.Background(), 1, "user:write", "user:read")
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if !has {
		t.Error("HasAnyPermission should be true")
	}
}

func TestUserRoleService_HasAllPermissions(t *testing.T) {
	urRepo := newMockUserRoleRepo()
	urRepo.hasPermResult[fmt_key(1, "user:read")] = true
	urRepo.hasPermResult[fmt_key(1, "user:write")] = true
	svc := NewUserRoleService(newMockRoleRepo(), urRepo, logger.GetLogger("rbac_test"))

	has, err := svc.HasAllPermissions(context.Background(), 1, "user:read", "user:write")
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if !has {
		t.Error("HasAllPermissions should be true")
	}
}

func TestUserRoleService_HasRole_RepoErr(t *testing.T) {
	urRepo := newMockUserRoleRepo()
	urRepo.hasRoleErr = errors.New("has role failed")
	svc := NewUserRoleService(newMockRoleRepo(), urRepo, logger.GetLogger("rbac_test"))
	_, err := svc.HasRole(context.Background(), 1, "admin")
	if err == nil || err.Error() != "has role failed" {
		t.Errorf("err = %v", err)
	}
}

func TestUserRoleService_HasAnyRole_RepoErr(t *testing.T) {
	urRepo := newMockUserRoleRepo()
	urRepo.hasRoleErr = errors.New("error")
	svc := NewUserRoleService(newMockRoleRepo(), urRepo, logger.GetLogger("rbac_test"))
	_, err := svc.HasAnyRole(context.Background(), 1, "admin")
	if err == nil {
		t.Error("expected error")
	}
}

func TestUserRoleService_HasAllRoles_RepoErr(t *testing.T) {
	urRepo := newMockUserRoleRepo()
	urRepo.hasRoleErr = errors.New("error")
	svc := NewUserRoleService(newMockRoleRepo(), urRepo, logger.GetLogger("rbac_test"))
	_, err := svc.HasAllRoles(context.Background(), 1, "admin")
	if err == nil {
		t.Error("expected error")
	}
}

func TestUserRoleService_HasPermission_RepoErr(t *testing.T) {
	urRepo := newMockUserRoleRepo()
	urRepo.hasPermErr = errors.New("error")
	svc := NewUserRoleService(newMockRoleRepo(), urRepo, logger.GetLogger("rbac_test"))
	_, err := svc.HasPermission(context.Background(), 1, "user:read")
	if err == nil {
		t.Error("expected error")
	}
}

func TestUserRoleService_HasAnyPermission_RepoErr(t *testing.T) {
	urRepo := newMockUserRoleRepo()
	urRepo.hasPermErr = errors.New("error")
	svc := NewUserRoleService(newMockRoleRepo(), urRepo, logger.GetLogger("rbac_test"))
	_, err := svc.HasAnyPermission(context.Background(), 1, "user:read")
	if err == nil {
		t.Error("expected error")
	}
}

func TestUserRoleService_HasAllPermissions_RepoErr(t *testing.T) {
	urRepo := newMockUserRoleRepo()
	urRepo.hasPermErr = errors.New("error")
	svc := NewUserRoleService(newMockRoleRepo(), urRepo, logger.GetLogger("rbac_test"))
	_, err := svc.HasAllPermissions(context.Background(), 1, "user:read")
	if err == nil {
		t.Error("expected error")
	}
}
