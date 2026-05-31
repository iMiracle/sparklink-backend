package repository

import (
	"time"

	"sparklink-backend/model"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByID(id uint) (*model.User, error)
	FindByPhone(phone string) (*model.User, error)
	FindByDeviceID(deviceID string) (*model.User, error)
	FindByInviteCode(code string) (*model.User, error)
	Save(user *model.User) error
	AddBalance(userID uint, minutes int) error
	CountDevicesByUserID(userID uint) (int, error)
	FindDevicesByUserID(userID uint) ([]model.Device, error)
	CreateDevice(device *model.Device) error
	DeactivateDevice(deviceID string) error
}

type NodeRepository interface {
	FindAll(protocol, visibility, region string) ([]model.Node, error)
	FindByNodeID(nodeID string) (*model.Node, error)
	FindByID(id uint) (*model.Node, error)
	UpdatePing(nodeID string, latency int) error
	UpdateLatency(id uint, latency int) error
	UpdateLoad(id uint, load int) error
	Create(node *model.Node) error
	UpdateByNodeID(nodeID string, updates map[string]interface{}) error
	AddFavorite(userID uint, nodeID string) error
	RemoveFavorite(userID uint, nodeID string) error
	GetFavorites(userID uint) ([]model.Favorite, error)
	GetRegions() ([]string, error)
}

type RewardRepository interface {
	CreateAdLog(log *model.AdLog) error
	FindAdLogByTransactionID(txnID string) (*model.AdLog, error)
	FindRecentAdLog(userID uint, adType string) (*model.AdLog, error)
	CreateDailyCheckin(checkin *model.DailyCheckin) error
	FindDailyCheckin(userID uint, date time.Time) (*model.DailyCheckin, error)
	CreateInvite(invite *model.Invite) error
	IsNonceUsed(nonce string) (bool, error)
	SetNonce(nonce string, expire time.Duration) error
}

type SubscriptionRepository interface {
	FindAllPlans() ([]model.Plan, error)
	FindPlanByID(planID string) (*model.Plan, error)
	CreatePlan(plan *model.Plan) error
	UpdatePlan(planID string, updates map[string]interface{}) error
	DeletePlan(planID string) error
	Create(sub *model.Subscription) error
	FindActiveByUserID(userID uint) (*model.Subscription, error)
}

type ConnectRepository interface {
	CreateSession(session *model.ConnectSession) error
	FindActiveSession(userID uint) (*model.ConnectSession, error)
	UpdateSession(session *model.ConnectSession) error
}

type VerificationRepository interface {
	Create(code *model.VerificationCode) error
	FindValidCode(phone, code string) (*model.VerificationCode, error)
	MarkUsed(id uint) error
	CreateQRSession(session *model.QRSession) error
	FindQRSessionBySessionID(sessionID string) (*model.QRSession, error)
	UpdateQRSessionStatus(sessionID string, status string, userID *uint, token string) error
}

type AdminRepository interface {
	Create(admin *model.AdminUser) error
	FindByID(id uint) (*model.AdminUser, error)
	FindByUsername(username string) (*model.AdminUser, error)
	FindAll() ([]model.AdminUser, error)
	UpdateStatus(id uint, status string) error
	Delete(id uint) error
}

type AnnouncementRepository interface {
	Create(a *model.Announcement) error
	FindByID(id uint) (*model.Announcement, error)
	FindAll() ([]model.Announcement, error)
	Update(a *model.Announcement) error
	Delete(id uint) error
}

type AuditLogRepository interface {
	Create(log *model.AuditLog) error
	FindAll(page, pageSize int) ([]model.AuditLog, int64, error)
}

var _ UserRepository = (*GormUserRepository)(nil)
var _ NodeRepository = (*GormNodeRepository)(nil)
var _ RewardRepository = (*GormRewardRepository)(nil)
var _ SubscriptionRepository = (*GormSubscriptionRepository)(nil)
var _ ConnectRepository = (*GormConnectRepository)(nil)
var _ VerificationRepository = (*GormVerificationRepository)(nil)
