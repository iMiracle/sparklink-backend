package handler

import (
	"net/http"

	"sparklink-backend/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type SendCodeRequest struct {
	Phone string `json:"phone" binding:"required"`
}

func (h *AuthHandler) SendCode(c *gin.Context) {
	var req SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "phone required"})
		return
	}

	// TODO: 发送验证码
	// 这里需要接入短信服务

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "code sent"})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Phone   string `json:"phone"`
	Password string `json:"password"`
	DeviceID string `json:"device_id" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request"})
		return
	}

	user, token, err := h.authService.Login(req.Email, req.Phone, req.Password, req.DeviceID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"token": token,
		"user": gin.H{
			"id":            user.ID,
			"email":         user.Email,
			"nickname":      user.Nickname,
			"ad_credits":   user.AdCredits,
			"expire_time":  user.ExpireTime,
			"referral_code": user.ReferralCode,
		},
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	userID := c.GetUint("user_id")

	token, err := h.authService.GenerateToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "token": token})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "logged out"})
}