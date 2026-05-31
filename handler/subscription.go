package handler

import (
	"fmt"
	"time"

	"sparklink-backend/pkg/response"
	"sparklink-backend/service"

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
		response.ServerError(c, "获取套餐列表失败")
		return
	}
	var result []gin.H
	for _, p := range plans {
		item := gin.H{
			"planId":        p.PlanID,
			"name":          p.Name,
			"duration":      p.DurationDays,
			"price":         p.Price,
			"originalPrice": p.OriginalPrice,
			"currency":      p.Currency,
			"popular":       p.Popular,
			"tag":           p.Tag,
			"features":      p.Features,
		}
		if p.DailyPrice != nil {
			item["dailyPrice"] = *p.DailyPrice
		}
		result = append(result, item)
	}
	if result == nil {
		result = []gin.H{}
	}
	response.Success(c, gin.H{"plans": result})
}

type CreateOrderRequest struct {
	PlanID string `json:"planId" binding:"required"`
}

func (h *SubscriptionHandler) CreateOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, response.ErrInvalidParams, "参数错误")
		return
	}
	sub, err := h.subService.CreateSubscription(userID, req.PlanID, 0)
	if err != nil {
		response.ServerError(c, "创建订阅失败")
		return
	}
	response.Success(c, gin.H{
		"orderId":   fmt.Sprintf("ord_%d", sub.ID),
		"amount":    sub.Amount,
		"status":    sub.Status,
		"createdAt": sub.CreatedAt.Format(time.RFC3339),
	})
}

func (h *SubscriptionHandler) Verify(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		Receipt string `json:"receipt" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, response.ErrInvalidParams, "参数错误")
		return
	}
	_, err := h.subService.VerifySubscription(userID)
	verified := err == nil
	resp := gin.H{"verified": verified}
	if verified {
		sub, _ := h.subService.FindActive(userID)
		if sub != nil {
			resp["vipExpiresAt"] = sub.ExpireTime.Format(time.RFC3339)
		}
	}
	response.Success(c, resp)
}
