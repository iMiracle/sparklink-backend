package service

import (
	"sparklink-backend/config"
	"sparklink-backend/repository"
)

type Services struct {
	Auth         *AuthService
	Node         *NodeService
	Reward       *RewardService
	Subscription *SubscriptionService
	Connect      *ConnectService
	User         *UserService
}

type Repos struct {
	User         repository.UserRepository
	Node         repository.NodeRepository
	Reward       repository.RewardRepository
	Subscription repository.SubscriptionRepository
	Connect      repository.ConnectRepository
	Verification repository.VerificationRepository
}

func NewServices(repos *Repos, cfg *config.Config) *Services {
	return &Services{
		Auth:         NewAuthService(repos.User, repos.Verification, cfg),
		Node:         NewNodeService(repos.Node),
		Reward:       NewRewardService(repos.Reward, repos.User),
		Subscription: NewSubscriptionService(repos.Subscription),
		Connect:      NewConnectService(repos.Connect, repos.Node),
		User:         NewUserService(repos.User),
	}
}
