package model

import "testing"

func TestPermission_TableName(t *testing.T) {
	var p Permission
	if got := p.TableName(); got != "sys_permissions" {
		t.Errorf("TableName() = %q, want %q", got, "sys_permissions")
	}
}
