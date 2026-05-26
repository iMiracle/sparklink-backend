package handler

import (
	"sparklink-backend/pkg/response"
	"sparklink-backend/service"
	"github.com/gin-gonic/gin"
)

type RewardHandler struct {
	rewardService *service.RewardService
}

func NewRewardHandler(rewardService *service.RewardService) *RewardHandler {
	return &RewardHandler{rewardService: rewardService}
}

type ClaimRewardRequest struct {
	AdID   string `json:"adId" binding:"required"`
	AdType string `json:"adType" binding:"required"`
}

func (h *RewardHandler) Claim(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req ClaimRewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, response.ErrInvalidParams, "参数错误")
		return
	}
	reward, balance, cooldownEndsAt, err := h.rewardService.ClaimReward(userID, req.AdID, req.AdType)
	if err != nil {
		response.BadRequest(c, response.ErrInvalidParams, err.Error())
		return
	}
	response.Success(c, gin.H{
		"reward":         reward,
		"balance":        balance,
		"cooldownEndsAt": cooldownEndsAt,
	})
}

func (h *RewardHandler) GetCooldown(c *gin.Context) {
	userID := c.GetUint("user_id")
	adType := c.Query("adType")
	inCooldown, remaining, cooldownEnd, err := h.rewardService.GetCooldown(userID, adType)
	if err != nil {
		response.ServerError(c, "查询冷却信息失败")
		return
	}
	response.Success(c, gin.H{
		"inCooldown":      inCooldown,
		"remainingSeconds": remaining,
		"cooldownEndsAt":  cooldownEnd,
	})
}

func (h *RewardHandler) GetInviteInfo(c *gin.Context) {
	userID := c.GetUint("user_id")
	code, invitedCount, totalReward, shareChannels, err := h.rewardService.GetInviteInfo(userID)
	if err != nil {
		response.ServerError(c, "获取邀请信息失败")
		return
	}
	response.Success(c, gin.H{
		"inviteCode":     code,
		"invitedCount":   invitedCount,
		"totalReward":    totalReward,
		"shareChannels":  shareChannels,
	})
}
