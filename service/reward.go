package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"sparklink-backend/model"
	"sparklink-backend/repository"
)

type RewardService struct {
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