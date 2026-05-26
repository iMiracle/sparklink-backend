package service

import (
	"errors"
<<<<<<< HEAD
	"fmt"
	"math/rand"
=======
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	"time"

	"sparklink-backend/config"
	"sparklink-backend/model"
<<<<<<< HEAD
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
=======
	"sparklink-backend/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repository.UserRepository
	cfg     *config.Config
}

func NewAuthService(userRepo *repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{userRepo: userRepo, cfg: cfg}
}

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *AuthService) Register(email, phone, password, deviceID string) (*model.User, string, error) {
	// 检查用户是否存在
	if email != "" {
		if _, err := s.userRepo.FindByEmail(email); err == nil {
			return nil, "", errors.New("email already exists")
		}
	}

	if phone != "" {
		if _, err := s.userRepo.FindByPhone(phone); err == nil {
			return nil, "", errors.New("phone already exists")
		}
	}

	// 检查设备是否已注册
	if _, err := s.userRepo.FindByDeviceID(deviceID); err == nil {
		user, _ := s.userRepo.FindByDeviceID(deviceID)
		token, _ := s.GenerateToken(user.ID)
		return user, token, nil
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	// 生成邀请码
	referralCode := generateReferralCode()

	user := &model.User{
		Email:         email,
		Phone:         phone,
		Password:      string(hashedPassword),
		DeviceID:      deviceID,
		ReferralCode:  referralCode,
		AdCredits:    120, // 新用户赠送120分钟
		CreatedAt:     time.Now(),
		UpdatedAt:    time.Now(),
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, "", err
	}

	token, err := s.GenerateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

<<<<<<< HEAD
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
=======
	// 创建设备记录
	s.userRepo.CreateDevice(&model.Device{
		UserID:     user.ID,
		DeviceID:  deviceID,
		LastLogin:  time.Now(),
		IsActive:   true,
		CreatedAt: time.Now(),
	})

	return user, token, nil
}

func (s *AuthService) Login(email, phone, password, deviceID string) (*model.User, string, error) {
	var user *model.User
	var err error

	if email != "" {
		user, err = s.userRepo.FindByEmail(email)
	} else if phone != "" {
		user, err = s.userRepo.FindByPhone(phone)
	} else {
		return nil, "", errors.New("email or phone required")
	}

	if err != nil {
		return nil, "", err
	}

	if password != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			return nil, "", errors.New("invalid password")
		}
	}
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f

	token, err := s.GenerateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

<<<<<<< HEAD
=======
	// 更新设备记录
	s.userRepo.CreateDevice(&model.Device{
		UserID:     user.ID,
		DeviceID:  deviceID,
		LastLogin:  time.Now(),
		IsActive:   true,
		CreatedAt: time.Now(),
	})

>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	return user, token, nil
}

func (s *AuthService) GenerateToken(userID uint) (string, error) {
<<<<<<< HEAD
	claims := auth.Claims{
=======
	claims := &Claims{
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.cfg.JWTExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
<<<<<<< HEAD
=======

>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

<<<<<<< HEAD
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
=======
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func generateReferralCode() string {
	// 简单的邀请码生成
	return "SPARK" + time.Now().Format("060102150405")
}
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
