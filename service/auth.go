package service

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
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

func (s *AuthService) sendSMS(phone, code string) {
	if s.cfg.SMSAPIKey == "" || s.cfg.SMSEndpoint == "" {
		log.Printf("[SMS] 开发模式: 向 %s 发送验证码 %s", phone, code)
		return
	}
	go func() {
		body := fmt.Sprintf(`{"phone":"%s","code":"%s","apiKey":"%s","appSecret":"%s"}`,
			phone, code, s.cfg.SMSAPIKey, s.cfg.SMSAppSecret)
		_, err := http.Post(s.cfg.SMSEndpoint, "application/json", strings.NewReader(body))
		if err != nil {
			log.Printf("[SMS] 发送失败: %v", err)
		}
	}()
}

func (s *AuthService) SendCode(phone string) (time.Time, error) {
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	expiresAt := time.Now().Add(5 * time.Minute)
	vcode := &model.VerificationCode{
		Phone:     phone,
		Code:      code,
		ExpiresAt: expiresAt,
	}
	if err := s.verifRepo.Create(vcode); err != nil {
		return time.Time{}, err
	}
	s.sendSMS(phone, code)
	return expiresAt, nil
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
			now := time.Now()
			vipExpiry := now.Add(24 * time.Hour)
			if referrer.VipStatus == "active" && referrer.VipExpireAt != nil && referrer.VipExpireAt.After(now) {
				vipExpiry = referrer.VipExpireAt.Add(24 * time.Hour)
			}
			referrer.VipStatus = "active"
			referrer.VipExpireAt = &vipExpiry
			s.userRepo.Save(referrer)
		}
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, "", err
	}

	if inviteCode != "" {
		now := time.Now()
		vipExpiry := now.Add(24 * time.Hour)
		user.VipStatus = "active"
		user.VipExpireAt = &vipExpiry
		s.userRepo.Save(user)
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

func (s *AuthService) CreateQRSession() (string, error) {
	session := &model.QRSession{
		SessionID: repository.GenerateSessionID("qr"),
		Status:    "pending",
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	if err := s.verifRepo.CreateQRSession(session); err != nil {
		return "", err
	}
	return session.SessionID, nil
}

func (s *AuthService) GetQRStatus(sessionID string) (string, interface{}, interface{}) {
	session, err := s.verifRepo.FindQRSessionBySessionID(sessionID)
	if err != nil {
		return "expired", nil, nil
	}
	if time.Now().After(session.ExpiresAt) {
		return "expired", nil, nil
	}
	return session.Status, session.Token, session.ExpiresAt.Format(time.RFC3339)
}

func (s *AuthService) ScanQR(sessionID string) error {
	session, err := s.verifRepo.FindQRSessionBySessionID(sessionID)
	if err != nil {
		return errors.New("会话不存在")
	}
	if time.Now().After(session.ExpiresAt) {
		return errors.New("二维码已过期")
	}
	if session.Status != "pending" {
		return errors.New("二维码状态异常")
	}
	return s.verifRepo.UpdateQRSessionStatus(sessionID, "scanned", nil, "")
}

func (s *AuthService) AcceptTerms(userID uint) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	user.TermsAccepted = true
	return s.userRepo.Save(user)
}

func (s *AuthService) ConfirmQR(userID uint, sessionID string) (string, error) {
	session, err := s.verifRepo.FindQRSessionBySessionID(sessionID)
	if err != nil {
		return "", errors.New("会话不存在")
	}
	if time.Now().After(session.ExpiresAt) {
		return "", errors.New("二维码已过期")
	}
	if session.Status != "scanned" {
		return "", errors.New("请先扫描二维码")
	}
	token, err := s.GenerateToken(userID)
	if err != nil {
		return "", err
	}
	if err := s.verifRepo.UpdateQRSessionStatus(sessionID, "confirmed", &userID, token); err != nil {
		return "", err
	}
	return token, nil
}
