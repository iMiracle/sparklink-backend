package model

import "time"

type User struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	Phone          string     `gorm:"uniqueIndex;size:32" json:"phone"`
	Nickname       string     `gorm:"size:64" json:"nickname"`
	Avatar         string     `gorm:"size:256" json:"avatar"`
	DeviceID       string     `gorm:"index;size:128" json:"device_id"`
	VipStatus      string     `gorm:"default:inactive;size:16" json:"vip_status"`
	VipExpireAt    *time.Time `json:"vip_expire_at,omitempty"`
	BalanceMinutes int        `gorm:"default:0" json:"balance_minutes"`
	InviteCode     string     `gorm:"uniqueIndex;size:32" json:"invite_code"`
	InvitedCount   int        `gorm:"default:0" json:"invited_count"`
	AdCredits      int        `gorm:"default:0" json:"ad_credits"`
	ExpireTime     *time.Time `json:"expire_time,omitempty"`
	ReferredBy     *uint      `json:"referred_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type VerificationCode struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Phone     string    `gorm:"index;size:32" json:"phone"`
	Code      string    `gorm:"size:8" json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `gorm:"default:false" json:"used"`
	CreatedAt time.Time `json:"created_at"`
}

type QRSession struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SessionID string    `gorm:"uniqueIndex;size:64" json:"session_id"`
	Status    string    `gorm:"default:pending" json:"status"` // pending / scanned / confirmed / expired
	UserID    *uint     `json:"user_id,omitempty"`
	Token     string    `json:"token,omitempty"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Device struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	DeviceID   string    `gorm:"uniqueIndex;size:128" json:"device_id"`
	DeviceName string    `json:"device_name"`
	DeviceType string    `json:"device_type"`
	IsActive   bool      `gorm:"default:true" json:"is_active"`
	LastActive time.Time `json:"last_active"`
	CreatedAt  time.Time `json:"created_at"`
}

type Node struct {
	ID              uint   `gorm:"primaryKey" json:"-"`
	NodeId          string `gorm:"uniqueIndex;size:32" json:"node_id"`
	Name            string `gorm:"size:64" json:"name"`
	Protocol        string `gorm:"size:16" json:"protocol"`
	Latency         int    `json:"latency"`
	Load            int    `json:"load"`
	RegionCode      string `json:"region_code"`
	RegionName      string `json:"region_name"`
	Tags            string `gorm:"size:128" json:"tags"` // comma-separated
	VisibilityLevel string `gorm:"default:free;size:16" json:"visibility_level"`
	Priority        int    `gorm:"default:0" json:"priority"`
	Host            string `gorm:"size:128" json:"host,omitempty"`
	Port            int    `json:"port,omitempty"`
	PublicKey       string `gorm:"size:256" json:"public_key,omitempty"`
	Protocols       string `gorm:"size:64" json:"protocols,omitempty"` // comma-separated
	Distance        int    `json:"distance,omitempty"`
	BandwidthLimit  int    `json:"bandwidth_limit,omitempty"`
	Status          string `gorm:"default:online;size:16" json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type ConnectSession struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	SessionID     string    `gorm:"uniqueIndex;size:64" json:"session_id"`
	UserID        uint      `gorm:"index" json:"user_id"`
	NodeID        string    `gorm:"size:32" json:"node_id"`
	Protocol      string    `gorm:"size:16" json:"protocol"`
	Status        string    `gorm:"default:active;size:16" json:"status"` // active / stopped
	StartedAt     time.Time `json:"started_at"`
	StoppedAt     *time.Time `json:"stopped_at,omitempty"`
	BytesSent     int64     `gorm:"default:0" json:"bytes_sent"`
	BytesReceived int64     `gorm:"default:0" json:"bytes_received"`
	CreatedAt     time.Time `json:"created_at"`
}

type Subscription struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	PlanID     string    `gorm:"size:32" json:"plan_id"`
	Amount     float64   `json:"amount"`
	StartTime  time.Time `json:"start_time"`
	ExpireTime time.Time `json:"expire_time"`
	Status     string    `gorm:"default:active;size:16" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type AdLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	AdID       string    `gorm:"size:64" json:"ad_id"`
	AdType     string    `gorm:"index;size:16" json:"ad_type"` // reward_video / sign
	Reward     int       `json:"reward"`
	TransactionID string    `json:"transaction_id"`
	CooldownEnd   time.Time `json:"cooldown_end"`
	CreatedAt     time.Time `json:"created_at"`
}

type DailyCheckin struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index" json:"user_id"`
	CheckinDate time.Time `json:"checkin_date"`
	Reward      int       `json:"reward"`
	CreatedAt   time.Time `json:"created_at"`
}

type Invite struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	InviteCode string    `gorm:"index;size:32" json:"invite_code"`
	InvitedUID uint      `json:"invited_uid"`
	Reward     int       `json:"reward"`
	CreatedAt  time.Time `json:"created_at"`
}

type Plan struct {
	ID            uint     `gorm:"primaryKey" json:"-"`
	PlanID        string   `gorm:"uniqueIndex;size:32" json:"plan_id"`
	Name          string   `gorm:"size:64" json:"name"`
	Price         float64  `json:"price"`
	OriginalPrice float64  `json:"original_price"`
	Currency      string   `gorm:"default:USD;size:8" json:"currency"`
	DurationDays  int      `json:"duration_days"`
	DailyPrice    *float64 `json:"daily_price,omitempty"`
	Popular       bool     `json:"popular"`
	Tag           string   `gorm:"default:''" json:"tag"`
	Features      string   `json:"features,omitempty"`
}
