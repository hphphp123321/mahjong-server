package controller

import (
	"fmt"
	"strings"
)

func JoinMsg(name string) string {
	return fmt.Sprintf("Player %s Join the Room", name)
}

func LeaveMsg(name string) string {
	return fmt.Sprintf("Player %s Leave the Room", name)
}

func GetReadyMsg(name string) string {
	return fmt.Sprintf("Player %s Get Ready", name)
}

func CancelReadyMsg(name string) string {
	return fmt.Sprintf("Player %s Cancel Ready", name)
}

func LeaveRoomMsg(name string) string {
	return fmt.Sprintf("Player %s Leave the Room", name)
}

func RemovePlayerMsg(name string) string {
	return fmt.Sprintf("Player %s Removed from the Room", name)
}

func AddRobotMsg(name string) string {
	return fmt.Sprintf("Robot %s Added to the Room", name)
}

func ListRobotMsg(names []string) string {
	return fmt.Sprintf("Robot List: %s", strings.Join(names, ", "))
}

func ChatMsg(name, msg string) string {
	return fmt.Sprintf("Player %s: %s", name, msg)
}

func StartGameMsg() string {
	return fmt.Sprintf("Start Game!")
}
