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

func (r *Simple) ChooseAction(boardState mahjong.BoardState) (actionIdx int, err error) {
	return rand.Intn(len(boardState.ValidActions)), nil
}
