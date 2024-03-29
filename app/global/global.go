package global

import (
	"github.com/hphphp123321/mahjong-server/app/component/config"
	"github.com/hphphp123321/mahjong-server/app/service/robot"
	"go.uber.org/zap"
)

var (
	// ProjectRoot [BaseLoader]项目根目录
	ProjectRoot string

	// C [ConfigLoader]app.yaml全局配置实例，这里做初始化是为了Decode不产生Error
	C *config.Config = &config.Config{}

	// Log [ZapLoggerLoader]zap-logger实例
	Log *zap.SugaredLogger

	// RobotRegistry [RobotLoader]机器人注册表
	RobotRegistry *robot.Registry = robot.NewRegistry()
)
