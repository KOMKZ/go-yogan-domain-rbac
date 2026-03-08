package model

import "testing"

func TestUserRole_TableName(t *testing.T) {
	var ur UserRole
	if got := ur.TableName(); got != "sys_user_roles" {
		t.Errorf("TableName() = %q, want %q", got, "sys_user_roles")
	}
}
