package handler

import (
	"fmt"
	"strings"
	"time"

	"sparklink-backend/middleware"
	"sparklink-backend/pkg/response"
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
	Type  string `json:"type"`
}

type VerifyCodeRequest struct {
	Phone      string `json:"phone" binding:"required"`
	Code       string `json:"code" binding:"required"`
	InviteCode string `json:"inviteCode"`
}

func (h *AuthHandler) SendCode(c *gin.Context) {
	var req SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, response.ErrInvalidParams, "手机号不能为空")
		return
	}
	expiresAt, err := h.authService.SendCode(req.Phone)
	if err != nil {
		response.ServerError(c, "验证码发送失败")
		return
	}
	response.Success(c, gin.H{
		"message":   "验证码已发送",
		"expiresAt": expiresAt.Format(time.RFC3339),
	})
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
		"userId":     fmt.Sprintf("u_%d", user.ID),
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
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	userID := c.GetUint("user_id")
	token, err := h.authService.GenerateToken(userID)
	if err != nil {
		response.ServerError(c, "令牌刷新失败")
		return
	}
	response.Success(c, gin.H{"token": token})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	tokenStr := authHeader
	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenStr = authHeader[7:]
	}
	if tokenStr != "" {
		middleware.BlacklistToken(tokenStr)
	}
	response.Success(c, gin.H{"message": "已退出登录"})
}

func (h *AuthHandler) QrCode(c *gin.Context) {
	sessionID, err := h.authService.CreateQRSession()
	if err != nil {
		response.ServerError(c, "创建二维码失败")
		return
	}
	response.Success(c, gin.H{
		"sessionId": sessionID,
		"qrData":    "sparklink://auth?session=" + sessionID,
	})
}

func (h *AuthHandler) QrCodeStatus(c *gin.Context) {
	sessionID := c.Query("sessionId")
	if sessionID == "" {
		response.BadRequest(c, response.ErrInvalidParams, "请提供 sessionId")
		return
	}
	status, token, expiresAt := h.authService.GetQRStatus(sessionID)
	response.Success(c, gin.H{
		"status":    status,
		"token":     token,
		"expiresAt": expiresAt,
	})
}

func (h *AuthHandler) QrScan(c *gin.Context) {
	var req struct {
		SessionID string `json:"sessionId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, response.ErrInvalidParams, "请提供 sessionId")
		return
	}
	if err := h.authService.ScanQR(req.SessionID); err != nil {
		response.BadRequest(c, response.ErrInvalidParams, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "扫码成功"})
}

func (h *AuthHandler) QrConfirm(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		SessionID string `json:"sessionId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, response.ErrInvalidParams, "请提供 sessionId")
		return
	}
	token, err := h.authService.ConfirmQR(userID, req.SessionID)
	if err != nil {
		response.BadRequest(c, response.ErrInvalidParams, err.Error())
		return
	}
	response.Success(c, gin.H{"token": token})
}

func (h *AuthHandler) AcceptTerms(c *gin.Context) {
	userID := c.GetUint("user_id")
	if err := h.authService.AcceptTerms(userID); err != nil {
		response.NotFound(c, "用户不存在")
		return
	}
	response.Success(c, gin.H{"message": "已同意服务条款"})
}
