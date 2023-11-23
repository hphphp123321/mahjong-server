package server

import (
	"github.com/hphphp123321/mahjong-go/mahjong"
	"github.com/hphphp123321/mahjong-server/app/global"
	"github.com/hphphp123321/mahjong-server/app/model/game"
	"github.com/hphphp123321/mahjong-server/app/service/robot"
)

func TestRobot(r robot.Robot) error {
	var g = mahjong.NewMahjongGame(0, nil)
	var players = []*mahjong.Player{
		mahjong.NewMahjongPlayer(),
		mahjong.NewMahjongPlayer(),
		mahjong.NewMahjongPlayer(),
		mahjong.NewMahjongPlayer()}
	var validActions = g.Reset(players, game.GetTestTiles(-1))[mahjong.East]
	var events = g.GetPosEvents(mahjong.East, 0)

	_, err := r.ChooseAction(events, validActions)
	return err
}

func TestAllRobots() {
	var robotRegistry = global.RobotRegistry
	var g = mahjong.NewMahjongGame(0, nil)
	var players = []*mahjong.Player{
		mahjong.NewMahjongPlayer(),
		mahjong.NewMahjongPlayer(),
		mahjong.NewMahjongPlayer(),
		mahjong.NewMahjongPlayer()}
	var validActions = g.Reset(players, nil)[mahjong.East]
	var events = g.GetPosEvents(mahjong.East, 0)

	for robotType, r := range robotRegistry.Robots {
		if _, err := r.ChooseAction(events, validActions); err != nil {
			global.Log.Infoln("robotType:", robotType, "err:", err, "delete robot")
			delete(robotRegistry.Robots, robotType)
		}
	}
	return
}
