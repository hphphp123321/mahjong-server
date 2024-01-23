package dbloader

import (
	"context"
	"github.com/hphphp123321/mahjong-server/app/dao"
	"github.com/hphphp123321/mahjong-server/app/entity"
	"github.com/hphphp123321/mahjong-server/app/query"
)

type DBLoader struct {
}

func (D DBLoader) Load(ctx context.Context, env map[string]string) error {
	if err := dao.InitMySql(); err != nil {
		return err
	}

	// 自动迁移
	dao.GormDB.AutoMigrate(&entity.User{})
	dao.GormDB.AutoMigrate(&entity.Log{})
	dao.GormDB.AutoMigrate(&entity.UserLog{})

	// Query初始化
	query.SetDefault(dao.GormDB)

	return nil
}

func (D DBLoader) Name() string {
	return "DBLoader"
}

func (D DBLoader) Require() []string {
	return []string{"RobotLoader"}
}
