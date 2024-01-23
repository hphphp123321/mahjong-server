package entity

import "gorm.io/gorm"

// Log 对应 logs 表
type Log struct {
	gorm.Model
	ID      uint    `gorm:"primaryKey"`
	Content string  `gorm:"type:json;not null"`
	Users   []*User `gorm:"many2many:user_logs;"` // 多对多关联
}
