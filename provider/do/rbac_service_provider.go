package do

import (
	"github.com/KOMKZ/go-yogan-domain-rbac/repository"
	"github.com/KOMKZ/go-yogan-domain-rbac/service"
	"github.com/KOMKZ/go-yogan-framework/logger"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

// ---- Repository Providers ----

func ProvideRoleRepository(i do.Injector) (repository.RoleRepository, error) {
	db, err := do.Invoke[*gorm.DB](i)
	if err != nil {
		return nil, err
	}
	return repository.NewRoleMySQLRepository(db), nil
}

func ProvidePermissionRepository(i do.Injector) (repository.PermissionRepository, error) {
	db, err := do.Invoke[*gorm.DB](i)
	if err != nil {
		return nil, err
	}
	return repository.NewPermissionMySQLRepository(db), nil
}

func ProvideUserRoleRepository(i do.Injector) (repository.UserRoleRepository, error) {
	db, err := do.Invoke[*gorm.DB](i)
	if err != nil {
		return nil, err
	}
	return repository.NewUserRoleMySQLRepository(db), nil
}

// ---- Service Providers ----

func ProvideRoleService(i do.Injector) (*service.RoleService, error) {
	roleRepo, err := do.Invoke[repository.RoleRepository](i)
	if err != nil {
		return nil, err
	}
	permRepo, err := do.Invoke[repository.PermissionRepository](i)
	if err != nil {
		return nil, err
	}
	log := do.MustInvokeNamed[*logger.CtxZapLogger](i, "rbac")
	return service.NewRoleService(roleRepo, permRepo, log), nil
}

func ProvideUserRoleService(i do.Injector) (*service.UserRoleService, error) {
	roleRepo, err := do.Invoke[repository.RoleRepository](i)
	if err != nil {
		return nil, err
	}
	userRoleRepo, err := do.Invoke[repository.UserRoleRepository](i)
	if err != nil {
		return nil, err
	}
	log := do.MustInvokeNamed[*logger.CtxZapLogger](i, "rbac")
	return service.NewUserRoleService(roleRepo, userRoleRepo, log), nil
}
