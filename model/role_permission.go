package model

import "time"

type RolePermission struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	RoleID       uint      `gorm:"not null;index" json:"role_id"`
	PermissionID uint      `gorm:"not null;index" json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
}

func (RolePermission) TableName() string {
	return "sys_role_permissions"
}
