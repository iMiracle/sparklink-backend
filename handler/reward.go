package handler

import (
	"net/http"

	"sparklink-backend/service"

	"github.com/gin-gonic/gin"
)

type RewardHandler struct {
	rewardService *service.RewardService
}

func NewRewardHandler(rewardService *service.RewardService) *RewardHandler {
	return &RewardHandler{rewardService: rewardService}
}

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
}

func (h *RewardHandler) GetInviteInfo(c *gin.Context) {
	userID := c.GetUint("user_id")

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