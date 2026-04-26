package repository

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"sparklink-backend/model"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByPhone(phone string) (*model.User, error) {
	var user model.User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByDeviceID(deviceID string) (*model.User, error) {
	var user model.User
	err := r.db.Where("device_id = ?", deviceID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Save(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) AddCredits(userID uint, minutes int) error {
	return r.db.Model(&model.User{}).
		Where("id = ?", userID).
		Update("ad_credits", gorm.Expr("ad_credits + ?", minutes)).Error
}

func (r *UserRepository) SetExpireTime(userID uint, expireTime time.Time) error {
	return r.db.Model(&model.User{}).
		Where("id = ?", userID).
		Update("expire_time", expireTime).Error
}

// Device operations
func (r *UserRepository) CreateDevice(device *model.Device) error {
	return r.db.Create(device).Error
}

func (r *UserRepository) FindDevicesByUserID(userID uint) ([]model.Device, error) {
	var devices []model.Device
	err := r.db.Where("user_id = ? AND is_active = ?", userID, true).Find(&devices).Error
	return devices, err
}

func (r *UserRepository) DeactivateDevice(deviceID string) error {
	return r.db.Model(&model.Device{}).
		Where("device_id = ?", deviceID).
		Update("is_active", false).Error
}

type NodeRepository struct {
	db *gorm.DB
}

func NewNodeRepository(db *gorm.DB) *NodeRepository {
	return &NodeRepository{db: db}
}

func (r *NodeRepository) FindAll(protocol, nodeType, country string) ([]model.Node, error) {
	var nodes []model.Node
	query := r.db.Where("status = ?", "online")

	if protocol != "" {
		query = query.Where("protocol = ?", protocol)
	}
	if nodeType != "" {
		query = query.Where("node_type = ?", nodeType)
	}
	if country != "" {
		query = query.Where("country = ?", country)
	}

	err := query.Order("latency ASC").Find(&nodes).Error
	return nodes, err
}

func (r *NodeRepository) FindByID(id uint) (*model.Node, error) {
	var node model.Node
	err := r.db.First(&node, id).Error
	if err != nil {
		return nil, err
	}
	return &node, nil
}

func (r *NodeRepository) UpdateLatency(id uint, latency int) error {
	return r.db.Model(&model.Node{}).
		Where("id = ?", id).
		Update("latency", latency).Error
}

func (r *NodeRepository) UpdateLoad(id uint, load int) error {
	return r.db.Model(&model.Node{}).
		Where("id = ?", id).
		Update("load", load).Error
}

type RewardRepository struct {
	db *gorm.DB
	rdb *redis.Client
}

func NewRewardRepository(db *gorm.DB) *RewardRepository {
	return &RewardRepository{db: db, rdb: nil}
}

func (r *RewardRepository) SetRedis(rdb *redis.Client) {
	r.rdb = rdb
}

func (r *RewardRepository) CreateAdLog(log *model.AdLog) error {
	return r.db.Create(log).Error
}

func (r *RewardRepository) FindAdLogByTransactionID(txnID string) (*model.AdLog, error) {
	var log model.AdLog
	err := r.db.Where("transaction_id = ?", txnID).First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *RewardRepository) CreateDailyCheckin(checkin *model.DailyCheckin) error {
	return r.db.Create(checkin).Error
}

func (r *RewardRepository) FindDailyCheckin(userID uint, date time.Time) (*model.DailyCheckin, error) {
	var checkin model.DailyCheckin
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	err := r.db.Where("user_id = ? AND checkin_date = ?", userID, startOfDay).First(&checkin).Error
	if err != nil {
		return nil, err
	}
	return &checkin, nil
}

func (r *RewardRepository) CreateInvite(invite *model.Invite) error {
	return r.db.Create(invite).Error
}

func (r *RewardRepository) IsNonceUsed(nonce string) (bool, error) {
	if r.rdb == nil {
		return false, nil
	}
	ctx := context.Background()
	return r.db.Exists(ctx, "nonce:"+nonce).Result()
}

func (r *RewardRepository) SetNonce(nonce string, expire time.Duration) error {
	if r.rdb == nil {
		return nil
	}
	ctx := context.Background()
	return r.db.SetEX(ctx, "nonce:"+nonce, "1", expire).Err()
}

type SubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(sub *model.Subscription) error {
	return r.db.Create(sub).Error
}

func (r *SubscriptionRepository) FindByUserID(userID uint) ([]model.Subscription, error) {
	var subs []model.Subscription
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&subs).Error
	return subs, err
}

func (r *SubscriptionRepository) FindActiveByUserID(userID uint) (*model.Subscription, error) {
	var sub model.Subscription
	err := r.db.Where("user_id = ? AND status = ? AND expire_time > ?", userID, "active", time.Now()).First(&sub).Error
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *SubscriptionRepository) FindAll() ([]model.Plan, error) {
	var plans []model.Plan
	err := r.db.Find(&plans).Error
	return plans, err
}