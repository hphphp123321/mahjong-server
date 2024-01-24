package controller

import "strings"

func LoginMsg(id string) string {
	return "Player Login, ID: " + id
}

func RegisterMsg(id string) string {
	return "Player Register, ID: " + id
}

func LogoutMsg() string {
	return "Player Logout"
}

func CreateMsg(id string) string {
	return " Create Room, ID: " + id
}

func JoinRoomMsg(roomName string) string {
	return "Join Room, Name: " + roomName
}

func ListRoomsMsg(roomNames []string) string {
	return "List Room, Names: " + strings.Join(roomNames, ", ")
}

func RegisterRobotMsg(robotName string) string {
	return "Register Robot, Name: " + robotName
}
