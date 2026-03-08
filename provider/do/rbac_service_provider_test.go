package do

import (
	"context"
	"testing"

	"github.com/KOMKZ/go-yogan-domain-rbac/model"
	"github.com/KOMKZ/go-yogan-domain-rbac/repository"
	"github.com/KOMKZ/go-yogan-domain-rbac/service"
	"github.com/KOMKZ/go-yogan-framework/logger"
	"github.com/samber/do/v2"
)

type stubRoleRepo struct{}

func (stubRoleRepo) Create(_ context.Context, _ *model.Role) error             { return nil }
func (stubRoleRepo) Update(_ context.Context, _ *model.Role) error             { return nil }
func (stubRoleRepo) Delete(_ context.Context, _ uint) error                    { return nil }
func (stubRoleRepo) FindByID(_ context.Context, _ uint) (*model.Role, error)   { return nil, nil }
func (stubRoleRepo) FindByCode(_ context.Context, _ string) (*model.Role, error) { return nil, nil }
func (stubRoleRepo) FindByIDs(_ context.Context, _ []uint) ([]model.Role, error) { return nil, nil }
func (stubRoleRepo) Paginate(_ context.Context, _, _ int, _ string) ([]model.Role, int64, error) {
	return nil, 0, nil
}
func (stubRoleRepo) Count(_ context.Context) (int64, error) { return 0, nil }

type stubPermRepo struct{}

func (stubPermRepo) FindAll(_ context.Context) ([]model.Permission, error)          { return nil, nil }
func (stubPermRepo) FindByIDs(_ context.Context, _ []uint) ([]model.Permission, error) { return nil, nil }
func (stubPermRepo) FindByRoleID(_ context.Context, _ uint) ([]model.Permission, error) { return nil, nil }
func (stubPermRepo) AssignToRole(_ context.Context, _ uint, _ []uint) error         { return nil }

type stubUserRoleRepo struct{}

func (stubUserRoleRepo) FindRolesByUserID(_ context.Context, _ uint) ([]model.Role, error) { return nil, nil }
func (stubUserRoleRepo) FindRoleIDsByUserID(_ context.Context, _ uint) ([]uint, error) { return nil, nil }
func (stubUserRoleRepo) AssignRolesToUser(_ context.Context, _ uint, _ []uint) error { return nil }
func (stubUserRoleRepo) FindPermissionsByUserID(_ context.Context, _ uint) ([]model.Permission, error) { return nil, nil }
func (stubUserRoleRepo) FindPermissionCodesByUserID(_ context.Context, _ uint) ([]string, error) { return nil, nil }
func (stubUserRoleRepo) HasRole(_ context.Context, _ uint, _ string) (bool, error) { return false, nil }
func (stubUserRoleRepo) HasPermission(_ context.Context, _ uint, _ string) (bool, error) { return false, nil }

var (
	_ repository.RoleRepository     = stubRoleRepo{}
	_ repository.PermissionRepository = stubPermRepo{}
	_ repository.UserRoleRepository  = stubUserRoleRepo{}
)

func TestProvideRoleService(t *testing.T) {
	injector := do.New()
	do.ProvideNamedValue(injector, "rbac", logger.GetLogger("rbac_test"))
	do.Provide(injector, func(_ do.Injector) (repository.RoleRepository, error) {
		return stubRoleRepo{}, nil
	})
	do.Provide(injector, func(_ do.Injector) (repository.PermissionRepository, error) {
		return stubPermRepo{}, nil
	})
	do.Provide(injector, ProvideRoleService)

	svc, err := do.Invoke[*service.RoleService](injector)
	if err != nil {
		t.Fatalf("ProvideRoleService err = %v", err)
	}
	if svc == nil {
		t.Fatal("RoleService should not be nil")
	}
}

func TestProvideUserRoleService(t *testing.T) {
	injector := do.New()
	do.ProvideNamedValue(injector, "rbac", logger.GetLogger("rbac_test"))
	do.Provide(injector, func(_ do.Injector) (repository.RoleRepository, error) {
		return stubRoleRepo{}, nil
	})
	do.Provide(injector, func(_ do.Injector) (repository.UserRoleRepository, error) {
		return stubUserRoleRepo{}, nil
	})
	do.Provide(injector, ProvideUserRoleService)

	svc, err := do.Invoke[*service.UserRoleService](injector)
	if err != nil {
		t.Fatalf("ProvideUserRoleService err = %v", err)
	}
	if svc == nil {
		t.Fatal("UserRoleService should not be nil")
	}
}

func TestProvideRoleService_MissingRoleRepo(t *testing.T) {
	injector := do.New()
	do.ProvideNamedValue(injector, "rbac", logger.GetLogger("rbac_test"))
	do.Provide(injector, func(_ do.Injector) (repository.PermissionRepository, error) {
		return stubPermRepo{}, nil
	})
	do.Provide(injector, ProvideRoleService)

	_, err := do.Invoke[*service.RoleService](injector)
	if err == nil {
		t.Fatal("should fail without RoleRepository")
	}
}

func TestProvideRoleService_MissingPermRepo(t *testing.T) {
	injector := do.New()
	do.ProvideNamedValue(injector, "rbac", logger.GetLogger("rbac_test"))
	do.Provide(injector, func(_ do.Injector) (repository.RoleRepository, error) {
		return stubRoleRepo{}, nil
	})
	do.Provide(injector, ProvideRoleService)

	_, err := do.Invoke[*service.RoleService](injector)
	if err == nil {
		t.Fatal("should fail without PermissionRepository")
	}
}

func TestProvideUserRoleService_MissingRoleRepo(t *testing.T) {
	injector := do.New()
	do.ProvideNamedValue(injector, "rbac", logger.GetLogger("rbac_test"))
	do.Provide(injector, func(_ do.Injector) (repository.UserRoleRepository, error) {
		return stubUserRoleRepo{}, nil
	})
	do.Provide(injector, ProvideUserRoleService)

	_, err := do.Invoke[*service.UserRoleService](injector)
	if err == nil {
		t.Fatal("should fail without RoleRepository")
	}
}

func TestProvideUserRoleService_MissingUserRoleRepo(t *testing.T) {
	injector := do.New()
	do.ProvideNamedValue(injector, "rbac", logger.GetLogger("rbac_test"))
	do.Provide(injector, func(_ do.Injector) (repository.RoleRepository, error) {
		return stubRoleRepo{}, nil
	})
	do.Provide(injector, ProvideUserRoleService)

	_, err := do.Invoke[*service.UserRoleService](injector)
	if err == nil {
		t.Fatal("should fail without UserRoleRepository")
	}
}
