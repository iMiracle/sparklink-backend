package service

import (
	"time"
<<<<<<< HEAD
=======

>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	"sparklink-backend/model"
	"sparklink-backend/repository"
)

type SubscriptionService struct {
<<<<<<< HEAD
	subRepo repository.SubscriptionRepository
}

func NewSubscriptionService(subRepo repository.SubscriptionRepository) *SubscriptionService {
=======
	subRepo *repository.SubscriptionRepository
}

func NewSubscriptionService(subRepo *repository.SubscriptionRepository) *SubscriptionService {
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	return &SubscriptionService{subRepo: subRepo}
}

func (s *SubscriptionService) GetPlans() ([]model.Plan, error) {
<<<<<<< HEAD
	return s.subRepo.FindAllPlans()
}

func (s *SubscriptionService) CreateSubscription(userID uint, planID string, amount float64) (*model.Subscription, error) {
	sub := &model.Subscription{
		UserID:    userID,
		PlanID:    planID,
		Amount:    amount,
		StartTime: time.Now(),
		ExpireTime: time.Now().AddDate(0, 1, 0),
		Status:    "active",
		CreatedAt: time.Now(),
	}
	err := s.subRepo.Create(sub)
	if err != nil {
		return nil, err
	}
=======
	defaultPlans := []model.Plan{
		{Plan: "weekly", Name: "Weekly", Price: 2.99, Duration: 7 * 24 * 60, Features: "VIP nodes"},
		{Plan: "monthly", Name: "Monthly", Price: 9.99, Duration: 30 * 24 * 60, IsPopular: true, Features: "VIP nodes + No ads"},
		{Plan: "quarterly", Name: "Quarterly", Price: 24.99, Duration: 90 * 24 * 60, Features: "VIP nodes + No ads"},
		{Plan: "yearly", Name: "Yearly", Price: 79.99, Duration: 365 * 24 * 60, Features: "All features + Priority support"},
	}

	return defaultPlans, nil
}

func (s *SubscriptionService) CreateSubscription(userID uint, plan string, amount float64) (*model.Subscription, error) {
	plans, _ := s.GetPlans()
	var selectedPlan model.Plan
	for _, p := range plans {
		if p.Plan == plan {
			selectedPlan = p
			break
		}
	}

	if selectedPlan.Plan == "" {
		return nil, nil
	}

	sub := &model.Subscription{
		UserID:     userID,
		Plan:       plan,
		Amount:     amount,
		StartTime:  time.Now(),
		ExpireTime: time.Now().Add(time.Duration(selectedPlan.Duration) * time.Minute),
		Status:     "active",
		CreatedAt:  time.Now(),
	}

	if err := s.subRepo.Create(sub); err != nil {
		return nil, err
	}

>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	return sub, nil
}

func (s *SubscriptionService) VerifySubscription(userID uint) (bool, error) {
<<<<<<< HEAD
	_, err := s.subRepo.FindActiveByUserID(userID)
	return err == nil, nil
}
=======
	sub, err := s.subRepo.FindActiveByUserID(userID)
	if err != nil {
		return false, nil
	}
	return sub != nil, nil
}
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
