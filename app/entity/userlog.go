package entity

// UserLog 对应 user_logs 表
type UserLog struct {
	UserID uint `gorm:"primaryKey"`
	LogID  uint `gorm:"primaryKey"`
}
