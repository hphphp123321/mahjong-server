package robot

import (
	"github.com/hphphp123321/mahjong-go/mahjong"
)

type Robot interface {
	GetRobotType() string
	ChooseAction(events mahjong.Events, validActions mahjong.Calls) (actionIdx int, err error)
}
