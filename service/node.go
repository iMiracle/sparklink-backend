package service

import (
	"sparklink-backend/model"
	"sparklink-backend/repository"
)

type NodeService struct {
	nodeRepo repository.NodeRepository
}

func NewNodeService(nodeRepo repository.NodeRepository) *NodeService {
	return &NodeService{nodeRepo: nodeRepo}
}

func (s *NodeService) GetNodes(protocol, visibility, region string) ([]model.Node, error) {
	return s.nodeRepo.FindAll(protocol, visibility, region)
}

func (s *NodeService) GetNode(nodeID string) (*model.Node, error) {
	return s.nodeRepo.FindByNodeID(nodeID)
}

func (s *NodeService) Speedtest(nodeID string) (downloadSpeed, uploadSpeed float64, latency int, err error) {
	node, err := s.nodeRepo.FindByNodeID(nodeID)
	if err != nil {
		return 0, 0, 0, err
	}
	return 85.5, 42.3, node.Latency, nil
}

func (s *NodeService) UpdatePing(nodeID string, latency int) error {
	return s.nodeRepo.UpdatePing(nodeID, latency)
}
