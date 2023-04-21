package simple

import (
	"github.com/hphphp123321/mahjong-go/mahjong"
	"math/rand"
)

type BaseRobot struct {
}

func (r *BaseRobot) GetRobotType() string {
	return "base"
}

func (r *BaseRobot) ChooseAction(boardState mahjong.BoardState) (actionIdx int, err error) {
	return rand.Intn(len(boardState.ValidActions)), nil
}
