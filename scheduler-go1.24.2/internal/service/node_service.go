package service

import (
	"jetjob/internal/model"
	"jetjob/internal/storage"
	"time"
)

// NodeService 提供节点相关业务方法
type NodeService struct{}

// UpdateNodeHeartbeat 更新节点心跳与负载
func (s *NodeService) UpdateNodeHeartbeat(name string, load float64) error {
	var node model.Node
	if err := storage.DB.Where("name = ?", name).First(&node).Error; err != nil {
		return err
	}
	node.LastHeartbeat = time.Now()
	node.LoadAvg = load
	node.Status = "active"
	return storage.DB.Save(&node).Error
}

// OfflineTimeoutNodes 超时未上报的节点批量下线，返回下线节点列表
func (s *NodeService) OfflineTimeoutNodes(timeout time.Duration) ([]model.Node, error) {
	var nodes []model.Node
	storage.DB.Find(&nodes)
	now := time.Now()
	var offlined []model.Node
	for _, node := range nodes {
		if node.Status == "active" && now.Sub(node.LastHeartbeat) > timeout {
			storage.DB.Model(&model.Node{}).Where("id = ?", node.ID).Update("status", "offline")
			offlined = append(offlined, node)
		}
	}
	return offlined, nil
}
