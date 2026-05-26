package repository

import (
	"errors"
	"fmt"
	"sparklink-backend/mock"
	"sparklink-backend/model"
	"strings"
	"sync"
	"time"
)

type MockUserRepository struct {
	mu   sync.RWMutex
	data *mock.MockData
}

func NewMockUserRepository(data *mock.MockData) *MockUserRepository {
	return &MockUserRepository{data: data}
}

func (r *MockUserRepository) Create(user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	user.ID = uint(len(r.data.Users) + 1)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	r.data.Users = append(r.data.Users, *user)
	return nil
}

func (r *MockUserRepository) FindByID(id uint) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for i := range r.data.Users {
		if r.data.Users[i].ID == id {
			return &r.data.Users[i], nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *MockUserRepository) FindByPhone(phone string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for i := range r.data.Users {
		if r.data.Users[i].Phone == phone {
			return &r.data.Users[i], nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *MockUserRepository) FindByInviteCode(code string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for i := range r.data.Users {
		if r.data.Users[i].InviteCode == code {
			return &r.data.Users[i], nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *MockUserRepository) Save(user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.data.Users {
		if r.data.Users[i].ID == user.ID {
			user.UpdatedAt = time.Now()
			r.data.Users[i] = *user
			return nil
		}
	}
	return errors.New("user not found")
}

func (r *MockUserRepository) AddBalance(userID uint, minutes int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.data.Users {
		if r.data.Users[i].ID == userID {
			r.data.Users[i].BalanceMinutes += minutes
			r.data.Users[i].UpdatedAt = time.Now()
			return nil
		}
	}
	return errors.New("user not found")
}

func (r *MockUserRepository) FindDevicesByUserID(userID uint) ([]model.Device, error) {
	return []model.Device{}, nil
}

func (r *MockUserRepository) CreateDevice(device *model.Device) error {
	return nil
}

func (r *MockUserRepository) DeactivateDevice(deviceID string) error {
	return nil
}

type MockNodeRepository struct {
	mu   sync.RWMutex
	data *mock.MockData
}

func NewMockNodeRepository(data *mock.MockData) *MockNodeRepository {
	return &MockNodeRepository{data: data}
}

func (r *MockNodeRepository) FindAll(protocol, visibility, region string) ([]model.Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []model.Node
	for _, n := range r.data.Nodes {
		if protocol != "" && n.Protocol != protocol {
			continue
		}
		if visibility != "" && n.VisibilityLevel != visibility {
			continue
		}
		if region != "" && n.RegionCode != region {
			continue
		}
		result = append(result, n)
	}
	return result, nil
}

func (r *MockNodeRepository) FindByNodeID(nodeID string) (*model.Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for i := range r.data.Nodes {
		if r.data.Nodes[i].NodeId == nodeID {
			return &r.data.Nodes[i], nil
		}
	}
	return nil, errors.New("node not found")
}

func (r *MockNodeRepository) UpdatePing(nodeID string, latency int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.data.Nodes {
		if r.data.Nodes[i].NodeId == nodeID {
			r.data.Nodes[i].Latency = latency
			r.data.Nodes[i].UpdatedAt = time.Now()
			return nil
		}
	}
	return errors.New("node not found")
}

func (r *MockNodeRepository) UpdateLoad(nodeID string, load int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.data.Nodes {
		if r.data.Nodes[i].NodeId == nodeID {
			r.data.Nodes[i].Load = load
			r.data.Nodes[i].UpdatedAt = time.Now()
			return nil
		}
	}
	return errors.New("node not found")
}

type MockRewardRepository struct {
	mu   sync.RWMutex
	data *mock.MockData
}

func NewMockRewardRepository(data *mock.MockData) *MockRewardRepository {
	return &MockRewardRepository{data: data}
}

func (r *MockRewardRepository) CreateAdLog(log *model.AdLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	log.ID = uint(len(r.data.AdLogs) + 1)
	log.CreatedAt = time.Now()
	r.data.AdLogs = append(r.data.AdLogs, *log)
	return nil
}

func (r *MockRewardRepository) FindRecentAdLog(userID uint, adType string) (*model.AdLog, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for i := len(r.data.AdLogs) - 1; i >= 0; i-- {
		if r.data.AdLogs[i].UserID == userID && r.data.AdLogs[i].AdType == adType {
			return &r.data.AdLogs[i], nil
		}
	}
	return nil, errors.New("no ad log found")
}

type MockSubscriptionRepository struct {
	mu   sync.RWMutex
	data *mock.MockData
}

func NewMockSubscriptionRepository(data *mock.MockData) *MockSubscriptionRepository {
	return &MockSubscriptionRepository{data: data}
}

func (r *MockSubscriptionRepository) FindAllPlans() ([]model.Plan, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.data.Plans, nil
}

func (r *MockSubscriptionRepository) Create(sub *model.Subscription) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	sub.ID = uint(len(r.data.Subscriptions) + 1)
	sub.CreatedAt = time.Now()
	r.data.Subscriptions = append(r.data.Subscriptions, *sub)
	return nil
}

func (r *MockSubscriptionRepository) FindActiveByUserID(userID uint) (*model.Subscription, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for i := range r.data.Subscriptions {
		s := r.data.Subscriptions[i]
		if s.UserID == userID && s.Status == "active" && s.ExpireTime.After(time.Now()) {
			return &s, nil
		}
	}
	return nil, errors.New("no active subscription")
}

type MockConnectRepository struct {
	mu   sync.RWMutex
	data *mock.MockData
}

func NewMockConnectRepository(data *mock.MockData) *MockConnectRepository {
	return &MockConnectRepository{data: data}
}

func (r *MockConnectRepository) CreateSession(session *model.ConnectSession) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	session.ID = uint(len(r.data.ConnectSessions) + 1)
	session.CreatedAt = time.Now()
	r.data.ConnectSessions = append(r.data.ConnectSessions, *session)
	return nil
}

func (r *MockConnectRepository) FindActiveSession(userID uint) (*model.ConnectSession, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for i := range r.data.ConnectSessions {
		if r.data.ConnectSessions[i].UserID == userID && r.data.ConnectSessions[i].Status == "active" {
			return &r.data.ConnectSessions[i], nil
		}
	}
	return nil, errors.New("no active session")
}

func (r *MockConnectRepository) UpdateSession(session *model.ConnectSession) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.data.ConnectSessions {
		if r.data.ConnectSessions[i].SessionID == session.SessionID {
			r.data.ConnectSessions[i] = *session
			return nil
		}
	}
	return errors.New("session not found")
}

type MockVerificationRepository struct {
	mu   sync.RWMutex
	data *mock.MockData
}

func NewMockVerificationRepository(data *mock.MockData) *MockVerificationRepository {
	return &MockVerificationRepository{data: data}
}

func (r *MockVerificationRepository) Create(code *model.VerificationCode) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	code.ID = uint(len(r.data.VerificationCodes) + 1)
	code.CreatedAt = time.Now()
	r.data.VerificationCodes = append(r.data.VerificationCodes, *code)
	return nil
}

func (r *MockVerificationRepository) FindValidCode(phone, code string) (*model.VerificationCode, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for i := len(r.data.VerificationCodes) - 1; i >= 0; i-- {
		v := r.data.VerificationCodes[i]
		if v.Phone == phone && v.Code == code && !v.Used && v.ExpiresAt.After(time.Now()) {
			return &v, nil
		}
	}
	return nil, errors.New("invalid or expired code")
}

func (r *MockVerificationRepository) MarkUsed(id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.data.VerificationCodes {
		if r.data.VerificationCodes[i].ID == id {
			r.data.VerificationCodes[i].Used = true
			return nil
		}
	}
	return errors.New("code not found")
}

type MockRepository struct {
	User         *MockUserRepository
	Node         *MockNodeRepository
	Reward       *MockRewardRepository
	Subscription *MockSubscriptionRepository
	Connect      *MockConnectRepository
	Verification *MockVerificationRepository
}

func NewMockRepository(data *mock.MockData) *MockRepository {
	return &MockRepository{
		User:         NewMockUserRepository(data),
		Node:         NewMockNodeRepository(data),
		Reward:       NewMockRewardRepository(data),
		Subscription: NewMockSubscriptionRepository(data),
		Connect:      NewMockConnectRepository(data),
		Verification: NewMockVerificationRepository(data),
	}
}

func GenerateSessionID(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

func GenerateInviteCode() string {
	return fmt.Sprintf("SPARK-%s", strings.ToUpper(time.Now().Format("150405")))
}
