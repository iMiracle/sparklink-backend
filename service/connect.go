package service

import (
	"errors"
	"time"

	"sparklink-backend/model"
	"sparklink-backend/repository"
)

type ConnectService struct {
	connRepo repository.ConnectRepository
	nodeRepo repository.NodeRepository
}

func NewConnectService(connRepo repository.ConnectRepository, nodeRepo repository.NodeRepository) *ConnectService {
	return &ConnectService{
		connRepo: connRepo,
		nodeRepo: nodeRepo,
	}
}

func (s *ConnectService) Start(userID uint, nodeID, protocol string) (*model.ConnectSession, *model.Node, error) {
	node, err := s.nodeRepo.FindByNodeID(nodeID)
	if err != nil {
		return nil, nil, errors.New("节点不存在")
	}

	existing, _ := s.connRepo.FindActiveSession(userID)
	if existing != nil {
		existing.Status = "stopped"
		now := time.Now()
		existing.StoppedAt = &now
		s.connRepo.UpdateSession(existing)
	}

	session := &model.ConnectSession{
		SessionID: repository.GenerateSessionID("conn"),
		UserID:    userID,
		NodeID:    nodeID,
		Protocol:  protocol,
		Status:    "active",
		StartedAt: time.Now(),
	}
	if err := s.connRepo.CreateSession(session); err != nil {
		return nil, nil, err
	}

	return session, node, nil
}

func (s *ConnectService) Stop(userID uint, sessionID string) error {
	session, err := s.connRepo.FindActiveSession(userID)
	if err != nil {
		return errors.New("无活跃连接")
	}
	if sessionID != "" && session.SessionID != sessionID {
		return errors.New("会话ID不匹配")
	}
	session.Status = "stopped"
	now := time.Now()
	session.StoppedAt = &now
	return s.connRepo.UpdateSession(session)
}

func (s *ConnectService) Report(userID uint, sessionID, status string, duration int) error {
	session, err := s.connRepo.FindActiveSession(userID)
	if err != nil {
		return errors.New("无活跃连接")
	}
	if session.SessionID != sessionID {
		return errors.New("会话ID不匹配")
	}
	session.Status = status
	if status == "error" || status == "stopped" {
		now := time.Now()
		session.StoppedAt = &now
	}
	if duration > 0 {
		session.BytesSent = int64(duration)
	}
	return s.connRepo.UpdateSession(session)
}
