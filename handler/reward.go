package handler

import (
<<<<<<< HEAD
	"sparklink-backend/pkg/response"
	"sparklink-backend/service"
=======
	"net/http"

	"sparklink-backend/service"

>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	"github.com/gin-gonic/gin"
)

type RewardHandler struct {
	rewardService *service.RewardService
}

func NewRewardHandler(rewardService *service.RewardService) *RewardHandler {
	return &RewardHandler{rewardService: rewardService}
}

<<<<<<< HEAD
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
=======
type VideoRewardRequest struct {
	AdPlatform     string `json:"ad_platform" binding:"required"`
	AdID           string `json:"ad_id" binding:"required"`
	TransactionID string `json:"transaction_id" binding:"required"`
	Nonce          string `json:"nonce" binding:"required"`
}

func (h *RewardHandler) VideoReward(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req VideoRewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request"})
		return
	}

	err := h.rewardService.ClaimVideoReward(userID, req.AdPlatform, req.AdID, req.TransactionID, req.Nonce)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "reward": 120})
}

func (h *RewardHandler) DailyCheckin(c *gin.Context) {
	userID := c.GetUint("user_id")

	reward, err := h.rewardService.DailyCheckin(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "reward": reward})
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
}

func (h *RewardHandler) GetInviteInfo(c *gin.Context) {
	userID := c.GetUint("user_id")
<<<<<<< HEAD
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
=======

	code, reward, _ := h.rewardService.GetInviteInfo(userID)

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"code":        code,
		"reward_minutes": reward,
	})
}

type BindInviteRequest struct {
	Code string `json:"code" binding:"required"`
}

func (h *RewardHandler) BindInvite(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req BindInviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "code required"})
		return
	}

	err := h.rewardService.BindInvite(userID, req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "reward": 1440})
}

func (h *RewardHandler) GetCredits(c *gin.Context) {
	userID := c.GetUint("user_id")

	// 这里需要从 userRepo 获取
	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"credits": 0,
	})
}
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
