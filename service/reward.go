package service

import (
<<<<<<< HEAD
	"errors"
	"time"
=======
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	"sparklink-backend/model"
	"sparklink-backend/repository"
)

type RewardService struct {
<<<<<<< HEAD
	rewardRepo repository.RewardRepository
	userRepo   repository.UserRepository
}

func NewRewardService(rewardRepo repository.RewardRepository, userRepo repository.UserRepository) *RewardService {
	return &RewardService{
		rewardRepo: rewardRepo,
		userRepo:   userRepo,
	}
}

const (
	RewardVideoMinutes   = 30
	RewardCheckinMinutes = 30
	RewardCooldownMinutes = 30
	DailyRewardLimit     = 3
)

func (s *RewardService) ClaimReward(userID uint, adID, adType string) (reward int, balance int, cooldownEndsAt string, err error) {
	recent, _ := s.rewardRepo.FindRecentAdLog(userID, adType)
	if recent != nil {
		cooldownEnd := recent.CooldownEnd
		if time.Now().Before(cooldownEnd) {
			return 0, 0, cooldownEnd.Format(time.RFC3339), errors.New("冷却中")
		}
	}

	if adType == "sign" && recent != nil {
		lastSignDay := recent.CreatedAt.Truncate(24 * time.Hour)
		if time.Now().Truncate(24 * time.Hour).Equal(lastSignDay) {
			return 0, 0, "", errors.New("已达每日签到上限")
		}
	}

	switch adType {
	case "reward_video":
		reward = RewardVideoMinutes
	case "sign":
		reward = RewardCheckinMinutes
	default:
		return 0, 0, "", errors.New("不支持的激励类型")
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return 0, 0, "", err
	}

	cooldownEnd := time.Now().Add(RewardCooldownMinutes * time.Minute)
	log := &model.AdLog{
		UserID:      userID,
		AdID:        adID,
		AdType:      adType,
		Reward:      reward,
		CooldownEnd: cooldownEnd,
	}
	s.rewardRepo.CreateAdLog(log)
	s.userRepo.AddBalance(userID, reward)

	user, _ = s.userRepo.FindByID(userID)
	return reward, user.BalanceMinutes, cooldownEnd.Format(time.RFC3339), nil
}

func (s *RewardService) GetCooldown(userID uint, adType string) (inCooldown bool, remainingSeconds int, cooldownEndsAt string, err error) {
	recent, _ := s.rewardRepo.FindRecentAdLog(userID, adType)
	if recent == nil {
		return false, 0, "", nil
	}
	remaining := time.Until(recent.CooldownEnd)
	if remaining <= 0 {
		return false, 0, "", nil
	}
	return true, int(remaining.Seconds()), recent.CooldownEnd.Format(time.RFC3339), nil
}

func (s *RewardService) GetInviteInfo(userID uint) (inviteCode string, invitedCount int, totalRewardHours int, shareChannels []string, err error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return "", 0, 0, nil, err
	}
	return user.InviteCode, user.InvitedCount, user.InvitedCount * 24, []string{"wechat", "qq", "moments", "link", "qrcode"}, nil
}
=======
	rewardRepo *repository.RewardRepository
}

func NewRewardService(rewardRepo *repository.RewardRepository) *RewardService {
	return &RewardService{rewardRepo: rewardRepo}
}

func (s *RewardService) ClaimVideoReward(userID uint, adPlatform, adID, transactionID, nonce string) error {
	// 检查是否已经处理过
	existing, _ := s.rewardRepo.FindAdLogByTransactionID(transactionID)
	if existing != nil {
		return errors.New("transaction already processed")
	}

	// 检查 nonce
	used, _ := s.rewardRepo.IsNonceUsed(nonce)
	if used {
		return errors.New("nonce already used")
	}

	// 创建广告日志
	adLog := &model.AdLog{
		UserID:         userID,
		AdPlatform:    adPlatform,
		AdID:          adID,
		TransactionID: transactionID,
		RewardAmount:  120,
		Status:        "success",
		CreatedAt:     time.Now(),
	}

	if err := s.rewardRepo.CreateAdLog(adLog); err != nil {
		return err
	}

	// 记录nonce防止重放
	s.rewardRepo.SetNonce(nonce, 24*time.Hour)

	return nil
}

func (s *RewardService) DailyCheckin(userID uint) (int, error) {
	// 检查今天是否已签到
	checkin, _ := s.rewardRepo.FindDailyCheckin(userID, time.Now())
	if checkin != nil {
		return 0, errors.New("already checked in today")
	}

	// 创建签到记录
	newCheckin := &model.DailyCheckin{
		UserID:        userID,
		CheckinDate:   time.Now(),
		RewardAmount:  30,
		CreatedAt:    time.Now(),
	}

	if err := s.rewardRepo.CreateDailyCheckin(newCheckin); err != nil {
		return 0, err
	}

	return newCheckin.RewardAmount, nil
}

func (s *RewardService) GetInviteInfo(userID uint) (string, int, error) {
	// 这里应该从 userRepo 获取
	return "", 1440, nil // 24小时 = 1440分钟
}

func (s *RewardService) BindInvite(userID uint, referralCode string) error {
	// 检查邀请码是否存在
	// 如果存在，给邀请人和被邀请人发放奖励
	return nil
}

func (s *RewardService) GetRemainingTime(userID uint) (int, error) {
	// 计算剩余时间（分钟）
	return 0, nil
}

func GenerateNonce() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
