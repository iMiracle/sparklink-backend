package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"sparklink-backend/model"
)

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
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) FindByPhone(phone string) (*model.User, error) {
	var user model.User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) FindByDeviceID(deviceID string) (*model.User, error) {
	var user model.User
	err := r.db.Where("device_id = ?", deviceID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) FindByInviteCode(code string) (*model.User, error) {
	var user model.User
	err := r.db.Where("invite_code = ?", code).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) Save(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *GormUserRepository) AddBalance(userID uint, minutes int) error {
	return r.db.Model(&model.User{}).
		Where("id = ?", userID).
		Update("balance_minutes", gorm.Expr("balance_minutes + ?", minutes)).Error
}

func (r *GormUserRepository) CountDevicesByUserID(userID uint) (int, error) {
	var count int64
	err := r.db.Model(&model.Device{}).Where("user_id = ? AND is_active = ?", userID, true).Count(&count).Error
	return int(count), err
}

func (r *GormUserRepository) FindDevicesByUserID(userID uint) ([]model.Device, error) {
	var devices []model.Device
	err := r.db.Where("user_id = ? AND is_active = ?", userID, true).Find(&devices).Error
	return devices, err
}

func (r *GormUserRepository) CreateDevice(device *model.Device) error {
	return r.db.Create(device).Error
}

func (r *GormUserRepository) DeactivateDevice(deviceID string) error {
	return r.db.Model(&model.Device{}).
		Where("device_id = ?", deviceID).
		Update("is_active", false).Error
}

type GormNodeRepository struct {
	db *gorm.DB
}

func NewNodeRepository(db *gorm.DB) *GormNodeRepository {
	return &GormNodeRepository{db: db}
}

func (r *GormNodeRepository) FindAll(protocol, visibility, region string) ([]model.Node, error) {
	var nodes []model.Node
	query := r.db.Where("status = ?", "online")

	if protocol != "" {
		query = query.Where("protocol = ?", protocol)
	}
	if visibility != "" {
		query = query.Where("visibility_level = ?", visibility)
	}
	if region != "" {
		query = query.Where("region_code = ?", region)
	}

	err := query.Order("latency ASC").Find(&nodes).Error
	return nodes, err
}

func (r *GormNodeRepository) FindByNodeID(nodeID string) (*model.Node, error) {
	var node model.Node
	err := r.db.Where("node_id = ?", nodeID).First(&node).Error
	if err != nil {
		return nil, err
	}
	return &node, nil
}

func (r *GormNodeRepository) FindByID(id uint) (*model.Node, error) {
	var node model.Node
	err := r.db.First(&node, id).Error
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

func (r *GormNodeRepository) UpdateLatency(id uint, latency int) error {
	return r.db.Model(&model.Node{}).
		Where("id = ?", id).
		Update("latency", latency).Error
}

func (r *GormNodeRepository) GetRegions() ([]string, error) {
	var regions []string
	err := r.db.Model(&model.Node{}).
		Select("DISTINCT region_code").
		Where("status = ?", "online").
		Order("region_code ASC").
		Pluck("region_code", &regions).Error
	return regions, err
}

func (r *GormNodeRepository) UpdateLoad(id uint, load int) error {
	return r.db.Model(&model.Node{}).
		Where("id = ?", id).
		Update("load", load).Error
}

func (r *GormNodeRepository) Create(node *model.Node) error {
	return r.db.Create(node).Error
}

func (r *GormNodeRepository) UpdateByNodeID(nodeID string, updates map[string]interface{}) error {
	return r.db.Model(&model.Node{}).Where("node_id = ?", nodeID).Updates(updates).Error
}

func (r *GormNodeRepository) AddFavorite(userID uint, nodeID string) error {
	fav := &model.Favorite{
		UserID: userID,
		NodeID: nodeID,
	}
	return r.db.Create(fav).Error
}

func (r *GormNodeRepository) RemoveFavorite(userID uint, nodeID string) error {
	return r.db.Where("user_id = ? AND node_id = ?", userID, nodeID).
		Delete(&model.Favorite{}).Error
}

func (r *GormNodeRepository) GetFavorites(userID uint) ([]model.Favorite, error) {
	var favs []model.Favorite
	err := r.db.Where("user_id = ?", userID).Find(&favs).Error
	return favs, err
}

type GormRewardRepository struct {
	db  *gorm.DB
	rdb *redis.Client
}

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
	var log model.AdLog
	err := r.db.Where("transaction_id = ?", txnID).First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

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
	var checkin model.DailyCheckin
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	err := r.db.Where("user_id = ? AND checkin_date = ?", userID, startOfDay).First(&checkin).Error
	if err != nil {
		return nil, err
	}
	return &checkin, nil
}

func (r *GormRewardRepository) CreateInvite(invite *model.Invite) error {
	return r.db.Create(invite).Error
}

func (r *GormRewardRepository) IsNonceUsed(nonce string) (bool, error) {
	if r.rdb == nil {
		return false, nil
	}
	ctx := context.Background()
	return r.rdb.Exists(ctx, "nonce:"+nonce).Val() > 0, nil
}

func (r *GormRewardRepository) SetNonce(nonce string, expire time.Duration) error {
	if r.rdb == nil {
		return nil
	}
	ctx := context.Background()
	return r.rdb.SetEX(ctx, "nonce:"+nonce, "1", expire).Err()
}

type GormSubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *GormSubscriptionRepository {
	return &GormSubscriptionRepository{db: db}
}

func (r *GormSubscriptionRepository) FindPlanByID(planID string) (*model.Plan, error) {
	var plan model.Plan
	err := r.db.Where("plan_id = ?", planID).First(&plan).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *GormSubscriptionRepository) FindAllPlans() ([]model.Plan, error) {
	var plans []model.Plan
	err := r.db.Find(&plans).Error
	return plans, err
}

func (r *GormSubscriptionRepository) Create(sub *model.Subscription) error {
	return r.db.Create(sub).Error
}

func (r *GormSubscriptionRepository) FindActiveByUserID(userID uint) (*model.Subscription, error) {
	var sub model.Subscription
	err := r.db.Where("user_id = ? AND status = ? AND expire_time > ?", userID, "active", time.Now()).First(&sub).Error
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *GormSubscriptionRepository) CreatePlan(plan *model.Plan) error {
	return r.db.Create(plan).Error
}

func (r *GormSubscriptionRepository) UpdatePlan(planID string, updates map[string]interface{}) error {
	return r.db.Model(&model.Plan{}).Where("plan_id = ?", planID).Updates(updates).Error
}

func (r *GormSubscriptionRepository) DeletePlan(planID string) error {
	return r.db.Where("plan_id = ?", planID).Delete(&model.Plan{}).Error
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
}

func (r *GormVerificationRepository) CreateQRSession(session *model.QRSession) error {
	return r.db.Create(session).Error
}

func (r *GormVerificationRepository) FindQRSessionBySessionID(sessionID string) (*model.QRSession, error) {
	var session model.QRSession
	err := r.db.Where("session_id = ?", sessionID).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *GormVerificationRepository) UpdateQRSessionStatus(sessionID string, status string, userID *uint, token string) error {
	return r.db.Model(&model.QRSession{}).
		Where("session_id = ?", sessionID).
		Updates(map[string]interface{}{
			"status":  status,
			"user_id": userID,
			"token":   token,
		}).Error
}

type GormAdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *GormAdminRepository {
	return &GormAdminRepository{db: db}
}

func (r *GormAdminRepository) Create(admin *model.AdminUser) error {
	return r.db.Create(admin).Error
}

func (r *GormAdminRepository) FindByID(id uint) (*model.AdminUser, error) {
	var admin model.AdminUser
	err := r.db.First(&admin, id).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *GormAdminRepository) FindByUsername(username string) (*model.AdminUser, error) {
	var admin model.AdminUser
	err := r.db.Where("username = ?", username).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *GormAdminRepository) FindAll() ([]model.AdminUser, error) {
	var admins []model.AdminUser
	err := r.db.Find(&admins).Error
	return admins, err
}

func (r *GormAdminRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&model.AdminUser{}).Where("id = ?", id).Update("status", status).Error
}

func (r *GormAdminRepository) Delete(id uint) error {
	return r.db.Delete(&model.AdminUser{}, id).Error
}

type GormAnnouncementRepository struct {
	db *gorm.DB
}

func NewAnnouncementRepository(db *gorm.DB) *GormAnnouncementRepository {
	return &GormAnnouncementRepository{db: db}
}

func (r *GormAnnouncementRepository) Create(a *model.Announcement) error {
	return r.db.Create(a).Error
}

func (r *GormAnnouncementRepository) FindByID(id uint) (*model.Announcement, error) {
	var a model.Announcement
	err := r.db.First(&a, id).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *GormAnnouncementRepository) FindAll() ([]model.Announcement, error) {
	var list []model.Announcement
	err := r.db.Order("created_at DESC").Find(&list).Error
	return list, err
}

func (r *GormAnnouncementRepository) Update(a *model.Announcement) error {
	return r.db.Save(a).Error
}

func (r *GormAnnouncementRepository) Delete(id uint) error {
	return r.db.Delete(&model.Announcement{}, id).Error
}

type GormAuditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) *GormAuditLogRepository {
	return &GormAuditLogRepository{db: db}
}

func (r *GormAuditLogRepository) Create(log *model.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *GormAuditLogRepository) FindAll(page, pageSize int) ([]model.AuditLog, int64, error) {
	var logs []model.AuditLog
	var total int64
	r.db.Model(&model.AuditLog{}).Count(&total)
	err := r.db.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs).Error
	return logs, total, err
}

func GenerateSessionID(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

func GenerateInviteCode() string {
	return fmt.Sprintf("SPARK-%s", strings.ToUpper(time.Now().Format("150405")))
}
