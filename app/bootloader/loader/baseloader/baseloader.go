package baseloader

import (
	"context"
	"fmt"
	"github.com/hphphp123321/mahjong-server/app/global"
	"os"
	"path/filepath"
	"strings"
)

// BaseLoader 相关依赖：globals.ProjectRoot
// 1.初始化项目根目录globals.ProjectRoot
type BaseLoader struct {
	ProjectRootDirNameForTest string
}

func (loader *BaseLoader) Require() []string {
	return []string{}
}

func (loader *BaseLoader) Load(ctx context.Context, env map[string]string) error {
	// run || debug || build&./target.exe：wd得到的是main.go/target.exe所在目录
	// test时候的到的wd是 test文件所在目录
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	if mode, ok := env["mode"]; ok && mode == "test" {
		limits := 100

		for ; limits > 0 && !strings.HasSuffix(wd, loader.ProjectRootDirNameForTest); limits-- {
			wd = filepath.Dir(wd)
		}
		if limits <= 0 {
			return fmt.Errorf("未找到项目根目录名：%s", loader.ProjectRootDirNameForTest)
		}

		global.ProjectRoot = wd
		return nil
	}
	global.ProjectRoot = wd
	return nil
}

func (loader *BaseLoader) Name() string {
	return "BaseLoader"
}
