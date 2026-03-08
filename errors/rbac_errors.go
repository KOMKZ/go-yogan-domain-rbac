package errors

import (
	"net/http"

	"github.com/KOMKZ/go-yogan-framework/errcode"
)

const ModuleRBAC = 26

var (
	ErrRoleNotFound = errcode.Register(errcode.New(
		ModuleRBAC, 1001, "rbac",
		"error.rbac.role_not_found", "角色不存在",
		http.StatusNotFound,
	))
	ErrRoleCodeExists = errcode.Register(errcode.New(
		ModuleRBAC, 1002, "rbac",
		"error.rbac.role_code_exists", "角色编码已存在",
		http.StatusConflict,
	))
	ErrPermissionNotFound = errcode.Register(errcode.New(
		ModuleRBAC, 1003, "rbac",
		"error.rbac.permission_not_found", "权限不存在",
		http.StatusNotFound,
	))
	ErrRoleDisabled = errcode.Register(errcode.New(
		ModuleRBAC, 1004, "rbac",
		"error.rbac.role_disabled", "角色已禁用",
		http.StatusBadRequest,
	))
	ErrUserRoleRepoNil = errcode.Register(errcode.New(
		ModuleRBAC, 1005, "rbac",
		"error.rbac.user_role_repo_nil", "用户角色仓库未初始化",
		http.StatusInternalServerError,
	))
)
