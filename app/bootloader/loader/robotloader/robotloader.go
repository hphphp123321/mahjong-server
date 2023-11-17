package robotloader

import (
	"context"

	"github.com/hphphp123321/mahjong-server/app/global"
	"github.com/hphphp123321/mahjong-server/app/service/robot"
)

type RobotLoader struct {
}

func (loader *RobotLoader) Load(ctx context.Context, env map[string]string) error {
	global.RobotRegistry.Register(&robot.SimpleRobot{})

	var conf = global.C.Openai
	if conf.Key != "" && conf.Model != "" {
		var chatgptRobot = &robot.ChatGPTRobot{
			BaseUrl:  conf.BaseURL,
			Key:      conf.Key,
			Model:    conf.Model,
			Lang:     conf.Lang,
			ProxyUrl: conf.ProxyUrl,
		}
		global.RobotRegistry.Register(chatgptRobot)
		global.Log.Infoln("chatgpt robot registered!")

	}
	return nil
}

func (loader *RobotLoader) Name() string {
	return "Robot Register Loader"
}

func (loader *RobotLoader) Require() []string {
	return []string{"ZapLoggerLoader"}
}
