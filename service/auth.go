package service

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"sparklink-backend/config"
	"sparklink-backend/model"
	"sparklink-backend/pkg/auth"
	"sparklink-backend/repository"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	userRepo  repository.UserRepository
	verifRepo repository.VerificationRepository
	cfg       *config.Config
}

func NewAuthService(userRepo repository.UserRepository, verifRepo repository.VerificationRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		verifRepo: verifRepo,
		cfg:       cfg,
	}
}

func (s *AuthService) SendCode(phone string) error {
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	vcode := &model.VerificationCode{
		Phone:     phone,
		Code:      code,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	return s.verifRepo.Create(vcode)
}

func (s *AuthService) Register(phone, code, inviteCode string) (*model.User, string, error) {
	existing, _ := s.userRepo.FindByPhone(phone)
	if existing != nil {
		return nil, "", errors.New("手机号已注册")
	}
	vcode, err := s.verifRepo.FindValidCode(phone, code)
	if err != nil {
		return nil, "", errors.New("验证码错误或已过期")
	}
	s.verifRepo.MarkUsed(vcode.ID)

	user := &model.User{
		Phone:          phone,
		Nickname:       "User",
		VipStatus:      "inactive",
		BalanceMinutes: 60,
		InviteCode:     repository.GenerateInviteCode(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if inviteCode != "" {
		referrer, err := s.userRepo.FindByInviteCode(inviteCode)
		if err == nil {
			referredBy := referrer.ID
			user.ReferredBy = &referredBy
			referrer.InvitedCount++
			s.userRepo.Save(referrer)
		}
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, "", err
	}

	token, err := s.GenerateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) Login(phone, code string) (*model.User, string, error) {
	user, err := s.userRepo.FindByPhone(phone)
	if err != nil {
		return nil, "", errors.New("手机号未注册")
	}

	vcode, err := s.verifRepo.FindValidCode(phone, code)
	if err != nil {
		return nil, "", errors.New("验证码错误或已过期")
	}
	s.verifRepo.MarkUsed(vcode.ID)

	token, err := s.GenerateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) GenerateToken(userID uint) (string, error) {
	claims := auth.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.cfg.JWTExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

func (s *AuthService) ValidateToken(tokenString string) (*auth.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*auth.Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func (s *AuthService) GetTokenExpiry() time.Time {
	return time.Now().Add(s.cfg.JWTExpire)
}

func (s *AuthService) GetQRStatus(sessionID string) (string, interface{}, interface{}) {
	return "pending", nil, nil
}
