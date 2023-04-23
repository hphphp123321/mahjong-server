package global

import "github.com/hphphp123321/mahjong-server/app/service/robot"

var RobotRegistry = robot.NewRegistry()

func init() {
	RobotRegistry.Register(&robot.Simple{})
}
