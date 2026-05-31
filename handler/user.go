package handler

import (
	"fmt"
	"time"

	"sparklink-backend/pkg/response"
	"sparklink-backend/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
	subService  *service.SubscriptionService
}

func NewUserHandler(userService *service.UserService, subService *service.SubscriptionService) *UserHandler {
	return &UserHandler{userService: userService, subService: subService}
}

func nullTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

func (h *UserHandler) Profile(c *gin.Context) {
	userID := c.GetUint("user_id")
	user, err := h.userService.GetProfile(userID)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}
	subInfo := gin.H{}
	sub, err := h.subService.FindActive(userID)
	if err == nil && sub != nil {
		subInfo = gin.H{
			"planId":     sub.PlanID,
			"expiresAt":  sub.ExpireTime.Format(time.RFC3339),
			"status":     sub.Status,
		}
	}
	response.Success(c, gin.H{
		"userId":            fmt.Sprintf("u_%d", user.ID),
		"phone":             user.Phone,
		"nickname":          user.Nickname,
		"avatar":            user.Avatar,
		"vipStatus":         user.VipStatus,
		"vipExpiresAt":      nullTime(user.VipExpireAt),
		"balanceMinutes":    user.BalanceMinutes,
		"inviteCode":        user.InviteCode,
		"invitedCount":      user.InvitedCount,
		"killSwitchEnabled": user.KillSwitchEnabled,
		"termsAccepted":     user.TermsAccepted,
		"subscription":      subInfo,
		"registeredAt":      user.CreatedAt.Format(time.RFC3339),
	})
}

func (h *UserHandler) Devices(c *gin.Context) {
	userID := c.GetUint("user_id")
	devices, err := h.userService.GetDevices(userID)
	if err != nil {
		response.ServerError(c, "获取设备列表失败")
		return
	}
	maxDevices, err := h.userService.MaxDevicesForUser(userID)
	if err != nil {
		response.ServerError(c, "获取设备上限失败")
		return
	}
	var result []gin.H
	for _, d := range devices {
		result = append(result, gin.H{
			"deviceId":   d.DeviceID,
			"deviceName": d.DeviceName,
			"platform":   d.DeviceType,
			"lastActive": d.LastActive.Format(time.RFC3339),
			"isOnline":   d.IsActive,
		})
	}
	if result == nil {
		result = []gin.H{}
	}
	response.Success(c, gin.H{
		"devices":    result,
		"total":      len(result),
		"maxDevices": maxDevices,
	})
}

func (h *UserHandler) RegisterDevice(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		DeviceID   string `json:"deviceId" binding:"required"`
		DeviceName string `json:"deviceName"`
		Platform   string `json:"platform"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, response.ErrInvalidParams, "请提供 deviceId")
		return
	}
	if err := h.userService.RegisterDevice(userID, req.DeviceID, req.DeviceName, req.Platform); err != nil {
		response.BadRequest(c, response.ErrDeviceLimit, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "设备注册成功"})
}

func (h *UserHandler) RemoveDevice(c *gin.Context) {
	deviceID := c.Param("deviceId")
	if err := h.userService.RemoveDevice(deviceID); err != nil {
		response.NotFound(c, "设备不存在")
		return
	}
	response.Success(c, gin.H{"message": "设备已移除"})
}

func (h *UserHandler) KillSwitchStatus(c *gin.Context) {
	userID := c.GetUint("user_id")
	enabled, err := h.userService.KillSwitchStatus(userID)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}
	response.Success(c, gin.H{"killSwitchEnabled": enabled})
}

func (h *UserHandler) ToggleKillSwitch(c *gin.Context) {
	userID := c.GetUint("user_id")
	enabled, err := h.userService.ToggleKillSwitch(userID)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}
	response.Success(c, gin.H{"killSwitchEnabled": enabled})
}
