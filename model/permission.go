package model

import "time"

type Permission struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	PermissionCode string    `gorm:"size:100;uniqueIndex;not null" json:"permission_code"`
	PermissionName string    `gorm:"size:100;not null" json:"permission_name"`
	PermissionType string    `gorm:"size:10;not null" json:"permission_type"` // READ, WRITE
	ResourceCode   string    `gorm:"size:50;not null" json:"resource_code"`
	GroupCode      string    `gorm:"size:50;not null" json:"group_code"`
	Description    string    `gorm:"size:200" json:"description"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (Permission) TableName() string {
	return "sys_permissions"
}
