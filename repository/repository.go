package repository

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"sparklink-backend/model"
)

<<<<<<< HEAD
type GormUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *GormUserRepository) FindByID(id uint) (*model.User, error) {
=======
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
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

<<<<<<< HEAD
func (r *GormUserRepository) FindByEmail(email string) (*model.User, error) {
=======
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

<<<<<<< HEAD
func (r *GormUserRepository) FindByPhone(phone string) (*model.User, error) {
=======
func (r *UserRepository) FindByPhone(phone string) (*model.User, error) {
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	var user model.User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

<<<<<<< HEAD
func (r *GormUserRepository) FindByDeviceID(deviceID string) (*model.User, error) {
=======
func (r *UserRepository) FindByDeviceID(deviceID string) (*model.User, error) {
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	var user model.User
	err := r.db.Where("device_id = ?", deviceID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

<<<<<<< HEAD
func (r *GormUserRepository) Save(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *GormUserRepository) AddCredits(userID uint, minutes int) error {
	return r.db.Model(&model.User{}).
		Where("id = ?", userID).
		Update("balance_minutes", gorm.Expr("balance_minutes + ?", minutes)).Error
}

func (r *GormUserRepository) SetExpireTime(userID uint, expireTime time.Time) error {
=======
func (r *UserRepository) Save(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) AddCredits(userID uint, minutes int) error {
	return r.db.Model(&model.User{}).
		Where("id = ?", userID).
		Update("ad_credits", gorm.Expr("ad_credits + ?", minutes)).Error
}

func (r *UserRepository) SetExpireTime(userID uint, expireTime time.Time) error {
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	return r.db.Model(&model.User{}).
		Where("id = ?", userID).
		Update("expire_time", expireTime).Error
}

<<<<<<< HEAD
func (r *GormUserRepository) FindByInviteCode(code string) (*model.User, error) {
	var user model.User
	err := r.db.Where("invite_code = ?", code).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) AddBalance(userID uint, minutes int) error {
	return r.db.Model(&model.User{}).
		Where("id = ?", userID).
		Update("balance_minutes", gorm.Expr("balance_minutes + ?", minutes)).Error
}

// Device operations
func (r *GormUserRepository) CreateDevice(device *model.Device) error {
	return r.db.Create(device).Error
}

func (r *GormUserRepository) FindDevicesByUserID(userID uint) ([]model.Device, error) {
=======
// Device operations
func (r *UserRepository) CreateDevice(device *model.Device) error {
	return r.db.Create(device).Error
}

func (r *UserRepository) FindDevicesByUserID(userID uint) ([]model.Device, error) {
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	var devices []model.Device
	err := r.db.Where("user_id = ? AND is_active = ?", userID, true).Find(&devices).Error
	return devices, err
}

<<<<<<< HEAD
func (r *GormUserRepository) DeactivateDevice(deviceID string) error {
=======
func (r *UserRepository) DeactivateDevice(deviceID string) error {
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	return r.db.Model(&model.Device{}).
		Where("device_id = ?", deviceID).
		Update("is_active", false).Error
}

<<<<<<< HEAD
type GormNodeRepository struct {
	db *gorm.DB
}

func NewNodeRepository(db *gorm.DB) *GormNodeRepository {
	return &GormNodeRepository{db: db}
}

func (r *GormNodeRepository) FindAll(protocol, visibility, region string) ([]model.Node, error) {
=======
type NodeRepository struct {
	db *gorm.DB
}

func NewNodeRepository(db *gorm.DB) *NodeRepository {
	return &NodeRepository{db: db}
}

func (r *NodeRepository) FindAll(protocol, nodeType, country string) ([]model.Node, error) {
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	var nodes []model.Node
	query := r.db.Where("status = ?", "online")

	if protocol != "" {
		query = query.Where("protocol = ?", protocol)
	}
<<<<<<< HEAD
	if visibility != "" {
		query = query.Where("visibility_level = ?", visibility)
	}
	if region != "" {
		query = query.Where("region_code = ?", region)
=======
	if nodeType != "" {
		query = query.Where("node_type = ?", nodeType)
	}
	if country != "" {
		query = query.Where("country = ?", country)
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	}

	err := query.Order("latency ASC").Find(&nodes).Error
	return nodes, err
}

<<<<<<< HEAD
func (r *GormNodeRepository) FindByID(id uint) (*model.Node, error) {
=======
func (r *NodeRepository) FindByID(id uint) (*model.Node, error) {
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	var node model.Node
	err := r.db.First(&node, id).Error
	if err != nil {
		return nil, err
	}
	return &node, nil
}

<<<<<<< HEAD
func (r *GormNodeRepository) UpdateLatency(id uint, latency int) error {
=======
func (r *NodeRepository) UpdateLatency(id uint, latency int) error {
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	return r.db.Model(&model.Node{}).
		Where("id = ?", id).
		Update("latency", latency).Error
}

<<<<<<< HEAD
func (r *GormNodeRepository) FindByNodeID(nodeID string) (*model.Node, error) {
	var node model.Node
	err := r.db.Where("node_id = ?", nodeID).First(&node).Error
	if err != nil {
		return nil, err
	}
	return &node, nil
}

func (r *GormNodeRepository) UpdatePing(nodeID string, latency int) error {
	return r.db.Model(&model.Node{}).
		Where("node_id = ?", nodeID).
		Update("latency", latency).Error
}

func (r *GormNodeRepository) UpdateLoad(id uint, load int) error {
=======
func (r *NodeRepository) UpdateLoad(id uint, load int) error {
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	return r.db.Model(&model.Node{}).
		Where("id = ?", id).
		Update("load", load).Error
}

<<<<<<< HEAD
func (r *GormNodeRepository) UpdateLoadByNodeID(nodeID string, load int) error {
	return r.db.Model(&model.Node{}).
		Where("node_id = ?", nodeID).
		Update("load", load).Error
}

type GormRewardRepository struct {
=======
type RewardRepository struct {
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	db *gorm.DB
	rdb *redis.Client
}

<<<<<<< HEAD
func NewRewardRepository(db *gorm.DB) *GormRewardRepository {
	return &GormRewardRepository{db: db, rdb: nil}
}

func (r *GormRewardRepository) SetRedis(rdb *redis.Client) {
	r.rdb = rdb
}

func (r *GormRewardRepository) CreateAdLog(log *model.AdLog) error {
	return r.db.Create(log).Error
}

func (r *GormRewardRepository) FindAdLogByTransactionID(txnID string) (*model.AdLog, error) {
=======
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
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	var log model.AdLog
	err := r.db.Where("transaction_id = ?", txnID).First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

<<<<<<< HEAD
func (r *GormRewardRepository) FindRecentAdLog(userID uint, adType string) (*model.AdLog, error) {
	var log model.AdLog
	err := r.db.Where("user_id = ? AND ad_type = ?", userID, adType).
		Order("created_at DESC").First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *GormRewardRepository) CreateDailyCheckin(checkin *model.DailyCheckin) error {
	return r.db.Create(checkin).Error
}

func (r *GormRewardRepository) FindDailyCheckin(userID uint, date time.Time) (*model.DailyCheckin, error) {
=======
func (r *RewardRepository) CreateDailyCheckin(checkin *model.DailyCheckin) error {
	return r.db.Create(checkin).Error
}

func (r *RewardRepository) FindDailyCheckin(userID uint, date time.Time) (*model.DailyCheckin, error) {
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	var checkin model.DailyCheckin
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	err := r.db.Where("user_id = ? AND checkin_date = ?", userID, startOfDay).First(&checkin).Error
	if err != nil {
		return nil, err
	}
	return &checkin, nil
}

<<<<<<< HEAD
func (r *GormRewardRepository) CreateInvite(invite *model.Invite) error {
	return r.db.Create(invite).Error
}

func (r *GormRewardRepository) IsNonceUsed(nonce string) (bool, error) {
=======
func (r *RewardRepository) CreateInvite(invite *model.Invite) error {
	return r.db.Create(invite).Error
}

func (r *RewardRepository) IsNonceUsed(nonce string) (bool, error) {
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	if r.rdb == nil {
		return false, nil
	}
	ctx := context.Background()
<<<<<<< HEAD
	return r.rdb.Exists(ctx, "nonce:"+nonce).Val() > 0, nil
}

func (r *GormRewardRepository) SetNonce(nonce string, expire time.Duration) error {
=======
	return r.db.Exists(ctx, "nonce:"+nonce).Result()
}

func (r *RewardRepository) SetNonce(nonce string, expire time.Duration) error {
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	if r.rdb == nil {
		return nil
	}
	ctx := context.Background()
<<<<<<< HEAD
	return r.rdb.SetEX(ctx, "nonce:"+nonce, "1", expire).Err()
}

type GormSubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *GormSubscriptionRepository {
	return &GormSubscriptionRepository{db: db}
}

func (r *GormSubscriptionRepository) Create(sub *model.Subscription) error {
	return r.db.Create(sub).Error
}

func (r *GormSubscriptionRepository) FindByUserID(userID uint) ([]model.Subscription, error) {
=======
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
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	var subs []model.Subscription
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&subs).Error
	return subs, err
}

<<<<<<< HEAD
func (r *GormSubscriptionRepository) FindActiveByUserID(userID uint) (*model.Subscription, error) {
=======
func (r *SubscriptionRepository) FindActiveByUserID(userID uint) (*model.Subscription, error) {
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	var sub model.Subscription
	err := r.db.Where("user_id = ? AND status = ? AND expire_time > ?", userID, "active", time.Now()).First(&sub).Error
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

<<<<<<< HEAD
func (r *GormSubscriptionRepository) FindAll() ([]model.Plan, error) {
	var plans []model.Plan
	err := r.db.Find(&plans).Error
	return plans, err
}

func (r *GormSubscriptionRepository) FindAllPlans() ([]model.Plan, error) {
	var plans []model.Plan
	err := r.db.Find(&plans).Error
	return plans, err
}

type GormConnectRepository struct {
	db *gorm.DB
}

func NewConnectRepository(db *gorm.DB) *GormConnectRepository {
	return &GormConnectRepository{db: db}
}

func (r *GormConnectRepository) CreateSession(session *model.ConnectSession) error {
	return r.db.Create(session).Error
}

func (r *GormConnectRepository) FindActiveSession(userID uint) (*model.ConnectSession, error) {
	var session model.ConnectSession
	err := r.db.Where("user_id = ? AND status = ?", userID, "active").First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *GormConnectRepository) UpdateSession(session *model.ConnectSession) error {
	return r.db.Save(session).Error
}

type GormVerificationRepository struct {
	db *gorm.DB
}

func NewVerificationRepository(db *gorm.DB) *GormVerificationRepository {
	return &GormVerificationRepository{db: db}
}

func (r *GormVerificationRepository) Create(code *model.VerificationCode) error {
	return r.db.Create(code).Error
}

func (r *GormVerificationRepository) FindValidCode(phone, code string) (*model.VerificationCode, error) {
	var vcode model.VerificationCode
	err := r.db.Where("phone = ? AND code = ? AND used = ? AND expires_at > ?", phone, code, false, time.Now()).
		Order("created_at DESC").First(&vcode).Error
	if err != nil {
		return nil, err
	}
	return &vcode, nil
}

func (r *GormVerificationRepository) MarkUsed(id uint) error {
	return r.db.Model(&model.VerificationCode{}).Where("id = ?", id).Update("used", true).Error
=======
func (r *SubscriptionRepository) FindAll() ([]model.Plan, error) {
	var plans []model.Plan
	err := r.db.Find(&plans).Error
	return plans, err
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
}