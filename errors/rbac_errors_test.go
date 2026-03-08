package errors

import (
	"errors"
	"testing"
)

func TestErrRoleNotFound(t *testing.T) {
	if ErrRoleNotFound == nil {
		t.Fatal("ErrRoleNotFound should not be nil")
	}
	if !errors.Is(ErrRoleNotFound, ErrRoleNotFound) {
		t.Error("errors.Is should match ErrRoleNotFound")
	}
	if ErrRoleNotFound.Error() != "role not found" {
		t.Errorf("Error() = %q, want %q", ErrRoleNotFound.Error(), "role not found")
	}
}

func TestErrRoleCodeExists(t *testing.T) {
	if ErrRoleCodeExists == nil {
		t.Fatal("ErrRoleCodeExists should not be nil")
	}
	if ErrRoleCodeExists.Error() != "role code already exists" {
		t.Errorf("Error() = %q, want %q", ErrRoleCodeExists.Error(), "role code already exists")
	}
}

func TestErrPermissionNotFound(t *testing.T) {
	if ErrPermissionNotFound == nil {
		t.Fatal("ErrPermissionNotFound should not be nil")
	}
	if ErrPermissionNotFound.Error() != "permission not found" {
		t.Errorf("Error() = %q, want %q", ErrPermissionNotFound.Error(), "permission not found")
	}
}

func TestErrRoleDisabled(t *testing.T) {
	if ErrRoleDisabled == nil {
		t.Fatal("ErrRoleDisabled should not be nil")
	}
	if ErrRoleDisabled.Error() != "role is disabled" {
		t.Errorf("Error() = %q, want %q", ErrRoleDisabled.Error(), "role is disabled")
	}
}

func TestErrUserRoleRepoNil(t *testing.T) {
	if ErrUserRoleRepoNil == nil {
		t.Fatal("ErrUserRoleRepoNil should not be nil")
	}
	if ErrUserRoleRepoNil.Error() != "user role repository not initialized" {
		t.Errorf("Error() = %q, want %q", ErrUserRoleRepoNil.Error(), "user role repository not initialized")
	}
}
