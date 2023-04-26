package robot

import (
	"github.com/hphphp123321/mahjong-go/mahjong"
	"math/rand"
)

type Simple struct {
}

func (r *Simple) GetRobotType() string {
	return "base"
}

func (r *Simple) ChooseAction(events mahjong.Events, validActions mahjong.Calls) (actionIdx int) {
	return rand.Intn(len(validActions))
}
