package handler

import (
	"net/http"
	"strconv"

	"sparklink-backend/service"

	"github.com/gin-gonic/gin"
)

type NodeHandler struct {
	nodeService *service.NodeService
}

func NewNodeHandler(nodeService *service.NodeService) *NodeHandler {
	return &NodeHandler{nodeService: nodeService}
}

func (h *NodeHandler) List(c *gin.Context) {
	protocol := c.Query("protocol")
	nodeType := c.Query("type")
	country := c.Query("region")

	nodes, err := h.nodeService.GetNodes(protocol, nodeType, country)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"nodes": nodes,
		"total": len(nodes),
	})
}

func (h *NodeHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid id"})
		return
	}

	node, err := h.nodeService.GetNode(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "node not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "node": node})
}

type PingRequest struct {
	NodeID  uint `json:"node_id" binding:"required"`
	Latency int  `json:"latency" binding:"required"`
}

func (h *NodeHandler) UpdatePing(c *gin.Context) {
	var req PingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request"})
		return
	}

	h.nodeService.UpdatePing(req.NodeID, req.Latency)

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *NodeHandler) Favorites(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "nodes": []}})
}

type FavoriteRequest struct {
	NodeID uint `json:"node_id" binding:"required"`
}

func (h *NodeHandler) AddFavorite(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *NodeHandler) RemoveFavorite(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true})
}