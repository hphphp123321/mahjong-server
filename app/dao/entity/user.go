package entity

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"type:varchar(255);unique;not null"`
	Password string `gorm:"type:varchar(255);not null"`
	Logs     []*Log `gorm:"many2many:user_logs;"` // 多对多关联
}

// 数据库表明自定义，默认为model的复数形式，比如这里默认为 users
func (User) TableName() string {
	return "users"
}
