package service

import (
<<<<<<< HEAD
=======
	"math/rand"
	"time"

>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	"sparklink-backend/model"
	"sparklink-backend/repository"
)

type NodeService struct {
<<<<<<< HEAD
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
=======
	nodeRepo *repository.NodeRepository
}

func NewNodeService(nodeRepo *repository.NodeRepository) *NodeService {
	return &NodeService{nodeRepo: nodeRepo}
}

func (s *NodeService) GetNodes(protocol, nodeType, country string) ([]model.Node, error) {
	return s.nodeRepo.FindAll(protocol, nodeType, country)
}

func (s *NodeService) GetNode(id uint) (*model.Node, error) {
	return s.nodeRepo.FindByID(id)
}

func (s *NodeService) GetOptimalNode(protocol, nodeType string) (*model.Node, error) {
	nodes, err := s.nodeRepo.FindAll(protocol, nodeType, "")
	if err != nil {
		return nil, err
	}

	if len(nodes) == 0 {
		return nil, nil
	}

	// 根据权重选择最优节点
	return s.weightedSelect(nodes), nil
}

func (s *NodeService) UpdatePing(id uint, latency int) error {
	return s.nodeRepo.UpdateLatency(id, latency)
}

func (s *NodeService) UpdateLoad(id uint, load int) error {
	return s.nodeRepo.UpdateLoad(id, load)
}

func (s *NodeService) weightedSelect(nodes []model.Node) *model.Node {
	var totalWeight int
	for _, n := range nodes {
		nodeWeight := 100 - n.Load
		if nodeWeight < 0 {
			nodeWeight = 0
		}
		totalWeight += nodeWeight
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	target := r.Intn(totalWeight)

	var cumulative int
	for _, n := range nodes {
		nodeWeight := 100 - n.Load
		if nodeWeight < 0 {
			nodeWeight = 0
		}
		cumulative += nodeWeight
		if cumulative >= target {
			return &n
		}
	}

	return &nodes[0]
}
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
