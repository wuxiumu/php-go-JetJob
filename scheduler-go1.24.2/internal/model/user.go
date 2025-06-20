// internal/model/node.go
package model

import (
	"time"
)

// User 模型映射（字段与表结构一一对应）
type User struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	Name            string     `gorm:"size:255;not null" json:"name"`
	Email           string     `gorm:"size:255;uniqueIndex;not null" json:"email"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	Password        string     `gorm:"size:255;not null" json:"-"`
	RememberToken   *string    `gorm:"size:100" json:"remember_token"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}
