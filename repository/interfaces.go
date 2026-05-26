package repository

import (
	"sparklink-backend/model"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByID(id uint) (*model.User, error)
	FindByPhone(phone string) (*model.User, error)
	FindByInviteCode(code string) (*model.User, error)
	Save(user *model.User) error
	AddBalance(userID uint, minutes int) error
	FindDevicesByUserID(userID uint) ([]model.Device, error)
	DeactivateDevice(deviceID string) error
}

type NodeRepository interface {
	FindAll(protocol, visibility, region string) ([]model.Node, error)
	FindByNodeID(nodeID string) (*model.Node, error)
	UpdatePing(nodeID string, latency int) error
}

type RewardRepository interface {
	CreateAdLog(log *model.AdLog) error
	FindRecentAdLog(userID uint, adType string) (*model.AdLog, error)
}

type SubscriptionRepository interface {
	FindAllPlans() ([]model.Plan, error)
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
}

// ensure mock and gorm types implement interfaces at compile time
var _ UserRepository = (*MockUserRepository)(nil)
var _ NodeRepository = (*MockNodeRepository)(nil)
var _ RewardRepository = (*MockRewardRepository)(nil)
var _ SubscriptionRepository = (*MockSubscriptionRepository)(nil)
var _ ConnectRepository = (*MockConnectRepository)(nil)
var _ VerificationRepository = (*MockVerificationRepository)(nil)

var _ UserRepository = (*GormUserRepository)(nil)
var _ NodeRepository = (*GormNodeRepository)(nil)
var _ RewardRepository = (*GormRewardRepository)(nil)
var _ SubscriptionRepository = (*GormSubscriptionRepository)(nil)
var _ ConnectRepository = (*GormConnectRepository)(nil)
var _ VerificationRepository = (*GormVerificationRepository)(nil)
