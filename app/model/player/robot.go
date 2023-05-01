package player

import (
	"github.com/hphphp123321/mahjong-go/mahjong"
	"github.com/hphphp123321/mahjong-server/app/errs"
	"github.com/hphphp123321/mahjong-server/app/global"
	"github.com/hphphp123321/mahjong-server/app/service/robot"
)

func InitRobotStream() chan *mahjong.Call {
	return make(chan *mahjong.Call, 20)
}

func StartRobotStream(r robot.Robot, ech chan mahjong.Events, vch chan mahjong.Calls, ach chan *mahjong.Call, error chan error) {
	go func() {
		var events = make(mahjong.Events, 0)
		for {
			select {
			case err := <-error:
				if err == errs.ErrGameEnd {
					global.Log.Debugf("robot %s end", r.GetRobotType())
					return
				}
				global.Log.Warnf("robot %s error: %v", r.GetRobotType(), err)
			case es := <-ech:
				for _, e := range es {
					events = append(events, e)
				}
			case validActions := <-vch:
				actionIdx := r.ChooseAction(events, validActions)
				action := validActions[actionIdx]
				ach <- action
			}
		}
	}()
}
