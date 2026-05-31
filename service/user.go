package service

import (
	"errors"
	"time"

	"sparklink-backend/model"
	"sparklink-backend/repository"
)

const (
	MaxDevicesFree = 1
	MaxDevicesVIP  = 10
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

func (s *UserService) MaxDevicesForUser(userID uint) (int, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return 0, err
	}
	if user.VipStatus == "active" {
		return MaxDevicesVIP, nil
	}
	return MaxDevicesFree, nil
}

func (s *UserService) RegisterDevice(userID uint, deviceID, deviceName, deviceType string) error {
	count, err := s.userRepo.CountDevicesByUserID(userID)
	if err != nil {
		return err
	}
	maxDevices, err := s.MaxDevicesForUser(userID)
	if err != nil {
		return err
	}
	if count >= maxDevices {
		return errors.New("device limit reached")
	}
	return s.userRepo.CreateDevice(&model.Device{
		UserID:     userID,
		DeviceID:   deviceID,
		DeviceName: deviceName,
		DeviceType: deviceType,
		LastActive: time.Now(),
	})
}

func (s *UserService) ToggleKillSwitch(userID uint) (bool, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return false, err
	}
	user.KillSwitchEnabled = !user.KillSwitchEnabled
	if err := s.userRepo.Save(user); err != nil {
		return false, err
	}
	return user.KillSwitchEnabled, nil
}

func (s *UserService) KillSwitchStatus(userID uint) (bool, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return false, err
	}
	return user.KillSwitchEnabled, nil
}
