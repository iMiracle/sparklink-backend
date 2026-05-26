package service

import (
	"sparklink-backend/model"
	"sparklink-backend/repository"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetProfile(userID uint) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *UserService) GetDevices(userID uint) ([]model.Device, error) {
	return s.userRepo.FindDevicesByUserID(userID)
}

func (s *UserService) RemoveDevice(deviceID string) error {
	return s.userRepo.DeactivateDevice(deviceID)
}
