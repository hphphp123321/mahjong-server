package robot

import (
	"github.com/hphphp123321/mahjong-go/mahjong"
)

type Robot interface {
	GetRobotType() string
	ChooseAction(boardState mahjong.BoardState) (actionIdx int, err error)
}
