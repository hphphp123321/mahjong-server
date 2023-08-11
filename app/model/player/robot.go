package player

import (
	"github.com/hphphp123321/mahjong-go/mahjong"
	"github.com/hphphp123321/mahjong-server/app/global"
	"github.com/hphphp123321/mahjong-server/app/service/robot"
)

func InitRobotStream() chan *mahjong.Call {
	return make(chan *mahjong.Call, 20)
}

func StartRobotStream(r robot.Robot, ech chan *GameEventChannel, ach chan *mahjong.Call) {
	go func() {
		var events = make(mahjong.Events, 0)
		for {
			select {
			case ge := <-ech:
				if ge.Events != nil {
					for _, e := range ge.Events {
						events = append(events, e)
					}
				} else if ge.Err != nil {
					global.Log.Warnf("robot %s error: %v", r.GetRobotType(), ge.Err)
				} else if ge.ValidActions != nil {
					actionIdx := r.ChooseAction(events, ge.ValidActions)
					action := ge.ValidActions[actionIdx]
					ach <- action
				}

			}
		}
	}()
}
