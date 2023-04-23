package bootloader

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"
)

type Loader interface {
	Load(ctx context.Context, env map[string]string) error // 加载方法
	Name() string                                          // 需返回加载器的名称，一般为类名
	Require() []string                                     // 需返回该加载器依赖的加载器名称列表
}
type LoaderList []Loader

const (
	Finally        = "FINALLY"
	KeyMode        = "mode"
	ValTestMode    = "test"
	KeyDisableLog  = "DisableLog"
	KeyDisableSort = "DisableSort"
)

var (
	once sync.Once
	env  = map[string]string{}
)

// Load 启动加载
func Load(ctx context.Context, list LoaderList) {
	once.Do(func() {
		t := time.Now()

		if _, ok := env[KeyDisableSort]; !ok {
			list = topoSort(list)
		}

		if _, ok := env[KeyDisableLog]; !ok {
			log.Printf("[Bootloader] 预计加载模块：\n")
			for i, loader := range list {
				log.Printf("[Bootloader] Module(%v/%v): %v\n",
					i+1,
					len(list),
					loader.Name())
			}
		}

		if _, ok := env[KeyDisableLog]; !ok {
			log.Printf("[Bootloader] start (total: %v)...\n", len(list))
		}

		for i, loader := range list {
			//log.Printf("[Bootloaders] start to load: %v\n", loader.Name())
			t := time.Now()
			if err := loader.Load(ctx, env); err != nil {
				panicInfo := fmt.Sprintf("[%s]%s", loader.Name(), err.Error())
				panic(panicInfo)
			}
			if _, ok := env[KeyDisableLog]; !ok {
				log.Printf("[Bootloader] successfully loaded(%v/%v, time: %vms): %v\n",
					i+1,
					len(list),
					time.Since(t).Milliseconds(),
					loader.Name())
			}
			time.Sleep(time.Millisecond * 20)
		}

		if _, ok := env[KeyDisableLog]; !ok {

			log.Printf("[Bootloader] load finished (total: %v), time: %vms\n",
				len(list),
				time.Since(t).Milliseconds())
		}
	})
}

// 根据require拓扑排序
func topoSort(loaderList LoaderList) LoaderList {
	name2Loader := map[string]Loader{}
	Q := make([]string, 0, len(loaderList))
	var finalLoader Loader
	// 建立依赖图结构
	// loader的入边（表示被依赖的）
	beRequiredMap := map[string]map[string]struct{}{}
	// loader的出边（表示依赖的）
	requireMap := map[string]map[string]struct{}{}

	for _, loader := range loaderList {
		// 检查重复
		if _, ok := name2Loader[loader.Name()]; ok {
			panic("重复的Loader, name:" + loader.Name())
		}
		name2Loader[loader.Name()] = loader
		requireList := loader.Require()
		// 检查特殊标记[ALL]
		if len(requireList) != 0 && strings.ToUpper(requireList[0]) == Finally {
			if finalLoader != nil {
				panic("Finally标记不唯一，请检查：[" + finalLoader.Name() + "]以及[" + loader.Name() + "]")
			}
			finalLoader = loader
			continue
		}

		// 将没有依赖的loader加入队列
		if len(requireList) == 0 {
			Q = append(Q, loader.Name())
		}
	}

	for _, loader := range loaderList {
		if finalLoader != nil && loader.Name() == finalLoader.Name() {
			continue
		}

		requireList := loader.Require()
		// 处理出度
		require := requireMap[loader.Name()]
		if require == nil {
			require = map[string]struct{}{}
		}
		for _, targetLoaderName := range requireList {
			if l, ok := name2Loader[targetLoaderName]; !ok || l == nil {
				errMsg := fmt.Sprintf("缺少或未找到%v，因此无法加载%v", targetLoaderName, loader.Name())
				panic(errMsg)
			}
			require[targetLoaderName] = struct{}{}
		}
		requireMap[loader.Name()] = require

		// 处理入度
		for _, targetLoaderName := range requireList {
			beRequired := beRequiredMap[targetLoaderName]
			if beRequired == nil {
				beRequired = map[string]struct{}{}
			}
			beRequired[loader.Name()] = struct{}{}
			beRequiredMap[targetLoaderName] = beRequired
		}
	}

	// 拓扑排序BFS
	result := make([]string, 0, len(loaderList))
	for len(Q) != 0 {
		size := len(Q)
		var subResult []string
		for i := 0; i < size; i++ {
			cur := Q[0]
			Q = Q[1:]
			subResult = append(subResult, cur)

			for targetName, _ := range beRequiredMap[cur] {
				// 删除target loader的该出度
				delete(requireMap[targetName], cur)
				if len(requireMap[targetName]) == 0 {
					Q = append(Q, targetName)
				}
			}
		}
		// 排序是为了稳定
		sort.Strings(subResult)
		result = append(result, subResult...)
	}

	sorted := make(LoaderList, 0, len(result)+1)

	for _, loaderName := range result {
		loader := name2Loader[loaderName]
		sorted = append(sorted, loader)
	}
	if finalLoader != nil {
		sorted = append(sorted, finalLoader)
	}

	return sorted
}

func DisableLog() {
	SetEnv(KeyDisableLog, "")
}

func DisableSort() {
	SetEnv(KeyDisableSort, "")
}

func SetEnv(k, v string) {
	env[k] = v
}

func TestMode() {
	SetEnv(KeyMode, ValTestMode)
}
