package model

import "time"

type AdminUser struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	AdminID     string    `gorm:"uniqueIndex;size:32" json:"admin_id"`
	Username    string    `gorm:"uniqueIndex;size:64" json:"username"`
	Password    string    `gorm:"size:128" json:"-"`
	Role        string    `gorm:"size:16" json:"role"`
	Status      string    `gorm:"default:active;size:16" json:"status"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Announcement struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"size:256" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	Type      string    `gorm:"size:32" json:"type"`
	Target    string    `gorm:"size:32" json:"target"`
	Status    string    `gorm:"default:draft;size:16" json:"status"`
	StartAt   time.Time `json:"start_at"`
	EndAt     time.Time `json:"end_at"`
	CreatedBy string    `gorm:"size:64" json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuditLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	AdminID    string    `gorm:"index;size:32" json:"admin_id"`
	Action     string    `gorm:"size:64" json:"action"`
	TargetType string    `gorm:"size:32" json:"target_type"`
	TargetID   string    `gorm:"size:64" json:"target_id"`
	Detail     string    `gorm:"type:text" json:"detail"`
	IP         string    `gorm:"size:64" json:"ip"`
	Result     string    `gorm:"size:16" json:"result"`
	CreatedAt  time.Time `json:"created_at"`
}
