package handler

import (
	"time"
	"sparklink-backend/pkg/response"
	"sparklink-backend/service"
	"github.com/gin-gonic/gin"
)

type ConnectHandler struct {
	connectService *service.ConnectService
}

func NewConnectHandler(connectService *service.ConnectService) *ConnectHandler {
	return &ConnectHandler{connectService: connectService}
}

type StartConnectRequest struct {
	NodeID   string `json:"nodeId" binding:"required"`
	Protocol string `json:"protocol" binding:"required"`
}

type StopConnectRequest struct {
	SessionID string `json:"sessionId"`
}

type ReportConnectRequest struct {
	SessionID string `json:"sessionId" binding:"required"`
	Status    string `json:"status" binding:"required"`
	Duration  int    `json:"duration"`
}

func (h *ConnectHandler) Start(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req StartConnectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, response.ErrInvalidParams, "参数错误")
		return
	}
	session, node, err := h.connectService.Start(userID, req.NodeID, req.Protocol)
	if err != nil {
		response.BadRequest(c, response.ErrInvalidParams, err.Error())
		return
	}
	config := gin.H{
		"host":      node.Host,
		"port":      node.Port,
		"publicKey": node.PublicKey,
		"protocol":  node.Protocol,
	}
	if node.Protocols != "" {
		config["protocols"] = parseTags(node.Protocols)
	}
	response.Success(c, gin.H{
		"sessionId": session.SessionID,
		"config":    config,
		"expiresAt": time.Now().Add(1 * time.Hour).Format(time.RFC3339),
	})
}

func (h *ConnectHandler) Stop(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req StopConnectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, response.ErrInvalidParams, "参数错误")
		return
	}
	if err := h.connectService.Stop(userID, req.SessionID); err != nil {
		response.BadRequest(c, response.ErrInvalidParams, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "连接已断开"})
}

func (h *ConnectHandler) Report(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req ReportConnectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, response.ErrInvalidParams, "参数错误")
		return
	}
	if err := h.connectService.Report(userID, req.SessionID, req.Status, req.Duration); err != nil {
		response.BadRequest(c, response.ErrInvalidParams, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "报告已接收"})
}
