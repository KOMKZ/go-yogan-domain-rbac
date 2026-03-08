package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	domainerrors "github.com/KOMKZ/go-yogan-domain-rbac/errors"
	"github.com/KOMKZ/go-yogan-domain-rbac/model"
	"github.com/KOMKZ/go-yogan-domain-rbac/repository"
	"github.com/KOMKZ/go-yogan-framework/logger"
)

// --- mock role repo ---

type mockRoleRepo struct {
	mu          sync.RWMutex
	roles       map[uint]*model.Role
	byCode      map[string]*model.Role
	nextID      uint
	createErr   error
	updateErr   error
	deleteErr   error
	findByIDErr error
	findCodeErr error
	paginateErr error
	countErr    error
	findByIDsErr error
}

func newMockRoleRepo() *mockRoleRepo {
	return &mockRoleRepo{
		roles:  make(map[uint]*model.Role),
		byCode: make(map[string]*model.Role),
		nextID: 1,
	}
}

func (m *mockRoleRepo) Create(_ context.Context, role *model.Role) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if role.ID == 0 {
		role.ID = m.nextID
		m.nextID++
	}
	cp := copyRole(role)
	m.roles[role.ID] = cp
	m.byCode[role.RoleCode] = cp
	return nil
}

func (m *mockRoleRepo) Update(_ context.Context, role *model.Role) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	old := m.roles[role.ID]
	if old != nil {
		delete(m.byCode, old.RoleCode)
	}
	cp := copyRole(role)
	m.roles[role.ID] = cp
	m.byCode[role.RoleCode] = cp
	return nil
}

func (m *mockRoleRepo) Delete(_ context.Context, id uint) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	r := m.roles[id]
	if r != nil {
		delete(m.byCode, r.RoleCode)
		delete(m.roles, id)
	}
	return nil
}

func (m *mockRoleRepo) FindByID(_ context.Context, id uint) (*model.Role, error) {
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	r, ok := m.roles[id]
	if !ok {
		return nil, nil
	}
	return copyRole(r), nil
}

func (m *mockRoleRepo) FindByCode(_ context.Context, code string) (*model.Role, error) {
	if m.findCodeErr != nil {
		return nil, m.findCodeErr
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	r, ok := m.byCode[code]
	if !ok {
		return nil, nil
	}
	return copyRole(r), nil
}

func (m *mockRoleRepo) FindByIDs(_ context.Context, ids []uint) ([]model.Role, error) {
	if m.findByIDsErr != nil {
		return nil, m.findByIDsErr
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []model.Role
	for _, id := range ids {
		if r, ok := m.roles[id]; ok {
			result = append(result, *copyRole(r))
		}
	}
	return result, nil
}

func (m *mockRoleRepo) Paginate(_ context.Context, page, pageSize int, _ string) ([]model.Role, int64, error) {
	if m.paginateErr != nil {
		return nil, 0, m.paginateErr
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	var list []model.Role
	for _, r := range m.roles {
		list = append(list, *copyRole(r))
	}
	total := int64(len(list))
	offset := (page - 1) * pageSize
	if offset >= len(list) {
		return []model.Role{}, total, nil
	}
	end := offset + pageSize
	if end > len(list) {
		end = len(list)
	}
	return list[offset:end], total, nil
}

func (m *mockRoleRepo) Count(_ context.Context) (int64, error) {
	if m.countErr != nil {
		return 0, m.countErr
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	return int64(len(m.roles)), nil
}

func copyRole(r *model.Role) *model.Role {
	if r == nil {
		return nil
	}
	cp := *r
	return &cp
}

var _ repository.RoleRepository = (*mockRoleRepo)(nil)

// --- mock permission repo ---

type mockPermRepo struct {
	mu          sync.RWMutex
	permissions []model.Permission
	rolePerms   map[uint][]uint // roleID -> permissionIDs
	findAllErr  error
	findByIDsErr error
	findByRoleErr error
	assignErr   error
}

func newMockPermRepo() *mockPermRepo {
	return &mockPermRepo{
		rolePerms: make(map[uint][]uint),
	}
}

func (m *mockPermRepo) FindAll(_ context.Context) ([]model.Permission, error) {
	if m.findAllErr != nil {
		return nil, m.findAllErr
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]model.Permission, len(m.permissions))
	copy(result, m.permissions)
	return result, nil
}

func (m *mockPermRepo) FindByIDs(_ context.Context, ids []uint) ([]model.Permission, error) {
	if m.findByIDsErr != nil {
		return nil, m.findByIDsErr
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	idSet := make(map[uint]bool, len(ids))
	for _, id := range ids {
		idSet[id] = true
	}
	var result []model.Permission
	for _, p := range m.permissions {
		if idSet[p.ID] {
			result = append(result, p)
		}
	}
	return result, nil
}

func (m *mockPermRepo) FindByRoleID(_ context.Context, roleID uint) ([]model.Permission, error) {
	if m.findByRoleErr != nil {
		return nil, m.findByRoleErr
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	permIDs := m.rolePerms[roleID]
	idSet := make(map[uint]bool, len(permIDs))
	for _, id := range permIDs {
		idSet[id] = true
	}
	var result []model.Permission
	for _, p := range m.permissions {
		if idSet[p.ID] {
			result = append(result, p)
		}
	}
	return result, nil
}

func (m *mockPermRepo) AssignToRole(_ context.Context, roleID uint, permissionIDs []uint) error {
	if m.assignErr != nil {
		return m.assignErr
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.rolePerms[roleID] = permissionIDs
	return nil
}

func (m *mockPermRepo) seedPermissions(perms []model.Permission) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.permissions = perms
}

var _ repository.PermissionRepository = (*mockPermRepo)(nil)

// --- tests ---

func intPtr(v int) *int { return &v }

func TestRoleService_Create_Success(t *testing.T) {
	roleRepo := newMockRoleRepo()
	permRepo := newMockPermRepo()
	svc := NewRoleService(roleRepo, permRepo, logger.GetLogger("rbac_test"))
	ctx := context.Background()

	role, err := svc.Create(ctx, CreateRoleInput{
		RoleCode: "admin", RoleName: "管理员", Description: "desc",
	})
	if err != nil {
		t.Fatalf("Create() err = %v", err)
	}
	if role.RoleCode != "admin" {
		t.Errorf("RoleCode = %q, want admin", role.RoleCode)
	}
	if role.Status != 1 {
		t.Errorf("Status = %d, want 1", role.Status)
	}
}

func TestRoleService_Create_CodeExists(t *testing.T) {
	roleRepo := newMockRoleRepo()
	permRepo := newMockPermRepo()
	svc := NewRoleService(roleRepo, permRepo, logger.GetLogger("rbac_test"))
	ctx := context.Background()

	_, _ = svc.Create(ctx, CreateRoleInput{RoleCode: "admin", RoleName: "A"})
	_, err := svc.Create(ctx, CreateRoleInput{RoleCode: "admin", RoleName: "B"})
	if err != domainerrors.ErrRoleCodeExists {
		t.Errorf("err = %v, want ErrRoleCodeExists", err)
	}
}

func TestRoleService_Create_FindByCodeErr(t *testing.T) {
	roleRepo := newMockRoleRepo()
	roleRepo.findCodeErr = errors.New("db error")
	svc := NewRoleService(roleRepo, newMockPermRepo(), logger.GetLogger("rbac_test"))

	_, err := svc.Create(context.Background(), CreateRoleInput{RoleCode: "x", RoleName: "X"})
	if err == nil || err.Error() != "db error" {
		t.Errorf("err = %v, want db error", err)
	}
}

func TestRoleService_Create_RepoCreateErr(t *testing.T) {
	roleRepo := newMockRoleRepo()
	roleRepo.createErr = errors.New("create failed")
	svc := NewRoleService(roleRepo, newMockPermRepo(), logger.GetLogger("rbac_test"))

	_, err := svc.Create(context.Background(), CreateRoleInput{RoleCode: "x", RoleName: "X"})
	if err == nil || err.Error() != "create failed" {
		t.Errorf("err = %v, want create failed", err)
	}
}

func TestRoleService_Update_Success(t *testing.T) {
	roleRepo := newMockRoleRepo()
	svc := NewRoleService(roleRepo, newMockPermRepo(), logger.GetLogger("rbac_test"))
	ctx := context.Background()

	created, _ := svc.Create(ctx, CreateRoleInput{RoleCode: "ed", RoleName: "Editor"})
	newName := "编辑"
	updated, err := svc.Update(ctx, created.ID, UpdateRoleInput{RoleName: &newName})
	if err != nil {
		t.Fatalf("Update() err = %v", err)
	}
	if updated.RoleName != "编辑" {
		t.Errorf("RoleName = %q, want 编辑", updated.RoleName)
	}
}

func TestRoleService_Update_NotFound(t *testing.T) {
	svc := NewRoleService(newMockRoleRepo(), newMockPermRepo(), logger.GetLogger("rbac_test"))
	n := "x"
	_, err := svc.Update(context.Background(), 999, UpdateRoleInput{RoleName: &n})
	if err != domainerrors.ErrRoleNotFound {
		t.Errorf("err = %v, want ErrRoleNotFound", err)
	}
}

func TestRoleService_Update_FindByIDErr(t *testing.T) {
	roleRepo := newMockRoleRepo()
	roleRepo.findByIDErr = errors.New("find failed")
	svc := NewRoleService(roleRepo, newMockPermRepo(), logger.GetLogger("rbac_test"))
	n := "x"
	_, err := svc.Update(context.Background(), 1, UpdateRoleInput{RoleName: &n})
	if err == nil || err.Error() != "find failed" {
		t.Errorf("err = %v, want find failed", err)
	}
}

func TestRoleService_Update_RepoUpdateErr(t *testing.T) {
	roleRepo := newMockRoleRepo()
	svc := NewRoleService(roleRepo, newMockPermRepo(), logger.GetLogger("rbac_test"))
	ctx := context.Background()
	created, _ := svc.Create(ctx, CreateRoleInput{RoleCode: "a", RoleName: "A"})
	roleRepo.updateErr = errors.New("update failed")
	n := "B"
	_, err := svc.Update(ctx, created.ID, UpdateRoleInput{RoleName: &n})
	if err == nil || err.Error() != "update failed" {
		t.Errorf("err = %v, want update failed", err)
	}
}

func TestRoleService_Update_StatusChange(t *testing.T) {
	roleRepo := newMockRoleRepo()
	svc := NewRoleService(roleRepo, newMockPermRepo(), logger.GetLogger("rbac_test"))
	ctx := context.Background()
	created, _ := svc.Create(ctx, CreateRoleInput{RoleCode: "a", RoleName: "A"})
	updated, err := svc.Update(ctx, created.ID, UpdateRoleInput{Status: intPtr(0)})
	if err != nil {
		t.Fatalf("Update() err = %v", err)
	}
	if updated.Status != 0 {
		t.Errorf("Status = %d, want 0", updated.Status)
	}
}

func TestRoleService_Delete_Success(t *testing.T) {
	roleRepo := newMockRoleRepo()
	svc := NewRoleService(roleRepo, newMockPermRepo(), logger.GetLogger("rbac_test"))
	ctx := context.Background()
	created, _ := svc.Create(ctx, CreateRoleInput{RoleCode: "a", RoleName: "A"})
	if err := svc.Delete(ctx, created.ID); err != nil {
		t.Fatalf("Delete() err = %v", err)
	}
	_, err := svc.GetByID(ctx, created.ID)
	if err != domainerrors.ErrRoleNotFound {
		t.Errorf("GetByID after delete: err = %v, want ErrRoleNotFound", err)
	}
}

func TestRoleService_Delete_NotFound(t *testing.T) {
	svc := NewRoleService(newMockRoleRepo(), newMockPermRepo(), logger.GetLogger("rbac_test"))
	err := svc.Delete(context.Background(), 999)
	if err != domainerrors.ErrRoleNotFound {
		t.Errorf("err = %v, want ErrRoleNotFound", err)
	}
}

func TestRoleService_Delete_FindByIDErr(t *testing.T) {
	roleRepo := newMockRoleRepo()
	roleRepo.findByIDErr = errors.New("find failed")
	svc := NewRoleService(roleRepo, newMockPermRepo(), logger.GetLogger("rbac_test"))
	err := svc.Delete(context.Background(), 1)
	if err == nil || err.Error() != "find failed" {
		t.Errorf("err = %v, want find failed", err)
	}
}

func TestRoleService_Delete_RepoDeleteErr(t *testing.T) {
	roleRepo := newMockRoleRepo()
	svc := NewRoleService(roleRepo, newMockPermRepo(), logger.GetLogger("rbac_test"))
	ctx := context.Background()
	created, _ := svc.Create(ctx, CreateRoleInput{RoleCode: "a", RoleName: "A"})
	roleRepo.deleteErr = errors.New("delete failed")
	err := svc.Delete(ctx, created.ID)
	if err == nil || err.Error() != "delete failed" {
		t.Errorf("err = %v, want delete failed", err)
	}
}

func TestRoleService_GetByID_Success(t *testing.T) {
	roleRepo := newMockRoleRepo()
	svc := NewRoleService(roleRepo, newMockPermRepo(), logger.GetLogger("rbac_test"))
	ctx := context.Background()
	created, _ := svc.Create(ctx, CreateRoleInput{RoleCode: "a", RoleName: "A"})
	got, err := svc.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetByID() err = %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("ID = %d, want %d", got.ID, created.ID)
	}
}

func TestRoleService_GetByID_NotFound(t *testing.T) {
	svc := NewRoleService(newMockRoleRepo(), newMockPermRepo(), logger.GetLogger("rbac_test"))
	_, err := svc.GetByID(context.Background(), 999)
	if err != domainerrors.ErrRoleNotFound {
		t.Errorf("err = %v, want ErrRoleNotFound", err)
	}
}

func TestRoleService_Paginate_Success(t *testing.T) {
	roleRepo := newMockRoleRepo()
	svc := NewRoleService(roleRepo, newMockPermRepo(), logger.GetLogger("rbac_test"))
	ctx := context.Background()
	_, _ = svc.Create(ctx, CreateRoleInput{RoleCode: "a", RoleName: "A"})
	_, _ = svc.Create(ctx, CreateRoleInput{RoleCode: "b", RoleName: "B"})
	res, err := svc.Paginate(ctx, 1, 10, "")
	if err != nil {
		t.Fatalf("Paginate() err = %v", err)
	}
	if res.Total != 2 {
		t.Errorf("Total = %d, want 2", res.Total)
	}
}

func TestRoleService_Paginate_NormalizePage(t *testing.T) {
	svc := NewRoleService(newMockRoleRepo(), newMockPermRepo(), logger.GetLogger("rbac_test"))
	res, err := svc.Paginate(context.Background(), 0, 10, "")
	if err != nil {
		t.Fatalf("Paginate() err = %v", err)
	}
	if res.Current != 1 {
		t.Errorf("Current = %d, want 1", res.Current)
	}
}

func TestRoleService_Paginate_NormalizePageSize(t *testing.T) {
	svc := NewRoleService(newMockRoleRepo(), newMockPermRepo(), logger.GetLogger("rbac_test"))
	res, err := svc.Paginate(context.Background(), 1, 0, "")
	if err != nil {
		t.Fatalf("Paginate() err = %v", err)
	}
	if res.Size != 10 {
		t.Errorf("Size = %d, want 10", res.Size)
	}
}

func TestRoleService_Paginate_PageSizeOver100(t *testing.T) {
	svc := NewRoleService(newMockRoleRepo(), newMockPermRepo(), logger.GetLogger("rbac_test"))
	res, err := svc.Paginate(context.Background(), 1, 200, "")
	if err != nil {
		t.Fatalf("Paginate() err = %v", err)
	}
	if res.Size != 10 {
		t.Errorf("Size = %d, want 10", res.Size)
	}
}

func TestRoleService_Paginate_PagesCalc(t *testing.T) {
	roleRepo := newMockRoleRepo()
	ctx := context.Background()
	for i := 0; i < 25; i++ {
		_ = roleRepo.Create(ctx, &model.Role{RoleCode: fmt.Sprintf("r%d", i), RoleName: fmt.Sprintf("R%d", i), Status: 1})
	}
	svc := NewRoleService(roleRepo, newMockPermRepo(), logger.GetLogger("rbac_test"))
	res, err := svc.Paginate(ctx, 1, 10, "")
	if err != nil {
		t.Fatalf("Paginate() err = %v", err)
	}
	if res.Pages != 3 {
		t.Errorf("Pages = %d, want 3", res.Pages)
	}
}

func TestRoleService_Paginate_RepoErr(t *testing.T) {
	roleRepo := newMockRoleRepo()
	roleRepo.paginateErr = errors.New("paginate failed")
	svc := NewRoleService(roleRepo, newMockPermRepo(), logger.GetLogger("rbac_test"))
	_, err := svc.Paginate(context.Background(), 1, 10, "")
	if err == nil || err.Error() != "paginate failed" {
		t.Errorf("err = %v, want paginate failed", err)
	}
}

func TestRoleService_GetAllPermissions_Success(t *testing.T) {
	permRepo := newMockPermRepo()
	permRepo.seedPermissions([]model.Permission{
		{ID: 1, PermissionCode: "user:read", PermissionName: "查看用户", PermissionType: "READ", ResourceCode: "user", GroupCode: "SYSTEM"},
		{ID: 2, PermissionCode: "user:write", PermissionName: "编辑用户", PermissionType: "WRITE", ResourceCode: "user", GroupCode: "SYSTEM"},
	})
	svc := NewRoleService(newMockRoleRepo(), permRepo, logger.GetLogger("rbac_test"))
	groups, err := svc.GetAllPermissions(context.Background())
	if err != nil {
		t.Fatalf("GetAllPermissions() err = %v", err)
	}
	if len(groups) != 1 {
		t.Fatalf("len(groups) = %d, want 1", len(groups))
	}
	if groups[0].GroupCode != "SYSTEM" {
		t.Errorf("GroupCode = %q, want SYSTEM", groups[0].GroupCode)
	}
	if groups[0].GroupName != "系统管理" {
		t.Errorf("GroupName = %q, want 系统管理", groups[0].GroupName)
	}
	if len(groups[0].Permissions) != 2 {
		t.Errorf("len(Permissions) = %d, want 2", len(groups[0].Permissions))
	}
}

func TestRoleService_GetAllPermissions_FindAllErr(t *testing.T) {
	permRepo := newMockPermRepo()
	permRepo.findAllErr = errors.New("find all failed")
	svc := NewRoleService(newMockRoleRepo(), permRepo, logger.GetLogger("rbac_test"))
	_, err := svc.GetAllPermissions(context.Background())
	if err == nil || err.Error() != "find all failed" {
		t.Errorf("err = %v, want find all failed", err)
	}
}

func TestRoleService_GetRolePermissions_Success(t *testing.T) {
	roleRepo := newMockRoleRepo()
	permRepo := newMockPermRepo()
	svc := NewRoleService(roleRepo, permRepo, logger.GetLogger("rbac_test"))
	ctx := context.Background()

	role, _ := svc.Create(ctx, CreateRoleInput{RoleCode: "admin", RoleName: "Admin"})
	permRepo.seedPermissions([]model.Permission{
		{ID: 1, PermissionCode: "user:read", PermissionName: "读", PermissionType: "READ", ResourceCode: "user", GroupCode: "USER"},
		{ID: 2, PermissionCode: "user:write", PermissionName: "写", PermissionType: "WRITE", ResourceCode: "user", GroupCode: "USER"},
	})
	_ = permRepo.AssignToRole(ctx, role.ID, []uint{1})

	vo, err := svc.GetRolePermissions(ctx, role.ID)
	if err != nil {
		t.Fatalf("GetRolePermissions() err = %v", err)
	}
	if len(vo.FlatPermissions) != 1 {
		t.Errorf("len(FlatPermissions) = %d, want 1", len(vo.FlatPermissions))
	}
	if vo.FlatPermissions[0] != "user:read" {
		t.Errorf("FlatPermissions[0] = %q, want user:read", vo.FlatPermissions[0])
	}

	selectedCount := 0
	for _, g := range vo.Groups {
		for _, p := range g.Permissions {
			if p.Selected {
				selectedCount++
			}
		}
	}
	if selectedCount != 1 {
		t.Errorf("selectedCount = %d, want 1", selectedCount)
	}
}

func TestRoleService_GetRolePermissions_RoleNotFound(t *testing.T) {
	svc := NewRoleService(newMockRoleRepo(), newMockPermRepo(), logger.GetLogger("rbac_test"))
	_, err := svc.GetRolePermissions(context.Background(), 999)
	if err != domainerrors.ErrRoleNotFound {
		t.Errorf("err = %v, want ErrRoleNotFound", err)
	}
}

func TestRoleService_AssignPermissions_Success(t *testing.T) {
	roleRepo := newMockRoleRepo()
	permRepo := newMockPermRepo()
	svc := NewRoleService(roleRepo, permRepo, logger.GetLogger("rbac_test"))
	ctx := context.Background()

	role, _ := svc.Create(ctx, CreateRoleInput{RoleCode: "a", RoleName: "A"})
	err := svc.AssignPermissions(ctx, role.ID, []uint{1, 2})
	if err != nil {
		t.Fatalf("AssignPermissions() err = %v", err)
	}
}

func TestRoleService_AssignPermissions_RoleNotFound(t *testing.T) {
	svc := NewRoleService(newMockRoleRepo(), newMockPermRepo(), logger.GetLogger("rbac_test"))
	err := svc.AssignPermissions(context.Background(), 999, []uint{1})
	if err != domainerrors.ErrRoleNotFound {
		t.Errorf("err = %v, want ErrRoleNotFound", err)
	}
}

func TestRoleService_AssignPermissions_AssignErr(t *testing.T) {
	roleRepo := newMockRoleRepo()
	permRepo := newMockPermRepo()
	permRepo.assignErr = errors.New("assign failed")
	svc := NewRoleService(roleRepo, permRepo, logger.GetLogger("rbac_test"))
	ctx := context.Background()
	role, _ := svc.Create(ctx, CreateRoleInput{RoleCode: "a", RoleName: "A"})
	err := svc.AssignPermissions(ctx, role.ID, []uint{1})
	if err == nil || err.Error() != "assign failed" {
		t.Errorf("err = %v, want assign failed", err)
	}
}

func TestGroupPermissions_MultipleGroups(t *testing.T) {
	perms := []model.Permission{
		{ID: 1, PermissionCode: "a:r", PermissionName: "AR", PermissionType: "READ", GroupCode: "A"},
		{ID: 2, PermissionCode: "b:r", PermissionName: "BR", PermissionType: "READ", GroupCode: "B"},
		{ID: 3, PermissionCode: "a:w", PermissionName: "AW", PermissionType: "WRITE", GroupCode: "A"},
	}
	groups := groupPermissions(perms, map[uint]bool{1: true})
	if len(groups) != 2 {
		t.Fatalf("len(groups) = %d, want 2", len(groups))
	}
	if groups[0].GroupCode != "A" || len(groups[0].Permissions) != 2 {
		t.Errorf("group A: %+v", groups[0])
	}
	if groups[1].GroupCode != "B" || len(groups[1].Permissions) != 1 {
		t.Errorf("group B: %+v", groups[1])
	}
	if !groups[0].Permissions[0].Selected {
		t.Error("permission ID 1 should be selected")
	}
	if groups[0].Permissions[1].Selected {
		t.Error("permission ID 3 should not be selected")
	}
}

func TestGroupPermissions_Empty(t *testing.T) {
	groups := groupPermissions(nil, nil)
	if len(groups) != 0 {
		t.Errorf("len(groups) = %d, want 0", len(groups))
	}
}

func TestRoleService_GetRolePermissions_FindAllErr(t *testing.T) {
	roleRepo := newMockRoleRepo()
	permRepo := newMockPermRepo()
	svc := NewRoleService(roleRepo, permRepo, logger.GetLogger("rbac_test"))
	ctx := context.Background()

	role, _ := svc.Create(ctx, CreateRoleInput{RoleCode: "a", RoleName: "A"})
	permRepo.findAllErr = errors.New("find all failed")
	_, err := svc.GetRolePermissions(ctx, role.ID)
	if err == nil || err.Error() != "find all failed" {
		t.Errorf("err = %v, want find all failed", err)
	}
}

func TestRoleService_GetRolePermissions_FindByRoleErr(t *testing.T) {
	roleRepo := newMockRoleRepo()
	permRepo := newMockPermRepo()
	svc := NewRoleService(roleRepo, permRepo, logger.GetLogger("rbac_test"))
	ctx := context.Background()

	role, _ := svc.Create(ctx, CreateRoleInput{RoleCode: "a", RoleName: "A"})
	permRepo.findByRoleErr = errors.New("find by role failed")
	_, err := svc.GetRolePermissions(ctx, role.ID)
	if err == nil || err.Error() != "find by role failed" {
		t.Errorf("err = %v, want find by role failed", err)
	}
}

func TestRoleService_GetRolePermissions_FindByIDErr(t *testing.T) {
	roleRepo := newMockRoleRepo()
	roleRepo.findByIDErr = errors.New("find failed")
	svc := NewRoleService(roleRepo, newMockPermRepo(), logger.GetLogger("rbac_test"))
	_, err := svc.GetRolePermissions(context.Background(), 1)
	if err == nil || err.Error() != "find failed" {
		t.Errorf("err = %v, want find failed", err)
	}
}

func TestRoleService_AssignPermissions_FindByIDErr(t *testing.T) {
	roleRepo := newMockRoleRepo()
	roleRepo.findByIDErr = errors.New("find failed")
	svc := NewRoleService(roleRepo, newMockPermRepo(), logger.GetLogger("rbac_test"))
	err := svc.AssignPermissions(context.Background(), 1, []uint{1})
	if err == nil || err.Error() != "find failed" {
		t.Errorf("err = %v, want find failed", err)
	}
}

func TestRoleService_GetByID_FindByIDErr(t *testing.T) {
	roleRepo := newMockRoleRepo()
	roleRepo.findByIDErr = errors.New("find failed")
	svc := NewRoleService(roleRepo, newMockPermRepo(), logger.GetLogger("rbac_test"))
	_, err := svc.GetByID(context.Background(), 1)
	if err == nil || err.Error() != "find failed" {
		t.Errorf("err = %v, want find failed", err)
	}
}

func TestGroupPermissions_UnknownGroupName(t *testing.T) {
	perms := []model.Permission{
		{ID: 1, PermissionCode: "x:r", PermissionName: "XR", PermissionType: "READ", GroupCode: "UNKNOWN"},
	}
	groups := groupPermissions(perms, nil)
	if len(groups) != 1 {
		t.Fatalf("len(groups) = %d, want 1", len(groups))
	}
	if groups[0].GroupName != "UNKNOWN" {
		t.Errorf("GroupName = %q, want UNKNOWN (fallback)", groups[0].GroupName)
	}
}
