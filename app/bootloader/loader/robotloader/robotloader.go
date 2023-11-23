package robotloader

import (
	"context"
	"github.com/hphphp123321/mahjong-server/app/service/robot/chatgpt"

	"github.com/hphphp123321/mahjong-server/app/global"
	"github.com/hphphp123321/mahjong-server/app/service/robot"
)

type RobotLoader struct {
}

func (loader *RobotLoader) Load(ctx context.Context, env map[string]string) error {
	_ = global.RobotRegistry.Register(&robot.SimpleRobot{})

	var conf = global.C.Openai
	if conf.Key != "" {
		for _, model := range conf.Models {
			var chatgptRobot = &chatgpt.ChatGPTRobot{
				BaseUrl:  conf.BaseURL,
				Key:      conf.Key,
				Model:    model,
				Lang:     conf.Lang,
				ProxyUrl: conf.ProxyUrl,
			}
			err := global.RobotRegistry.Register(chatgptRobot)
			if err != nil {
				global.Log.Errorf("chatgpt robot %s register failed: %v", model, err)
				continue
			}
			global.Log.Infof("chatgpt robot %s registered!", model)
		}

	}
	return nil
}

func (loader *RobotLoader) Name() string {
	return "Robot Register Loader"
}

func (loader *RobotLoader) Require() []string {
	return []string{"ZapLoggerLoader"}
}
