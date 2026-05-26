package handler

import (
	"time"
	"sparklink-backend/pkg/response"
	"sparklink-backend/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
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
	response.Success(c, gin.H{
		"userId":         user.ID,
		"phone":          user.Phone,
		"nickname":       user.Nickname,
		"avatar":         user.Avatar,
		"vipStatus":      user.VipStatus,
		"vipExpiresAt":   nullTime(user.VipExpireAt),
		"balanceMinutes": user.BalanceMinutes,
		"inviteCode":     user.InviteCode,
		"invitedCount":   user.InvitedCount,
		"registeredAt":   user.CreatedAt.Format(time.RFC3339),
	})
}

func (h *UserHandler) Devices(c *gin.Context) {
	userID := c.GetUint("user_id")
	devices, err := h.userService.GetDevices(userID)
	if err != nil {
		response.ServerError(c, "获取设备列表失败")
		return
	}
	var result []gin.H
	for _, d := range devices {
		result = append(result, gin.H{
			"deviceId":   d.DeviceID,
			"deviceName": d.DeviceName,
			"deviceType": d.DeviceType,
			"lastActive": d.LastActive.Format(time.RFC3339),
			"isActive":   d.IsActive,
		})
	}
	if result == nil {
		result = []gin.H{}
	}
	response.Success(c, gin.H{"devices": result})
}

func (h *UserHandler) RemoveDevice(c *gin.Context) {
	deviceID := c.Param("deviceId")
	if err := h.userService.RemoveDevice(deviceID); err != nil {
		response.NotFound(c, "设备不存在")
		return
	}
	response.Success(c, gin.H{"message": "设备已移除"})
}
