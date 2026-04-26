package handler

import (
	"net/http"

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
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "plans": plans})
}

type CreateSubscriptionRequest struct {
	Plan   string  `json:"plan" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
}

func (h *SubscriptionHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")

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
}

func (h *SubscriptionHandler) Verify(c *gin.Context) {
	userID := c.GetUint("user_id")

	valid, _ := h.subService.VerifySubscription(userID)

	c.JSON(http.StatusOK, gin.H{"success": true, "valid": valid})
}