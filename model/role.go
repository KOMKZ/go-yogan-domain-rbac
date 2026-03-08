package model

import "time"

type Role struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	RoleCode    string    `gorm:"size:50;uniqueIndex;not null" json:"role_code"`
	RoleName    string    `gorm:"size:100;not null" json:"role_name"`
	Description string    `gorm:"size:200" json:"description"`
	Status      int       `gorm:"default:1" json:"status"` // 1=启用, 0=禁用
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Role) TableName() string {
	return "sys_roles"
}

func (r *Role) IsEnabled() bool {
	return r.Status == 1
}
