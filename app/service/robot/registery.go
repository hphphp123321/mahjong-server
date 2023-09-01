package robot

import (
	"math/rand"
	"strings"
)

type Registry struct {
	robots map[string]Robot
}

func NewRegistry() *Registry {
	return &Registry{
		robots: make(map[string]Robot),
	}
}

func (r *Registry) Register(robot Robot) {
	r.robots[robot.GetRobotType()] = robot
}

func (r *Registry) Unregister(robotType string) {
	delete(r.robots, robotType)
}

func (r *Registry) GetRobot(robotName string) (robot Robot, ok bool) {
	var robotType string = robotName
	if strings.Contains(robotName, "(com)") {
		robotType = strings.Split(robotName, "(com)")[0]
	}
	robot, ok = r.robots[robotType]
	return robot, ok
}

func (r *Registry) GetRobotTypes() []string {
	robotTypes := make([]string, 0, len(r.robots))
	for robotType := range r.robots {
		robotTypes = append(robotTypes, robotType)
	}
	return robotTypes
}

func (r *Registry) GetRobotCount() int {
	return len(r.robots)
}

func (r *Registry) IsRobotTypeExist(robotType string) bool {
	_, ok := r.robots[robotType]
	return ok
}

func (r *Registry) GetRandomRobotType() string {
	robotTypes := r.GetRobotTypes()
	return robotTypes[rand.Intn(len(robotTypes))]
}

func (r *Registry) GetRandomRobot() Robot {
	robot, _ := r.GetRobot(r.GetRandomRobotType())
	return robot
}
