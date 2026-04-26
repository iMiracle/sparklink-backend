package model

import (
	"time"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex" json:"email"`
	Phone        string    `gorm:"index" json:"phone"`
	Password    string    `json:"-"`
	Nickname     string    `json:"nickname"`
	DeviceID    string    `gorm:"index" json:"device_id"`
	ExpireTime  *time.Time `json:"expire_time"`
	AdCredits   int       `gorm:"default:0" json:"ad_credits"`
	ReferralCode string   `gorm:"uniqueIndex" json:"referral_code"`
	ReferredBy  *uint    `json:"referred_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Device struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index" json:"user_id"`
	DeviceID   string    `gorm:"uniqueIndex" json:"device_id"`
	DeviceName string    `json:"device_name"`
	Platform   string    `json:"platform"`
	LastLogin  time.Time `json:"last_login"`
	IsActive   bool      `gorm:"default:true" json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
}

type Node struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Country     string    `json:"country"`
	City        string    `json:"city"`
	Protocol    string    `json:"protocol"`
	NodeType    string    `json:"node_type"`
	Host        string    `json:"host"`
	Port        int       `json:"port"`
	PublicKey  string    `json:"public_key,omitempty"`
	Password   string    `json:"password,omitempty"`
	Latency    int       `json:"latency"`
	Load       int       `json:"load"`
	Bandwidth  int64     `json:"bandwidth"`
	Status     string    `gorm:"default:online" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Subscription struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	Plan       string    `json:"plan"`
	Amount     float64   `json:"amount"`
	StartTime  time.Time `json:"start_time"`
	ExpireTime time.Time `json:"expire_time"`
	Status     string    `gorm:"default:active" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type AdLog struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `gorm:"index" json:"user_id"`
	AdPlatform     string    `json:"ad_platform"`
	AdID           string    `json:"ad_id"`
	TransactionID  string    `gorm:"uniqueIndex" json:"transaction_id"`
	RewardAmount  int       `json:"reward_amount"`
	Status        string    `gorm:"default:success" json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

type DailyCheckin struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	CheckinDate time.Time `gorm:"uniqueIndex:checkin_date" json:"checkin_date"`
	RewardAmount int      `json:"reward_amount"`
	CreatedAt  time.Time `json:"created_at"`
}

type Invite struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `gorm:"index" json:"user_id"`
	InvitedUserID   uint      `json:"invited_user_id"`
	RewardAmount    int       `json:"reward_amount"`
	CreatedAt       time.Time `json:"created_at"`
}

type Favorite struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	NodeID     uint      `json:"node_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type Plan struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Name       string   `json:"name"`
	Plan       string   `json:"plan"`
	Price      float64  `json:"price"`
	Duration   int      `json:"duration"`
	IsPopular  bool     `json:"is_popular"`
	Features   string   `json:"features"`
}