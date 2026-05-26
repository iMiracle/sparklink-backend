package handler

import (
<<<<<<< HEAD
	"fmt"
	"time"
	"sparklink-backend/pkg/response"
	"sparklink-backend/service"
=======
	"net/http"

	"sparklink-backend/service"

>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
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

<<<<<<< HEAD
type VerifyCodeRequest struct {
	Phone      string `json:"phone" binding:"required"`
	Code       string `json:"code" binding:"required"`
	InviteCode string `json:"inviteCode"`
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func (h *AuthHandler) SendCode(c *gin.Context) {
	var req SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, response.ErrInvalidParams, "手机号不能为空")
		return
	}
	if err := h.authService.SendCode(req.Phone); err != nil {
		response.ServerError(c, "验证码发送失败")
		return
	}
	response.Success(c, gin.H{"message": "验证码已发送"})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req VerifyCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, response.ErrInvalidParams, "参数错误")
		return
	}
	user, token, err := h.authService.Register(req.Phone, req.Code, req.InviteCode)
	if err != nil {
		response.BadRequest(c, response.ErrInvalidParams, err.Error())
		return
	}
	response.Success(c, gin.H{
		"userId":     user.ID,
		"inviteCode": user.InviteCode,
		"token":      token,
		"expiresAt":  h.authService.GetTokenExpiry().Format(time.RFC3339),
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req VerifyCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, response.ErrInvalidParams, "参数错误")
		return
	}
	_, token, err := h.authService.Login(req.Phone, req.Code)
	if err != nil {
		response.BadRequest(c, response.ErrInvalidParams, err.Error())
		return
	}
	response.Success(c, gin.H{
		"token":     token,
		"expiresAt": h.authService.GetTokenExpiry().Format(time.RFC3339),
=======
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
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	userID := c.GetUint("user_id")
<<<<<<< HEAD
	token, err := h.authService.GenerateToken(userID)
	if err != nil {
		response.ServerError(c, "令牌刷新失败")
		return
	}
	response.Success(c, gin.H{"token": token})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	response.Success(c, gin.H{"message": "已退出登录"})
}

func (h *AuthHandler) QrCode(c *gin.Context) {
	id := generateID()
	response.Success(c, gin.H{
		"sessionId": "qr_sess_" + id,
		"qrData":    "sparklink://auth?session=qr_sess_" + id,
	})
}

func (h *AuthHandler) QrCodeStatus(c *gin.Context) {
	sessionID := c.Param("sessionId")
	status, token, expiresAt := h.authService.GetQRStatus(sessionID)
	response.Success(c, gin.H{
		"status":    status,
		"token":     token,
		"expiresAt": expiresAt,
	})
}
=======

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
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
