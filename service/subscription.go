package service

import (
	"time"

	"sparklink-backend/model"
	"sparklink-backend/repository"
)

type SubscriptionService struct {
	subRepo repository.SubscriptionRepository
}

func NewSubscriptionService(subRepo repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{subRepo: subRepo}
}

func (s *SubscriptionService) GetPlans() ([]model.Plan, error) {
	return s.subRepo.FindAllPlans()
}

func (s *SubscriptionService) CreateSubscription(userID uint, planID string, amount float64) (*model.Subscription, error) {
	plan, err := s.subRepo.FindPlanByID(planID)
	if err != nil {
		return nil, err
	}
	sub := &model.Subscription{
		UserID:     userID,
		PlanID:     planID,
		Amount:     amount,
		StartTime:  time.Now(),
		ExpireTime: time.Now().AddDate(0, 0, plan.DurationDays),
		Status:     "active",
		CreatedAt:  time.Now(),
	}
	if err := s.subRepo.Create(sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *SubscriptionService) VerifySubscription(userID uint) (bool, error) {
	_, err := s.subRepo.FindActiveByUserID(userID)
	return err == nil, nil
}

func (s *SubscriptionService) FindActive(userID uint) (*model.Subscription, error) {
	return s.subRepo.FindActiveByUserID(userID)
}
