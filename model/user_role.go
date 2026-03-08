package model

import "time"

type UserRole struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	RoleID    uint      `gorm:"not null;index" json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (UserRole) TableName() string {
	return "sys_user_roles"
}
