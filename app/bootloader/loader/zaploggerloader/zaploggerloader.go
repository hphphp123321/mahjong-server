package zaploggerloader

import (
	"context"
	"github.com/hphphp123321/mahjong-server/app/component/zap"
	"github.com/hphphp123321/mahjong-server/app/global"
)

// ZapLoggerLoader 相关依赖：global.Log
// 1. 负责初始化zap logger
type ZapLoggerLoader struct {
}

func (z ZapLoggerLoader) Load(ctx context.Context, env map[string]string) error {
	var err error
	c := global.C.LogConfig
	mode := global.C.App.Mode

	cfg := &zap.LogConfig{
		Level:       c.Level,
		Filename:    c.Filename,
		MaxSize:     c.MaxSize,
		MaxAge:      c.MaxAge,
		MaxBackups:  c.MaxBackups,
		OutputFile:  false,
		OutputStdio: false,
	}
	if mode == "prod" {
		cfg.OutputFile = true
		cfg.OutputStdio = false
		cfg.Color = false
	} else if mode == "dev" {
		cfg.OutputFile = false
		cfg.OutputStdio = true
		cfg.Color = true
	} else {
		cfg.OutputFile = false
		cfg.OutputStdio = false
	}

	global.Log, err = zap.InitLogger(cfg)

	if err != nil {
		return err
	}

	return nil
}

func (z ZapLoggerLoader) Name() string {
	return "ZapLoggerLoader"
}

func (z ZapLoggerLoader) Require() []string {
	return []string{"ConfigLoader"}
}
