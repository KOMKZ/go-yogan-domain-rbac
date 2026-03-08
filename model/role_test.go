package model

import "testing"

func TestRole_TableName(t *testing.T) {
	var r Role
	if got := r.TableName(); got != "sys_roles" {
		t.Errorf("TableName() = %q, want %q", got, "sys_roles")
	}
}

func TestRole_IsEnabled(t *testing.T) {
	tests := []struct {
		name   string
		status int
		want   bool
	}{
		{"enabled when status 1", 1, true},
		{"disabled when status 0", 0, false},
		{"disabled when status 2", 2, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Role{Status: tt.status}
			if got := r.IsEnabled(); got != tt.want {
				t.Errorf("IsEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}
