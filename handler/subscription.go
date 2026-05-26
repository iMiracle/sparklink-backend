package handler

import (
<<<<<<< HEAD
	"time"
	"sparklink-backend/pkg/response"
	"sparklink-backend/service"
=======
	"net/http"

	"sparklink-backend/service"

>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	subService *service.SubscriptionService
}

func NewSubscriptionHandler(subService *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{subService: subService}
}

func (h *SubscriptionHandler) ListPlans(c *gin.Context) {
	plans, err := h.subService.GetPlans()
	if err != nil {
<<<<<<< HEAD
		response.ServerError(c, "获取套餐列表失败")
		return
	}
	var result []gin.H
	for _, p := range plans {
		result = append(result, gin.H{
			"planId":      p.PlanID,
			"name":        p.Name,
			"duration":    p.DurationDays,
			"price":       p.Price,
			"originalPrice": p.OriginalPrice,
			"popular":     p.Popular,
			"features":    p.Features,
		})
	}
	if result == nil {
		result = []gin.H{}
	}
	response.Success(c, gin.H{"plans": result})
=======
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "plans": plans})
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
}

type CreateSubscriptionRequest struct {
	Plan   string  `json:"plan" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
}

func (h *SubscriptionHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")
<<<<<<< HEAD
	var req CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, response.ErrInvalidParams, "参数错误")
		return
	}
	sub, err := h.subService.CreateSubscription(userID, req.Plan, req.Amount)
	if err != nil {
		response.ServerError(c, "创建订阅失败")
		return
	}
	response.Success(c, gin.H{
		"subscription": gin.H{
			"id":         sub.ID,
			"planId":     sub.PlanID,
			"amount":     sub.Amount,
			"startTime":  sub.StartTime.Format(time.RFC3339),
			"expireTime": sub.ExpireTime.Format(time.RFC3339),
			"status":     sub.Status,
		},
	})
=======

	var req CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request"})
		return
	}

	sub, err := h.subService.CreateSubscription(userID, req.Plan, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "subscription": sub})
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
}

func (h *SubscriptionHandler) Verify(c *gin.Context) {
	userID := c.GetUint("user_id")
<<<<<<< HEAD
	valid, _ := h.subService.VerifySubscription(userID)
	response.Success(c, gin.H{"valid": valid})
}
=======

	valid, _ := h.subService.VerifySubscription(userID)

	c.JSON(http.StatusOK, gin.H{"success": true, "valid": valid})
}
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
