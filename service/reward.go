package service

import (
	"errors"
	"time"
	"sparklink-backend/model"
	"sparklink-backend/repository"
)

type RewardService struct {
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
