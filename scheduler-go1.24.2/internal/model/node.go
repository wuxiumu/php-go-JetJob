// internal/model/node.go
package model

import (
	"time"
)

// Node 结构体，对应数据库表结构
type Node struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	Name          string     `gorm:"uniqueIndex;size:64;not null" json:"name"` // 节点名唯一
	Group         string     `gorm:"size:64;default:'default'" json:"group"`   // 新增字段
	Host          string     `gorm:"size:255" json:"host"`
	LastHeartbeat time.Time  `json:"last_heartbeat"`
	Status        string     `gorm:"size:16;default:'active'" json:"status"` // active, offline
	LoadAvg       float64    `json:"load_avg"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
}

// 节点在线判断
func (n *Node) IsOnline(timeout time.Duration) bool {
	return time.Since(n.LastHeartbeat) < timeout && n.Status == "active"
}
