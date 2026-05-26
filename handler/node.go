package handler

import (
	"strings"

	"sparklink-backend/pkg/response"
	"sparklink-backend/service"
	"github.com/gin-gonic/gin"
)

type NodeHandler struct {
	nodeService *service.NodeService
}

func NewNodeHandler(nodeService *service.NodeService) *NodeHandler {
	return &NodeHandler{nodeService: nodeService}
}

type SpeedtestRequest struct {
	NodeID string `json:"nodeId" binding:"required"`
}

func (h *NodeHandler) List(c *gin.Context) {
	protocol := c.Query("protocol")
	visibility := c.Query("type")
	region := c.Query("region")

	nodes, err := h.nodeService.GetNodes(protocol, visibility, region)
	if err != nil {
		response.ServerError(c, "获取节点列表失败")
		return
	}

	var result []gin.H
	for _, n := range nodes {
		result = append(result, gin.H{
			"nodeId":          n.NodeId,
			"name":            n.Name,
			"protocol":        n.Protocol,
			"latency":         n.Latency,
			"load":            n.Load,
			"regionCode":      n.RegionCode,
			"regionName":      n.RegionName,
			"tags":            parseTags(n.Tags),
			"visibilityLevel": n.VisibilityLevel,
			"priority":        n.Priority,
		})
	}
	if result == nil {
		result = []gin.H{}
	}
	response.Success(c, gin.H{
		"nodes": result,
		"total": len(result),
	})
}

func (h *NodeHandler) Get(c *gin.Context) {
	nodeID := c.Param("nodeId")
	node, err := h.nodeService.GetNode(nodeID)
	if err != nil {
		response.NotFound(c, "节点不存在")
		return
	}
	response.Success(c, gin.H{
		"nodeId":          node.NodeId,
		"name":            node.Name,
		"protocols":       parseTags(node.Protocols),
		"host":            node.Host,
		"port":            node.Port,
		"publicKey":       node.PublicKey,
		"latency":         node.Latency,
		"load":            node.Load,
		"regionCode":      node.RegionCode,
		"regionName":      node.RegionName,
		"tags":            parseTags(node.Tags),
		"distance":        node.Distance,
		"bandwidthLimit":  node.BandwidthLimit,
		"visibilityLevel": node.VisibilityLevel,
		"connectionStats": gin.H{
			"todayConnections": 128,
			"successRate":      98.5,
			"avgLatency":       42,
			"avgSpeed":         85.2,
		},
	})
}

func (h *NodeHandler) Speedtest(c *gin.Context) {
	var req SpeedtestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, response.ErrInvalidParams, "参数错误")
		return
	}
	download, upload, latency, err := h.nodeService.Speedtest(req.NodeID)
	if err != nil {
		response.BadRequest(c, response.ErrInvalidParams, "节点不存在")
		return
	}
	response.Success(c, gin.H{
		"downloadSpeed": download,
		"uploadSpeed":   upload,
		"latency":       latency,
	})
}

func (h *NodeHandler) UpdatePing(c *gin.Context) {
	var req struct {
		NodeID  string `json:"nodeId" binding:"required"`
		Latency int    `json:"latency" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, response.ErrInvalidParams, "参数错误")
		return
	}
	if err := h.nodeService.UpdatePing(req.NodeID, req.Latency); err != nil {
		response.BadRequest(c, response.ErrInvalidParams, "节点不存在")
		return
	}
	response.Success(c, gin.H{"message": "ping updated"})
}

func parseTags(tagsStr string) []string {
	if tagsStr == "" {
		return []string{}
	}
	return strings.Split(tagsStr, ",")
}
