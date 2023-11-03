package robot

import (
	"github.com/hphphp123321/mahjong-go/mahjong"
	"math/rand"
)

// 确保实现了Robot接口
var _ Robot = (*SimpleRobot)(nil)

type SimpleRobot struct {
}

func (r *SimpleRobot) GetRobotType() string {
	return "base"
}

func (r *SimpleRobot) ChooseAction(events mahjong.Events, validActions mahjong.Calls) (actionIdx int) {
	return rand.Intn(len(validActions))
}
