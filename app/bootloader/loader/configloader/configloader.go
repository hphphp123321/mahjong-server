package configloader

import (
	"context"
	"os"
	"strings"

	"github.com/hphphp123321/mahjong-server/app/global"
	"gopkg.in/yaml.v3"
)

// ConfigLoader 相关依赖：globals.C
// 1. 加载本地conf/app.yaml配置至globals.C
// 2. 解析Flag
type ConfigLoader struct {
}

func (loader *ConfigLoader) Require() []string {
	return []string{"BaseLoader"}
}

func (loader *ConfigLoader) Load(ctx context.Context, env map[string]string) error {
	target := global.ProjectRoot + "/config/app.yaml"
	if strings.Contains(target, "cmd") {
		target = global.ProjectRoot + "/../config/app.yaml"
	}
	fd, err := os.Open(target)
	if err != nil {
		return err
	}

	err = yaml.NewDecoder(fd).Decode(global.C)
	if err != nil {
		return err
	}

	return nil
}

func (*ConfigLoader) Name() string {
	return "ConfigLoader"
}
