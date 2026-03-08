package model

import "testing"

func TestRolePermission_TableName(t *testing.T) {
	var rp RolePermission
	if got := rp.TableName(); got != "sys_role_permissions" {
		t.Errorf("TableName() = %q, want %q", got, "sys_role_permissions")
	}
}
