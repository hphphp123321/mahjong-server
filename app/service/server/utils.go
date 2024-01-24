package server

import (
	"github.com/hphphp123321/mahjong-go/mahjong"
	query2 "github.com/hphphp123321/mahjong-server/app/dao/query"
	"github.com/hphphp123321/mahjong-server/app/global"
	"github.com/hphphp123321/mahjong-server/app/model/game"
	"github.com/hphphp123321/mahjong-server/app/service/robot"
	"strconv"
)

func CheckUserExist(userDo query2.IUserDo, userName string) (bool, error) {
	do := userDo.Where(query2.User.Name.Eq(userName))
	count, err := do.Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func ConvertStringToUInt(s string) (uint, error) {
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

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

	for _, r := range robotRegistry.ListNotBuiltInRobots() {
		if _, err := r.ChooseAction(events, validActions); err != nil {
			global.Log.Infoln("robotType:", r.GetRobotType(), "err:", err, "delete robot")
			delete(robotRegistry.Robots, r.GetRobotType())
		}
	}
}
